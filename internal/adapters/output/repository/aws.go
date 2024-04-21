// Package repository is Repository interface implementations
package repository

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/posilva/simpleboards/internal/core/ports"
)

type DynamoDBSettings struct {
	Table  string
	Logger ports.Logger
	Client DynamoDBClient
}

// This insterface is used to be able to execute unit tests
type DynamoDBClient interface {
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// NewDynamoDBClientFromConfig creates a new DynamoDB
func NewDynamoDBClientFromConfig(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

// DefaultLocalAWSClientConfig returns the default local AWS config
func DefaultLocalAWSClientConfig() aws.Config {
	host := "http://localhost:4566" // default value pointing for local stack
	region := "us-east-1"

	if v, ok := os.LookupEnv("AWS_ENDPOINT"); ok {
		host = v
	}
	if v, ok := os.LookupEnv("AWS_DEFAULT_REGION"); ok {
		region = v
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		// config.WithClientLogMode(aws.LogRetries|aws.LogRequest|aws.LogRequestWithBody|aws.LogResponse|aws.LogResponseWithBody|aws.LogDeprecatedUsage|aws.LogRequestEventMessage|aws.LogResponseEventMessage),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, opts ...any) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: host,
				}, nil
			},
		)),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}
