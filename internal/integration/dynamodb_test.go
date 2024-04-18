//go:build integration

package integration

import (
	"testing"

	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDynamoDBRepository_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	settings := testutil.NewDefaultDynamoDBSettings()

	r, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)

	entry := testutil.NewID()
	leaderboard := testutil.NewUnique(testutil.Name(t))
	score := float64(1)

	_, err = r.Add(entry, leaderboard, score)
	assert.NoError(t, err)

	v1, err := r.Add(entry, leaderboard, score)
	assert.NoError(t, err)

	assert.Equal(t, uint64(2*score), uint64(v1.Score))
}

func TestDynamoDBRepository_Max(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	settings := testutil.NewDefaultDynamoDBSettings()

	r, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)

	entry := testutil.NewID()
	leaderboard := testutil.NewUnique(testutil.Name(t))
	score := float64(10)

	_, err = r.Max(entry, leaderboard, score)
	assert.NoError(t, err)

	score2 := float64(0)

	v1, err := r.Max(entry, leaderboard, score2)
	assert.NoError(t, err)

	assert.False(t, v1.Done)
	assert.Equal(t, int(v1.Score), 0)
}

func TestDynamoDBRepository_Min(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	settings := testutil.NewDefaultDynamoDBSettings()

	r, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)

	entry := testutil.NewID()
	leaderboard := testutil.NewUnique(testutil.Name(t))
	score := float64(10)

	_, err = r.Min(entry, leaderboard, score)
	assert.NoError(t, err)

	score2 := float64(5)
	v1, err := r.Min(entry, leaderboard, score2)
	assert.NoError(t, err)

	assert.Equal(t, uint64(score2), uint64(v1.Score))
}

func TestDynamoDBRepository_Last(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	settings := testutil.NewDefaultDynamoDBSettings()

	r, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)

	entry := testutil.NewID()
	leaderboard := testutil.NewUnique(testutil.Name(t))
	score := float64(1)

	_, err = r.Add(entry, leaderboard, score)
	assert.NoError(t, err)

	v1, err := r.Add(entry, leaderboard, score)
	assert.NoError(t, err)

	assert.Equal(t, uint64(2*score), uint64(v1.Score))
}
