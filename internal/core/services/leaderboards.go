package services

import (
	"fmt"
	"time"

	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/posilva/simpleboards/internal/core/ports"
)

// LeaderboardsService
type LeaderboardsService struct {
	repository    ports.Repository
	scoreboard    ports.Scoreboard
	configuration ports.Provider[domain.LeaderboardsConfigMap]
}

// NewLeaderboardsServige creates a new leaderboards service
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

// ReportScore  register a new score to a given entry on a leaderboard
func (s *LeaderboardsService) ReportScore(entryID string, name string, score float64) (domain.ReportScoreOutput, error) {
	config, err := s.GetConfig(name)
	if err != nil {
		return domain.ReportScoreOutput{}, fmt.Errorf("failed to fetch configs: %v", err)
	}

	leaderboard, epoch, err := GetLeaderboardNameWithEpoch(name, config.Reset)
	if err != nil {
		return domain.ReportScoreOutput{}, fmt.Errorf("failed to generate name from configs: %v", err)
	}

	lbFn := func() (domain.ScoreUpdate, error) {
		return s.repository.Add(entryID, leaderboard, score)
	}
	switch config.Function {
	case domain.Max:
		lbFn = func() (domain.ScoreUpdate, error) {
			return s.repository.Max(entryID, leaderboard, score)
		}
	case domain.Min:
		lbFn = func() (domain.ScoreUpdate, error) {
			return s.repository.Min(entryID, leaderboard, score)
		}
	case domain.Last:
		lbFn = func() (domain.ScoreUpdate, error) {
			return s.repository.Last(entryID, leaderboard, score)
		}
	}
	v, err := lbFn()
	if err != nil {
		return domain.ReportScoreOutput{}, fmt.Errorf("failed to apply functoin to the  score: %v", err)
	}
	if v.Done {
		err = s.scoreboard.AddScore(entryID, leaderboard, v.Score)
		if err != nil {
			return domain.ReportScoreOutput{}, fmt.Errorf("failed to add score to scoreboard: %v", err)
		}
	}

	return domain.ReportScoreOutput{Update: v, Epoch: epoch}, nil
}

// ListScores returns a list of scores from leaderboards
func (s *LeaderboardsService) ListScores(name string) ([]domain.LeaderboardScores, int64, error) {
	config, err := s.GetConfig(name)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch configs: %v", err)
	}

	leaderboard, epoch, err := GetLeaderboardNameWithEpoch(name, config.Reset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate name from configs: %v", err)
	}

	scores, err := s.scoreboard.Get(leaderboard)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch scores: %v", err)
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

	return []domain.LeaderboardScores{resultScores}, epoch, nil
}

func (s *LeaderboardsService) GetResults(name string, epoch int64) ([]domain.LeaderboardScores, error) {
	config, err := s.GetConfig(name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch configs: %v", err)
	}

	leaderboard := getNameWithEpoch(name, epoch)

	scores, err := s.scoreboard.Get(leaderboard)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scores: %v", err)
	}
	// TODO: get the prize table ranks
	_ = config

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
	return []domain.LeaderboardScores{resultScores}, nil
}

// TODO: Create the name from the current timestamp and configuration settings
func GetLeaderboardNameWithEpoch(name string, resetType domain.LeaderboardResetType) (string, int64, error) {
	epoch, err := CalculateEpoch(resetType, time.Now().Unix())
	if err != nil {
		return "", 0, err
	}
	return getNameWithEpoch(name, epoch), epoch, nil
}

func getNameWithEpoch(name string, epoch int64) string {
	return fmt.Sprintf("%s::%d", name, epoch)
}

func CalculateEpoch(resetType domain.LeaderboardResetType, posixTs int64) (int64, error) {
	hour := posixTs / 60 / 60
	day := hour / 24
	week := day / 7
	month := day / 30

	switch resetType {
	case domain.Hourly:
		return hour, nil
	case domain.Daily:
		return day, nil
	case domain.Weekly:
		return week, nil
	case domain.Monthly:
		return month, nil
	}
	return 0, fmt.Errorf("unknown reset type: %v", resetType)
}
