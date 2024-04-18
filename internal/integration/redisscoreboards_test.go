package integration

import (
	"testing"

	"github.com/posilva/simpleboards/internal/adapters/output/scoreboard"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	redisAddr = "localhost:6379"
)

func TestNewRedisScoreboard(t *testing.T) {
	board, err := scoreboard.NewRedisScoreboard(redisAddr)
	assert.NoError(t, err)
	assert.NotNil(t, board)
}

func TestAddScore(t *testing.T) {
	board, err := scoreboard.NewRedisScoreboard(redisAddr)
	assert.NoError(t, err)
	assert.NotNil(t, board)

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()

	err = board.AddScore(entryID, lbName, 1)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	board, err := scoreboard.NewRedisScoreboard(redisAddr)
	assert.NoError(t, err)
	assert.NotNil(t, board)

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()
	entryID2 := testutil.NewID()

	err = board.AddScore(entryID, lbName, 5)
	assert.Nil(t, err)

	err = board.AddScore(entryID2, lbName, 10)
	assert.Nil(t, err)

	r, err := board.Get(lbName)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Len(t, r, 2)
	assert.Equal(t, entryID2, r[0].EntryID)
	assert.Equal(t, entryID, r[1].EntryID)
}

func TestGetRank(t *testing.T) {
	board, err := scoreboard.NewRedisScoreboard(redisAddr)
	assert.NoError(t, err)
	assert.NotNil(t, board)

	lbName := testutil.NewUnique(testutil.Name(t))

	entryID := testutil.NewID()
	err = board.AddScore(entryID, lbName, 5)
	assert.Nil(t, err)

	entryID = testutil.NewID()

	err = board.AddScore(entryID, lbName, 25)
	assert.Nil(t, err)

	entryID = testutil.NewID()
	err = board.AddScore(entryID, lbName, 50)
	assert.Nil(t, err)

	entryID = testutil.NewID()
	err = board.AddScore(entryID, lbName, 45)
	assert.Nil(t, err)

	r, err := board.GetRank(lbName, entryID)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r, uint64(2))
}
