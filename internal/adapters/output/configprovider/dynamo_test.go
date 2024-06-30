package configprovider

import (
	"sync"
	"testing"
	"time"

	"github.com/posilva/simpleboards/internal/core/services"
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

func TestDynamoDBConfigProvider(t *testing.T) {
	counter := testCounter{}
	s := services.NewScheduler(1, counter.Add)
	assert.NotNil(t, s)
	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, 2, counter.Get())
}
