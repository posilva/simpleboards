package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/google/uuid"
	"github.com/phayes/freeport"
	"github.com/posilva/simpleboards/cmd/simpleboards/app"
	"github.com/posilva/simpleboards/cmd/simpleboards/config"
	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/posilva/simpleboards/internal/testutil"

	"github.com/redis/rueidis"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	testcontainersredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	uniqueTestID      = uuid.NewString()
	defaultLbName     = "integration_lb_tests::" + uniqueTestID
	defaultLbNameSum  = defaultLbName + "::Sum"
	defaultLbNameMax  = defaultLbName + "::Max"
	defaultLbNameMin  = defaultLbName + "::Min"
	defaultLbNameLast = defaultLbName + "::Last"
	metadataDefault   = map[string]string{
		"country": "PT",
		"league":  "gold",
	}
)

type BaseTestSuite struct {
	suite.Suite
	Context            context.Context
	RedisContainer     *testcontainersredis.RedisContainer
	RedisClient        rueidis.Client
	DDBContainer       *localstack.LocalStackContainer
	AWSConfig          aws.Config
	DDBClient          *dynamodb.Client
	RedisEndpoint      string
	LocalstackEndpoint string
	ServiceEndpoint    string
}

func waitForService(suite *BaseTestSuite) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	c := &http.Client{Timeout: time.Second * 10}
	for {
		r, err := c.Get(fmt.Sprintf("http://%s/", suite.ServiceEndpoint))
		if err == nil && r.StatusCode == http.StatusOK {
			return
		}
		select {
		case <-time.After(time.Millisecond * 100):
		case <-ctx.Done():
			require.Fail(suite.T(), "timeout waiting for health check")
		}
	}
}

func setup(suite *BaseTestSuite) {
	log.Println("Running setup suite")

	suite.Context = context.Background()
	testcontainers.Logger = log.New(&ioutils.NopWriter{}, "", 0)
	setupRedisContainer(suite)
	setupDDBContainer(suite)
	createTable(suite)

	configLeaderboard(defaultLbName)
	f := domain.Sum
	r := domain.Hourly
	configLeaderboardFuncReset(defaultLbNameSum, f, r)
	f = domain.Max
	configLeaderboardFuncReset(defaultLbNameMax, f, r)
	f = domain.Min
	configLeaderboardFuncReset(defaultLbNameMin, f, r)
	f = domain.Last
	configLeaderboardFuncReset(defaultLbNameLast, f, r)

	port, err := freeport.GetFreePort()
	if err != nil {
		panic("failed to get free port: " + err.Error())
	}

	suite.ServiceEndpoint = fmt.Sprintf("127.0.0.1:%d", port)
	log.Println("Service endpoint: ", suite.ServiceEndpoint)
	config.SetAddr(suite.ServiceEndpoint)
	config.SetRedisAddr(suite.RedisEndpoint)
	config.SetDynamoDBTableName(testutil.DynamoDBLocalTableName)
	config.SetLocal(true)

	go func() {
		app.Run()
	}()
	waitForService(suite)
	log.Printf("Service is running on %s", suite.ServiceEndpoint)
}

func setupRedisContainer(suite *BaseTestSuite) {
	redisContainer, err := testcontainersredis.RunContainer(
		suite.Context,
		testcontainers.WithImage("redis:latest"),
		testcontainers.WithWaitStrategyAndDeadline(
			30*time.Second, wait.ForExposedPort()),
	)
	suite.NoError(err)

	ip, err := redisContainer.Host(suite.Context)
	suite.NoError(err)
	port, err := redisContainer.MappedPort(suite.Context, "6379")
	suite.NoError(err)

	endpoint := fmt.Sprintf("%s:%s", ip, port.Port())
	log.Printf("Redis endpoint: %s", endpoint)
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{endpoint},
	})

	suite.NoError(err)
	suite.RedisEndpoint = endpoint
	pingCmd := redisClient.B().Ping().Build()
	err = redisClient.Do(suite.Context, pingCmd).Error()
	suite.NoError(err)

	suite.RedisContainer = redisContainer
	suite.RedisClient = redisClient
}

