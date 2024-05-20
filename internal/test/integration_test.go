package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/rueidis"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	localstack "github.com/testcontainers/testcontainers-go/modules/localstack"
	testcontainersredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type LeaderboardsTestSuite struct {
	suite.Suite
	ctx         context.Context
	rdContainer *testcontainersredis.RedisContainer
	rdClient    rueidis.Client

	ddbContainer *localstack.LocalStackContainer
}

func (suite *LeaderboardsTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	redisContainer, err := testcontainersredis.RunContainer(
		suite.ctx,
		testcontainers.WithImage("redis:latest"),
		testcontainers.WithWaitStrategyAndDeadline(
			10*time.Second, wait.ForExposedPort()),
	)
	suite.NoError(err)

	ip, err := redisContainer.Host(suite.ctx)
	suite.NoError(err)
	port, err := redisContainer.MappedPort(suite.ctx, "6379")
	suite.NoError(err)

	endpoint := fmt.Sprintf("%s:%s", ip, port.Port())
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{endpoint},
	})

	suite.NoError(err)

	pingCmd := redisClient.B().Ping().Build()
	err = redisClient.Do(suite.ctx, pingCmd).Error()
	suite.NoError(err)

	suite.rdContainer = redisContainer
	suite.rdClient = redisClient

	ddbContainer, err := localstack.RunContainer(suite.ctx,
		testcontainers.WithImage("localstack/localstack:latest"),
		testcontainers.WithWaitStrategyAndDeadline(
			30*time.Second, wait.ForHealthCheck().WithPollInterval(1*time.Second)))
	suite.NoError(err)

	suite.ddbContainer = ddbContainer

}

func (suite *LeaderboardsTestSuite) TearDownSuite() {
	err := suite.rdContainer.Terminate(suite.ctx)
	suite.NoError(err)
	err = suite.ddbContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *LeaderboardsTestSuite) SetupTest() {

}

func (suite *LeaderboardsTestSuite) TearDownTest() {
	err := suite.rdClient.Do(suite.ctx, suite.rdClient.B().Flushall().Build()).Error()
	suite.NoError(err)
}

func (suite *LeaderboardsTestSuite) BeforeTest(_ string, testName string) {
	fmt.Printf("Running before test: %s\n", testName)
}

func (suite *LeaderboardsTestSuite) TestLeaderboards() {
}

func TestLeaderboards(t *testing.T) {
	suite.Run(t, new(LeaderboardsTestSuite))
}

func initDDB() {

}
