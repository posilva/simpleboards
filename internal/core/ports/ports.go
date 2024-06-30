package ports

import (
	"time"

	"github.com/posilva/simpleboards/internal/core/domain"
)

// Repository defines the interface to handle with
type Repository interface {
	Add(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error)
	Max(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error)
	Min(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error)
	Last(entry string, leaderboard string, value float64) (domain.ScoreUpdate, error)
	AddWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error)
	MaxWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error)
	MinWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error)
	LastWithMetadata(entry string, leaderboard string, value float64, meta domain.Metadata) (domain.ScoreUpdate, error)
}

// Logger defines a basic logger interface
type Logger interface {
	Debug(msg string, v ...interface{}) error
	Info(msg string, v ...interface{}) error
	Error(msg string, v ...interface{}) error
}

// LeaderboardsService defines the leaderboard service interface
type LeaderboardsService interface {
	GetConfig(name string) (domain.LeaderboardConfig, error)
	ReportScore(entryID string, name string, value float64) (domain.ReportScoreOutput, error)
	ReportScoreWithMetadata(entryID string, name string, value float64, meta domain.Metadata) (domain.ReportScoreOutput, error)
	ListScores(name string) ([]domain.LeaderboardScores, int64, error)
	ListScoresWithMetadata(name string, meta domain.Metadata) ([]domain.LeaderboardScores, int64, error)
	// TODO: we may have a dedicated data type to return in this call
	GetResults(name string, epoch int64) ([]domain.LeaderboardScores, error)
	GetResultsWithMetadata(name string, epoch int64, meta domain.Metadata) ([]domain.LeaderboardScores, error)
}

// Scoreboard ...
type Scoreboard interface {
	Get(name string) ([]domain.ScoreboardResult, error)
	GetTopN(name string, n int64) ([]domain.ScoreboardResult, error)
	AddScore(entryID string, name string, value float64) error
	GetRank(entryID string, name string) (uint64, error)
}

// Provider generic interface
type Provider[T any] interface {
	Provide() (T, error)
}

// ConfigProvider ...
type ConfigProvider interface {
	Provider[domain.LeaderboardsConfigMap]
	Refresh()
}

// ConfigGetter defines the interface to retrieve configs
type ConfigGetter interface {
	GetConfig() (domain.LeaderboardsConfigMap, error)
}

// ResetLocker defines the interface to lock during the Reset
type ResetLocker interface {
	ResetLock(leaderboar string, epoch int, duration time.Duration) (bool, error)
}

// TelemetryReporter defines the interface to report metrics
type TelemetryReporter interface {
	SetDefaultTags(tags map[string]string)
	ReportGauge(name string, value float64, tags map[string]string)
	ReportCounter(name string, value float64, tags map[string]string)
	ReportHistogram(name string, value float64, tags map[string]string)
	ReportSummary(name string, value float64, tags map[string]string)
}
