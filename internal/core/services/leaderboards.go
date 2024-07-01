package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/posilva/simpleboards/internal/core/ports"
)

// LeaderboardsService ...
type LeaderboardsService struct {
	repository    ports.Repository
	scoreboard    ports.Scoreboard
	configuration ports.Provider[domain.LeaderboardsConfigMap]
}

// NewLeaderboardsService creates a new leaderboards service
func NewLeaderboardsService(
	repo ports.Repository,
	scoreboard ports.Scoreboard,
	configProvider ports.ConfigProvider,
) *LeaderboardsService {
	return &LeaderboardsService{
		repository:    repo,
		scoreboard:    scoreboard,
		configuration: configProvider,
	}
}

// GetConfig returns the config for a given leaderboard identified by name
func (s *LeaderboardsService) GetConfig(name string) (domain.LeaderboardConfig, error) {
	configMap, err := s.configuration.Provide()
	if err != nil {
		return domain.LeaderboardConfig{}, fmt.Errorf("failed to provide configuration: %v", err)
	}
	config, ok := configMap[name]
	if !ok {
		return domain.LeaderboardConfig{}, fmt.Errorf("leaderboard config not found: %v", name)
	}
	return config, nil
}

// ReportScore ...
func (s *LeaderboardsService) ReportScore(entryID string, name string, score float64) (domain.ReportScoreOutput, error) {
	return s.ReportScoreWithMetadata(entryID, name, score, nil)
}

// ReportScoreWithMetadata ...
func (s *LeaderboardsService) ReportScoreWithMetadata(entryID string, name string, score float64, meta domain.Metadata) (domain.ReportScoreOutput, error) {

	// ReportScore  register a new score to a given entry on a leaderboard
	config, err := s.GetConfig(name)
	if err != nil {
		return domain.ReportScoreOutput{}, fmt.Errorf("failed to fetch configs: %v", err)
	}

	leaderboard, epoch, err := GetLeaderboardNameWithEpoch(name, config.CronExpression)

	if err != nil {
		return domain.ReportScoreOutput{}, fmt.Errorf("failed to generate name from configs: %v", err)
	}

	lbFn := s.applyFunction(entryID, leaderboard, score, config.Function, meta)
	v, err := lbFn()
	if err != nil {
		return domain.ReportScoreOutput{}, fmt.Errorf("failed to apply functoin to the  score: %v", err)
	}

	if v.Done {
		// Global scoreboard
		err = s.scoreboard.AddScore(entryID, leaderboard, v.Score)
		if err != nil {
			return domain.ReportScoreOutput{}, fmt.Errorf("failed to add score to scoreboard: %v", err)
		}
		// add to other scoreboards
		if config.Scoreboards != nil && len(config.Scoreboards) > 0 {
			for _, sb := range config.Scoreboards {
				// TODO: we may enforce to exist the config fields in the meta for correctness
				lb := s.sbNameFromType(name, epoch, sb, meta[sb.Field])
				err = s.scoreboard.AddScore(entryID, lb, v.Score)
				if err != nil {
					return domain.ReportScoreOutput{}, fmt.Errorf("failed to add score to scoreboard: %v", err)
				}
			}
		}
	}

	return domain.ReportScoreOutput{Update: v, Epoch: epoch}, nil
}

func (s *LeaderboardsService) sbNameFromType(lb string, epoch int64, sb domain.LeaderboardScoreBoardConfig, value string) string {
	name := fmt.Sprintf("%s::%d", lb, epoch)
	switch sb.Type {
	case domain.League:
		name = fmt.Sprintf("%s::%s::%s::%d", lb, "league", value, epoch)
	case domain.Country:
		name = fmt.Sprintf("%s::%s::%s::%d", lb, "country", value, epoch)
	}
	return strings.ToLower(name)
}

func (s *LeaderboardsService) applyFunction(entryID string, leaderboard string, score float64, configFunction domain.LeaderboardFunctionType, meta domain.Metadata) func() (domain.ScoreUpdate, error) {
	lbFn := func() (domain.ScoreUpdate, error) {
		return s.repository.AddWithMetadata(entryID, leaderboard, score, meta)
	}

	switch configFunction {
	case domain.Max:
		lbFn = func() (domain.ScoreUpdate, error) {
			return s.repository.MaxWithMetadata(entryID, leaderboard, score, meta)
		}
	case domain.Min:
		lbFn = func() (domain.ScoreUpdate, error) {
			return s.repository.MinWithMetadata(entryID, leaderboard, score, meta)
		}
	case domain.Last:
		lbFn = func() (domain.ScoreUpdate, error) {
			return s.repository.LastWithMetadata(entryID, leaderboard, score, meta)
		}
	}
	return lbFn

}

