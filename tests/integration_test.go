package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

var defaultLeaderboardName = "integration_lb_tests"

type E2ETestSuite struct {
	BaseTestSuite
}

func (suite *E2ETestSuite) SetupSuite() {
	setup(&suite.BaseTestSuite)
	fmt.Println("Running setup suite")
}

func (suite *E2ETestSuite) TearDownSuite() {
	fmt.Println("Running teardown suite")
	teardown(&suite.BaseTestSuite)
}

func (suite *E2ETestSuite) SetupTest() {
	fmt.Println("Running setup test")
}

func (suite *E2ETestSuite) TearDownTest() {
	fmt.Println("Running teardown test")
	err := suite.RedisClient.Do(suite.Context, suite.RedisClient.B().Flushall().Build()).Error()
	suite.NoError(err)
}

func (suite *E2ETestSuite) BeforeTest(_ string, testName string) {
	fmt.Printf("Running before test: %s\n", testName)
}

func (suite *E2ETestSuite) AfterTest(_ string, testName string) {
	fmt.Printf("Running after test: %s\n", testName)
}

func (suite *E2ETestSuite) TestLeaderboard() {
	suite.Assert().True(true, "This is true")
}

func TestLeaderboards(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
