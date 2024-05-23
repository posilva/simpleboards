package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	localstack "github.com/testcontainers/testcontainers-go/modules/localstack"
	testcontainersredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type BaseTestSuite struct {
	suite.Suite
	Context        context.Context
	RedisContainer *testcontainersredis.RedisContainer
	RedisClient    rueidis.Client
	DDBContainer   *localstack.LocalStackContainer
	AWSConfig      aws.Config
	DDBClient      *dynamodb.Client
}

func setup(suite *BaseTestSuite) {
	fmt.Println("Running setup suite")

	suite.Context = context.Background()
	setupRedisContainer(suite)
	setupDDBContainer(suite)
	createTable(suite)
}

func setupRedisContainer(suite *BaseTestSuite) {
	redisContainer, err := testcontainersredis.RunContainer(
		suite.Context,
		testcontainers.WithImage("redis:latest"),
		testcontainers.WithWaitStrategyAndDeadline(
			10*time.Second, wait.ForExposedPort()),
	)
	suite.NoError(err)

	ip, err := redisContainer.Host(suite.Context)
	suite.NoError(err)
	port, err := redisContainer.MappedPort(suite.Context, "6379")
	suite.NoError(err)

	endpoint := fmt.Sprintf("%s:%s", ip, port.Port())
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{endpoint},
	})

	suite.NoError(err)

	pingCmd := redisClient.B().Ping().Build()
	err = redisClient.Do(suite.Context, pingCmd).Error()
	suite.NoError(err)

	suite.RedisContainer = redisContainer
	suite.RedisClient = redisClient
}

func setupDDBContainer(suite *BaseTestSuite) {
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

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, opts ...any) (aws.Endpoint, error) {
				fmt.Println(endpoint)
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			},
		)),
	)
	suite.DDBContainer = ddbContainer
	suite.AWSConfig = cfg
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
		TableName: aws.String(testutil.DynamoDBLocalTableName)}, 10*time.Second)
	suite.NoError(err)
}
func teardown(suite *BaseTestSuite) {
	err := suite.RedisContainer.Terminate(suite.Context)
	suite.NoError(err)
	err = suite.DDBContainer.Terminate(suite.Context)
	suite.NoError(err)
}
