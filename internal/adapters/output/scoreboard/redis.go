package scoreboard

import (
	"context"
	"fmt"

	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/redis/rueidis"
)

type RedisScoreboard struct {
	options RedisScoreboardOptions
	client  rueidis.Client
}

type RedisScoreboardOptions struct {
	BatchSize int `json:"batch_size"`
}

// DefaultRedisScoreboardOptions returns the default optoins for redis cache
func DefaultRedisScoreboardOptions() RedisScoreboardOptions {
	return RedisScoreboardOptions{
		BatchSize: 50,
	}
}

// NewRedisScoreboard creates an instance of Redis Cache
func NewRedisScoreboard(address string) (*RedisScoreboard, error) {
	opts := rueidis.ClientOption{
		InitAddress: []string{
			address,
		},
	}
	c, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis host '%v': %v ", address, err)
	}
	return &RedisScoreboard{
		client:  c,
		options: DefaultRedisScoreboardOptions(),
	}, nil
}

// Get returns the list of results with batchsize
func (c *RedisScoreboard) Get(name string) ([]domain.ScoreboardResult, error) {
	return c.GetTopN(name, int64(c.options.BatchSize))
}

// GetTopN ...
func (c *RedisScoreboard) GetTopN(name string, n int64) ([]domain.ScoreboardResult, error) {
	cmd := c.client.B().Zrevrange().Key(name).Start(0).Stop(n).Withscores().Build()
	m, err := c.client.Do(context.Background(), cmd).AsZScores()
	if err != nil {
		return nil, err
	}

	results := []domain.ScoreboardResult{}
	for i, r := range m {
		results = append(results, domain.ScoreboardResult{
			EntryID: r.Member,
			Score:   r.Score,
			Rank:    int64(i + 1),
		})
	}
	return results, nil
}

// AddScore ...
func (c *RedisScoreboard) AddScore(entryID string, nameWithEpoch string, value float64) error {
	cmd := c.client.B().Zadd().Key(nameWithEpoch).ScoreMember().ScoreMember(value, entryID).Build()
	err := c.client.Do(context.Background(), cmd).Error()
	return err
}

// TODO: check the return of the functtion to match the Rank type in the result

// GetRank ...
func (c *RedisScoreboard) GetRank(nameWithEpoch string, entryID string) (uint64, error) {
	cmd := c.client.B().Zrevrank().Key(nameWithEpoch).Member(entryID).Build()
	score, err := c.client.Do(context.Background(), cmd).AsInt64()
	if err != nil {
		return 0, fmt.Errorf("failed to get rank: %v", err)
	}
	return uint64(score) + 1, nil
}
