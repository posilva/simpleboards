// Package testutil is used to share test utilities
package testutil

import (
	"strings"
	"testing"

	"github.com/posilva/simpleboards/internal/adapters/output/logging"
	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/core/domain"
	uuid "github.com/segmentio/ksuid"
)

const (
	// DynamoDBLocalTableName defines the DynamoDB table name for local development with LocalStack
	DynamoDBLocalTableName string = "sgs-gbl-dev-leaderboards"
	// RabbitMQLocalURL defines the local url to connect to Rabbit MQ
	RabbitMQLocalURL string = "amqp://guest:guest@localhost:5672/"
	// RabbitMQLocalURLSSL defines the local url to connect to Rabbit MQ using SSL
	RabbitMQLocalURLSSL string = "amqps://guest:guest@localhost:5671/"
)

// Name returns the name of the test
func Name(t *testing.T) string {
	return t.Name()
}

// NewID returns an ID for tests using kuid package
func NewID() string {
	return strings.ToLower(uuid.New().String())
}

// NewUnique appends to a string a UUID to allow for uniqueness
func NewUnique(prefix string) string {
	return strings.ToLower(prefix + NewID())
}

// NewLeaderboardConfig creates a new Leaderboard configuration struct
func NewLeaderboardConfig(name string, from uint64, to uint64, action string) domain.LeaderboardConfig {
	r := domain.ResetExpression{
		Type: domain.Hourly,
	}
	ce, _ := domain.NewCronExpression(r)
	return domain.LeaderboardConfig{
		Name:            name,
		Function:        domain.Sum,
		ResetExpression: r,
		PrizeTable: domain.LeaderboardPrizeTable{
			Table: []domain.LeaderboardPrize{
				{
					RankFrom: from,
					RankTo:   to,
					Action:   action,
				},
			},
		},
		CronExpression: ce,
	}
}

func NewLeaderboardConfigWithScoreboards(name string, resetType domain.LeaderboardResetType, function domain.LeaderboardFunctionType) domain.LeaderboardConfig {
	r := domain.ResetExpression{
		Type: resetType,
	}
	ce, _ := domain.NewCronExpression(r)
	return domain.LeaderboardConfig{
		Name:            name,
		Function:        function,
		ResetExpression: r,
		PrizeTable: domain.LeaderboardPrizeTable{
			Table: []domain.LeaderboardPrize{
				{
					RankFrom: 1,
					RankTo:   1,
					Action:   "reward 1",
				},
			},
		},
		Scoreboards: []domain.LeaderboardScoreBoardConfig{
			{
				Type:  domain.League,
				Field: "league",
			},
			{
				Type:  domain.Country,
				Field: "country",
			},
		},
		CronExpression: ce,
	}

}
func NewLeaderboardConfigWithFunctionResetWithScoreboards(name string, resetType domain.LeaderboardResetType, function domain.LeaderboardFunctionType) domain.LeaderboardConfig {
	return domain.LeaderboardConfig{
		Name:     name,
		Function: function,
		ResetExpression: domain.ResetExpression{
			Type: resetType,
		},
		PrizeTable: domain.LeaderboardPrizeTable{
			Table: []domain.LeaderboardPrize{
				{
					RankFrom: 1,
					RankTo:   1,
					Action:   "reward 1",
				},
			},
		},
		Scoreboards: []domain.LeaderboardScoreBoardConfig{
			{
				Type:  domain.League,
				Field: "league",
			},
			{
				Type:  domain.Country,
				Field: "country",
			},
		},
	}

}

func NewLeaderboardConfigWithFunctionReset(name string, resetType domain.LeaderboardResetType, function domain.LeaderboardFunctionType) domain.LeaderboardConfig {
	return domain.LeaderboardConfig{
		Name:     name,
		Function: function,
		ResetExpression: domain.ResetExpression{
			Type: resetType,
		},
		PrizeTable: domain.LeaderboardPrizeTable{
			Table: []domain.LeaderboardPrize{
				{
					RankFrom: 1,
					RankTo:   1,
					Action:   "reward 1",
				},
			},
		},
	}

}
func NewMockDefaultDynamoDBSettings(client repository.DynamoDBClient) repository.DynamoDBSettings {
	return repository.DynamoDBSettings{
		Client: client,
		Logger: logging.NewSimpleLogger(),
		Table:  DynamoDBLocalTableName,
	}
}

func NewDefaultDynamoDBSettings() repository.DynamoDBSettings {
	return repository.DynamoDBSettings{
		Client: repository.NewDynamoDBClientFromConfig(repository.DefaultLocalAWSClientConfig()),
		Logger: logging.NewSimpleLogger(),
		Table:  DynamoDBLocalTableName,
	}
}
