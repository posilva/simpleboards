package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/posilva/simpleboards/internal/core/ports"
)

const (
	// ddb fields
	hashKeyName string = "pk"
	sortKeyName string = "sk"
	// ddb prefixes
	pkUserPrefix        string = "USR#"
	skLeaderboardPrefix string = "LBRD#"

	queryTimeout   = 1 * time.Second
	pkConfigPrefix = "LBRD#CONFIG"
	skConfigPrefix = "LBRD#NAME#"
	scoreAttrib    = "score"
)

// DDBConfigItem ...
type DDBConfigItem struct {
	PK     string `dynamodbav:"pk"`
	SK     string `dynamodbav:"sk"`
	Config string `dynamodbav:"config"`
}

// LeaderboardEntryRecord represents a dynamodb table record
type LeaderboardEntryRecord struct {
	PK      string  `dynamodbav:"pk" json:"pk"`
	SK      string  `dynamodbav:"sk" json:"sk"`
	Score   float64 `dynamodbav:"score" json:"score"`
	Counter uint64  `dynamodbav:"counter" json:"counter"`
}

// DynamoDBRepository implements Repository interface for DynamoDB
type DynamoDBRepository struct {
	log       ports.Logger
	client    DynamoDBClient
	tableName string
}

// NewDynamoDBRepository creates a new DynamoDB repository
func NewDynamoDBRepository(settings DynamoDBSettings) (*DynamoDBRepository, error) {
	return &DynamoDBRepository{
		client:    settings.Client,
		tableName: settings.Table,
		log:       settings.Logger,
	}, nil
}

// AddWithMetadata the value to the entry
func (r *DynamoDBRepository) AddWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error) {
	builder := expression.NewBuilder()
	update := expression.Add(
		expression.Name("score"),
		expression.Value(value),
	).Add(
		expression.Name("counter"),
		expression.Value(1),
	)
	r.log.Debug("updating entry", "entry", entry, "leaderboard", leaderboard, "value", value, "function", "Add")
	builder = builder.WithUpdate(update)

	expr, err := builder.Build()
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to build update expression: %w", err)
	}
	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
		Key: map[string]types.AttributeValue{
			hashKeyName: &types.AttributeValueMemberS{Value: pkValue(entry)},
			sortKeyName: &types.AttributeValueMemberS{Value: skValue(leaderboard)},
		},
		UpdateExpression: expr.Update(),
	}

	output, err := r.client.UpdateItem(context.Background(), &input)
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to update item: %w", err)
	}

	s := LeaderboardEntryRecord{}
	err = attributevalue.UnmarshalMap(output.Attributes, &s)
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to process output: %w", err)
	}

	res := domain.ScoreUpdate{Score: s.Score, Counter: s.Counter, Done: true}
	return res, nil
}

// MaxWithMetadata ...
func (r *DynamoDBRepository) MaxWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error) {
	builder := expression.NewBuilder()
	update := expression.Set(
		expression.Name(scoreAttrib),
		expression.Value(value),
	).Add(
		expression.Name("counter"),
		expression.Value(1),
	)
	r.log.Info("updating entry", "entry", entry, "leaderboard", leaderboard, "value", value, "function", "Max")

	// just stores if the existing score value is less than the one to be stored
	condBuilder := expression.Name(scoreAttrib).LessThanEqual(expression.Value(value)).Or(expression.Name(scoreAttrib).AttributeNotExists())
	expr, err := builder.WithUpdate(update).WithCondition(condBuilder).Build()
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to build update expression: %w", err)
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
		ConditionExpression:       expr.Condition(),
		Key: map[string]types.AttributeValue{
			hashKeyName: &types.AttributeValueMemberS{Value: pkValue(entry)},
			sortKeyName: &types.AttributeValueMemberS{Value: skValue(leaderboard)},
		},
		UpdateExpression: expr.Update(),
	}

	output, err := r.client.UpdateItem(context.Background(), &input)
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return domain.ScoreUpdate{Done: false}, nil
		}
		return domain.ScoreUpdate{}, fmt.Errorf("failed to update item: %w", err)
	}
	s := LeaderboardEntryRecord{}
	err = attributevalue.UnmarshalMap(output.Attributes, &s)
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to process output: %w", err)
	}

	return domain.ScoreUpdate{Score: s.Score, Done: true}, nil
}