// ListScoresWithMetadata returns a list of scores from leaderboards with metadata
func (s *LeaderboardsService) ListScoresWithMetadata(name string, meta domain.Metadata) ([]domain.LeaderboardScores, int64, error) {
	config, err := s.GetConfig(name)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch configs: %v", err)
	}
	leaderboard, epoch, err := GetLeaderboardNameWithEpoch(name, config.CronExpression)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate name from configs: %v", err)
	}

	scores, err := s.scoreboard.Get(leaderboard)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch scores: %v", err)
	}
	var allLeaderboardScores []domain.LeaderboardScores

	resultScores := domain.LeaderboardScores{}
	resultScores.Name = leaderboard
	for _, score := range scores {
		resultScores.Scores = append(resultScores.Scores, domain.LeaderboardEntry{
			EntryID: score.EntryID,
			Score:   score.Score,
			Rank:    score.Rank,
		})
	}
	allLeaderboardScores = append(allLeaderboardScores, resultScores)

	if config.Scoreboards != nil && len(config.Scoreboards) > 0 {
		for _, sb := range config.Scoreboards {
			lb := s.sbNameFromType(name, epoch, sb, meta[sb.Field])
			scores, err := s.scoreboard.Get(lb)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to fetch scores for scoreboard: %v: %v", lb, err)
			}
			resultScores := domain.LeaderboardScores{}
			resultScores.Name = lb
			for _, score := range scores {
				resultScores.Scores = append(resultScores.Scores, domain.LeaderboardEntry{
					EntryID: score.EntryID,
					Score:   score.Score,
					Rank:    score.Rank,
				})
			}
			allLeaderboardScores = append(allLeaderboardScores, resultScores)
		}
	}
	return allLeaderboardScores, epoch, nil

}

// ListScores returns a list of scores from leaderboards
func (s *LeaderboardsService) ListScores(name string) ([]domain.LeaderboardScores, int64, error) {
	return s.ListScoresWithMetadata(name, nil)
}

// GetResults returns a list of scores from leaderboards
func (s *LeaderboardsService) GetResults(name string, epoch int64) ([]domain.LeaderboardScores, error) {
	return s.GetResultsWithMetadata(name, epoch, nil)
}

// GetResultsWithMetadata returns a list of scores from leaderboards
func (s *LeaderboardsService) GetResultsWithMetadata(name string, epoch int64, meta domain.Metadata) ([]domain.LeaderboardScores, error) {
	config, err := s.GetConfig(name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch configs: %v", err)
	}
	allResults := []domain.LeaderboardScores{}

	leaderboard := getNameWithEpoch(name, epoch)
	scores, err := s.scoreboard.Get(leaderboard)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scores: %v", err)
	}

	// TODO: should get the score boards
	resultScores := domain.LeaderboardScores{}
	resultScores.Name = leaderboard
	for _, score := range scores {
		resultScores.Scores = append(resultScores.Scores, domain.LeaderboardEntry{
			EntryID: score.EntryID,
			Score:   score.Score,
			Rank:    score.Rank,
		})
	}
	allResults = append(allResults, resultScores)

	if config.Scoreboards != nil && len(config.Scoreboards) > 0 {
		for _, sb := range config.Scoreboards {
			leaderboard = s.sbNameFromType(name, epoch, sb, meta[sb.Field])
			scores, err := s.scoreboard.Get(leaderboard)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch scores: %v", err)
			}

			resultScores := domain.LeaderboardScores{}
			resultScores.Name = leaderboard
			for _, score := range scores {
				resultScores.Scores = append(resultScores.Scores, domain.LeaderboardEntry{
					EntryID: score.EntryID,
					Score:   score.Score,
					Rank:    score.Rank,
				})
			}

			allResults = append(allResults, resultScores)
		}

	}

	return allResults, nil
}

func GetLeaderboardNameWithEpoch(name string, reset domain.CronExpression) (string, int64, error) {
	epoch := reset.GetEpochFromReferenceUnixTimestamp(time.Now().Unix())
	return strings.ToLower(getNameWithEpoch(name, epoch)), epoch, nil
}

func getNameWithEpoch(name string, epoch int64) string {
	return strings.ToLower(fmt.Sprintf("%s::%d", name, epoch))
}