func setupDDBContainer(suite *BaseTestSuite) {
	log.Println("Running setup DDB container")
	ddbContainer, err := localstack.RunContainer(suite.Context,
		testcontainers.WithImage("localstack/localstack:latest"),
		testcontainers.WithWaitStrategyAndDeadline(
			30*time.Second, wait.ForHealthCheck().WithPollInterval(1*time.Second)))
	suite.NoError(err)

	host, err := ddbContainer.Host(suite.Context)
	suite.NoError(err)

	port, err := ddbContainer.MappedPort(suite.Context, "4566")
	suite.NoError(err)

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())
	log.Printf("Localstack endpoint: %s", endpoint)
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(aws.AnonymousCredentials{}),
		awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, opts ...any) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			},
		)),
	)
	suite.DDBContainer = ddbContainer
	suite.AWSConfig = cfg
	suite.LocalstackEndpoint = endpoint

	os.Setenv("AWS_ENDPOINT", endpoint)
	suite.DDBClient = dynamodb.NewFromConfig(suite.AWSConfig)
}

func createTable(suite *BaseTestSuite) {
	_, err := suite.DDBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("pk"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("sk"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("pk"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("sk"),
			KeyType:       types.KeyTypeRange,
		}},
		TableName: aws.String(testutil.DynamoDBLocalTableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	suite.NoError(err)
	waiter := dynamodb.NewTableExistsWaiter(suite.DDBClient)
	err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(testutil.DynamoDBLocalTableName),
	}, 10*time.Second)
	suite.NoError(err)
}

func teardown(suite *BaseTestSuite) {
	err := suite.RedisContainer.Terminate(suite.Context)
	suite.NoError(err)
	err = suite.DDBContainer.Terminate(suite.Context)

	suite.NoError(err)
	//err = suite.RedisClient.Do(suite.Context, suite.RedisClient.B().Flushall().Build()).Error()
	//suite.NoError(err)
}

func configLeaderboard(lbName string) {
	settings := testutil.NewDefaultDynamoDBSettings()
	repo, err := repository.NewDynamoDBRepository(settings)
	if err != nil {
		panic(err)
	}
	lbConfig := testutil.NewLeaderboardConfigWithScoreboards(lbName, domain.Hourly, domain.Sum)
	err = repo.Update(lbName, lbConfig)
	if err != nil {
		panic(err)
	}
	_, err = repo.GetConfig()
	if err != nil {
		panic(err)
	}
}
func configLeaderboardFuncReset(lbName string, f domain.LeaderboardFunctionType, r domain.LeaderboardResetType) {
	settings := testutil.NewDefaultDynamoDBSettings()
	repo, err := repository.NewDynamoDBRepository(settings)
	if err != nil {
		panic(err)
	}
	lbConfig := testutil.NewLeaderboardConfigWithFunctionResetWithScoreboards(lbName, r, f)
	err = repo.Update(lbName, lbConfig)
	if err != nil {
		panic(err)
	}

	_, err = repo.GetConfig()
	if err != nil {
		panic(err)
	}
}

/**
TEMPLATE OF INTERGRATION TEST
package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExampleTestSuite struct {
	BaseTestSuite
}

func (suite *ExampleTestSuite) SetupSuite() {
	setup(&suite.BaseTestSuite)
	fmt.Println("Running setup suite")
}

func (suite *ExampleTestSuite) TearDownSuite() {
	fmt.Println("Running teardown suite")
	teardown(&suite.BaseTestSuite)
}

func (suite *ExampleTestSuite) SetupTest() {
	fmt.Println("Running setup test")
}

func (suite *ExampleTestSuite) TearDownTest() {
	fmt.Println("Running teardown test")
	err := suite.RedisClient.Do(suite.Context, suite.RedisClient.B().Flushall().Build()).Error()
	suite.NoError(err)
}

func (suite *ExampleTestSuite) BeforeTest(_ string, testName string) {
	fmt.Printf("Running before test: %s\n", testName)
}

func (suite *ExampleTestSuite) AfterTest(_ string, testName string) {
	fmt.Printf("Running after test: %s\n", testName)
}

func (suite *ExampleTestSuite) TestExample() {
}

func TestExample(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

*/
