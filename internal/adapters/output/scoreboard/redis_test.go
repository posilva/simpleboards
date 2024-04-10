// go:build integration
package scoreboard

import (
	"testing"

	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCache(t *testing.T) {
	cache, err := NewRedisScoreboard("localhost:6379")
	assert.NoError(t, err)
	assert.NotNil(t, cache)
}

func TestAddScore(t *testing.T) {
	cache, err := NewRedisScoreboard("localhost:6379")
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()

	err = cache.AddScore(entryID, lbName, 1)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	cache, err := NewRedisScoreboard("localhost:6379")
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()
	entryID2 := testutil.NewID()
	err = cache.AddScore(entryID, lbName, 5)
	assert.Nil(t, err)
	err = cache.AddScore(entryID2, lbName, 10)
	assert.Nil(t, err)

	r, err := cache.Get(lbName)
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Len(t, r, 2)
	assert.Equal(t, entryID2, r[0].EntryID)
	assert.Equal(t, entryID, r[1].EntryID)
}

func TestGetRank(t *testing.T) {
	cache, err := NewRedisScoreboard("localhost:6379")
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()
	err = cache.AddScore(testutil.NewID(), lbName, 5)
	assert.Nil(t, err)
	err = cache.AddScore(testutil.NewID(), lbName, 25)
	assert.Nil(t, err)
	err = cache.AddScore(testutil.NewID(), lbName, 50)
	assert.Nil(t, err)
	err = cache.AddScore(entryID, lbName, 45)
	assert.Nil(t, err)

	r, err := cache.GetRank(lbName, entryID)
	assert.Nil(t, err)
	assert.NotNil(t, r)
}
