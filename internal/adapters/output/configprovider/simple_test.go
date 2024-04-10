//go:build integration

package configprovider

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/core/services"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/assert"
)

type testCounter struct {
	sync.Mutex
	ct int
}

func (tc *testCounter) Add() {
	tc.Lock()
	defer tc.Unlock()
	tc.ct = tc.ct + 1
}

func (tc *testCounter) Get() int {
	tc.Lock()
	defer tc.Unlock()
	return tc.ct
}

func TestSimpleConfigProvider(t *testing.T) {
	fmt.Printf("Init Exec Job: %v\n", time.Now().Unix())
	counter := testCounter{}
	s := services.NewScheduler(1, counter.Add)
	assert.NotNil(t, s)
	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, 2, counter.Get())
}

func TestSetupConfig(t *testing.T) {
	settings := testutil.NewDefaultDynamoDBSettings()
	repo, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)

	cp := NewSimpleConfigProvider(repo, settings.Logger)
	assert.NotNil(t, cp)

	lbName := "some_leaderboard"
	lbConfig := testutil.NewLeaderboardConfig(lbName, 1, 1, "reward 1")
	err = repo.Update(lbName, lbConfig)
	assert.NoError(t, err)
}

func TestDynamoDBProviderUpdate(t *testing.T) {
	settings := testutil.NewDefaultDynamoDBSettings()
	repo, err := repository.NewDynamoDBRepository(settings)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	lbName := testutil.NewUnique(testutil.Name(t))
	lbConfig := testutil.NewLeaderboardConfig(lbName, 1, 1, "reward 1")
	err = repo.Update(lbName, lbConfig)

	assert.NoError(t, err)

	cp := NewSimpleConfigProvider(repo, settings.Logger)
	assert.NotNil(t, cp)
	time.Sleep(1 * time.Second)

	c, err := cp.Provide()

	assert.NotNil(t, c)
	assert.NoError(t, err)
}