// MinWithMetadata ...
func (r *DynamoDBRepository) MinWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error) {
	builder := expression.NewBuilder()
	update := expression.Set(
		expression.Name(scoreAttrib),
		expression.Value(value),
	).Add(
		expression.Name("counter"),
		expression.Value(1),
	)
	r.log.Debug("updating entry", "entry", entry, "leaderboard", leaderboard, "value", value, "function", "Min")

	// just stores if the existing score value is greater than the one to be stored
	condBuilder := expression.Name(scoreAttrib).GreaterThanEqual(expression.Value(value)).Or(expression.Name(scoreAttrib).AttributeNotExists())

	expr, err := builder.WithUpdate(update).WithCondition(condBuilder).Build()
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to build update expression: %w", err)
	}
	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
		ConditionExpression:       expr.Condition(),
		Key: map[string]types.AttributeValue{
			hashKeyName: &types.AttributeValueMemberS{Value: pkValue(entry)},
			sortKeyName: &types.AttributeValueMemberS{Value: skValue(leaderboard)},
		},
		UpdateExpression: expr.Update(),
	}

	output, err := r.client.UpdateItem(context.Background(), &input)
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return domain.ScoreUpdate{Done: false}, nil
		}
		return domain.ScoreUpdate{}, fmt.Errorf("failed to update item: %w", err)
	}

	s := LeaderboardEntryRecord{}
	err = attributevalue.UnmarshalMap(output.Attributes, &s)
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to process output: %w", err)
	}
	return domain.ScoreUpdate{Score: s.Score, Done: true}, nil
}

// LastWithMetadata ...
func (r *DynamoDBRepository) LastWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error) {
	builder := expression.NewBuilder()
	update := expression.Set(
		expression.Name(scoreAttrib),
		expression.Value(value),
	).Add(
		expression.Name("counter"),
		expression.Value(1),
	)

	r.log.Info("updating entry", "entry", entry, "leaderboard", leaderboard, "value", value, "function", "Last")

	expr, err := builder.WithUpdate(update).Build()
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to build update expression: %w", err)
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
		Key: map[string]types.AttributeValue{
			hashKeyName: &types.AttributeValueMemberS{Value: pkValue(entry)},
			sortKeyName: &types.AttributeValueMemberS{Value: skValue(leaderboard)},
		},
		UpdateExpression: expr.Update(),
	}

	output, err := r.client.UpdateItem(context.Background(), &input)
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to update item: %w", err)
	}

	s := LeaderboardEntryRecord{}
	err = attributevalue.UnmarshalMap(output.Attributes, &s)
	if err != nil {
		return domain.ScoreUpdate{}, fmt.Errorf("failed to process output: %w", err)
	}
	return domain.ScoreUpdate{Score: s.Score, Done: true}, nil
}

// GetConfig returns all existing leaderboards
func (r *DynamoDBRepository) GetConfig() (domain.LeaderboardsConfigMap, error) {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 1*time.Second, errors.New("get configuration timeout"))
	defer cancel()

	keyCond := expression.KeyAnd(
		expression.Key(hashKeyName).Equal(expression.Value(pkConfigPrefix)),
		expression.Key(sortKeyName).BeginsWith(skConfigPrefix),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression: %v", err)
	}
	input := dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	output, err := r.client.Query(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %v", err)
	}

	var configMap domain.LeaderboardsConfigMap = make(map[string]domain.LeaderboardConfig, len(output.Items))
	for _, item := range output.Items {
		var it DDBConfigItem
		err = attributevalue.UnmarshalMap(item, &it)
		if err != nil {
			break
		}
		var cfg domain.LeaderboardConfig
		err = json.Unmarshal([]byte(it.Config), &cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Json config for '%v': %v", it.SK, err)
		}
		configMap[cfg.Name] = cfg
	}
	return configMap, nil
}

// ResetLock implements the ResetLocker interface
func (r *DynamoDBRepository) ResetLock(leaderboard string, epoch int, duration time.Duration) (bool, error) {
	return false, nil
}

// Update configuration
func (r *DynamoDBRepository) Update(name string, config domain.LeaderboardConfig) error {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 1*time.Second, errors.New("get configuration timeout"))
	defer cancel()

	skValue := fmt.Sprintf("%s%s", skConfigPrefix, name)

	cfg, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %v ", err)
	}
	configItem := DDBConfigItem{
		PK:     pkConfigPrefix,
		SK:     skValue,
		Config: string(cfg),
	}
	item, err := attributevalue.MarshalMap(configItem)
	if err != nil {
		panic(fmt.Sprintf("failed to marshalMap: %v ", err))
	}
	input := dynamodb.PutItemInput{
		TableName:    aws.String(r.tableName),
		Item:         item,
		ReturnValues: types.ReturnValueNone,
	}

	_, err = r.client.PutItem(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}
	return nil
}

// Add ...
func (r *DynamoDBRepository) Add(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error) {
	return r.AddWithMetadata(entry, leaderboard, value, nil)
}

// Min ...
func (r *DynamoDBRepository) Min(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error) {
	return r.MinWithMetadata(entry, leaderboard, value, nil)
}

// Max ...
func (r *DynamoDBRepository) Max(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error) {
	return r.MaxWithMetadata(entry, leaderboard, value, nil)
}

// Last ...
func (r *DynamoDBRepository) Last(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error) {
	return r.LastWithMetadata(entry, leaderboard, value, nil)
}

func pkValue(value string) string {
	return fmt.Sprintf("%s%s", pkUserPrefix, value)
}

func skValue(value string) string {
	return fmt.Sprintf("%s%s", skLeaderboardPrefix, value)
}
