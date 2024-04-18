package scoreboard

import (
	"context"
	"testing"

	"github.com/posilva/simpleboards/internal/testutil"
	mock "github.com/redis/rueidis/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewRedisScoreboard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewClient(ctrl)
	board := NewRedisScoreboardWithClient(c)
	assert.NotNil(t, board)
}

func TestAddScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	board := NewRedisScoreboardWithClient(c)
	assert.NotNil(t, board)

	ctx := context.Background()
	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()

	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "1", entryID)).Return(mock.Result(mock.RedisString("does-not-matter")))

	err := board.AddScore(entryID, lbName, 1)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	board := NewRedisScoreboardWithClient(c)
	assert.NotNil(t, board)

	ctx := context.Background()
	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()
	entryID2 := testutil.NewID()

	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "5", entryID)).Return(mock.Result(mock.RedisString("does-not-matter")))
	err := board.AddScore(entryID, lbName, 5)
	assert.Nil(t, err)

	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "10", entryID2)).Return(mock.Result(mock.RedisString("does-not-matter")))
	err = board.AddScore(entryID2, lbName, 10)
	assert.Nil(t, err)

	c.EXPECT().Do(ctx, mock.Match("ZREVRANGE", lbName, "0", "50", "WITHSCORES")).Return(mock.Result(mock.RedisArray(
		mock.RedisString(entryID2),
		mock.RedisString("10"),
		mock.RedisString(entryID),
		mock.RedisString("5"),
	)))
	r, err := board.Get(lbName)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Len(t, r, 2)
	assert.Equal(t, entryID2, r[0].EntryID)
	assert.Equal(t, entryID, r[1].EntryID)
}

func TestGetRank(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	board := NewRedisScoreboardWithClient(c)
	assert.NotNil(t, board)
	ctx := context.Background()
	lbName := testutil.NewUnique(testutil.Name(t))

	entryID := testutil.NewID()
	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "5", entryID)).Return(mock.Result(mock.RedisString("does-not-matter")))
	err := board.AddScore(entryID, lbName, 5)
	assert.Nil(t, err)

	entryID = testutil.NewID()

	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "25", entryID)).Return(mock.Result(mock.RedisString("does-not-matter")))
	err = board.AddScore(entryID, lbName, 25)
	assert.Nil(t, err)

	entryID = testutil.NewID()
	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "50", entryID)).Return(mock.Result(mock.RedisString("does-not-matter")))
	err = board.AddScore(entryID, lbName, 50)
	assert.Nil(t, err)

	entryID = testutil.NewID()
	c.EXPECT().Do(ctx, mock.Match("ZADD", lbName, "45", entryID)).Return(mock.Result(mock.RedisString("does-not-matter")))
	err = board.AddScore(entryID, lbName, 45)
	assert.Nil(t, err)

	c.EXPECT().Do(ctx, mock.Match("ZREVRANK", lbName, entryID)).Return(mock.Result(mock.RedisInt64(2)))
	r, err := board.GetRank(lbName, entryID)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r, uint64(3))
}
