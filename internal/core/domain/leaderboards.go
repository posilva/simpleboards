package domain

// Metadata type
type Metadata map[string]string

// LeaderboardsConfigMap is a map for leaderboards configuration
type LeaderboardsConfigMap = map[string]LeaderboardConfig

// LeaderboardFunctionType enum for leaderboards function
type LeaderboardFunctionType int

const (
	// Last function that saves always the last value stored
	Last LeaderboardFunctionType = iota
	// Max function that saves always the max value stored
	Max
	// Min function that saves always the min value stored
	Min
	// Sum function that accumulated value stored
	Sum
)

// LeaderboardResetType enum for leaderboards reset type
type LeaderboardResetType int

const (
	Manually LeaderboardResetType = iota
	Hourly
	Daily
	Weekly
	Monthly
	Custom
)

type LeaderboardPrizeTable struct {
	Table []LeaderboardPrize `json:"table"`
}

func (t LeaderboardPrizeTable) Validate() error {
	return nil
}

// LeaderboardPrize holds data for configuration of rewards
type LeaderboardPrize struct {
	RankFrom uint64 `json:"rank_from"`
	RankTo   uint64 `json:"rank_to"`
	Action   string `json:"action"`
}

type ScoreboardType = int

const (
	League ScoreboardType = iota
	Country
)

type LeaderboardScoreBoardConfig struct {
	Type  ScoreboardType `json:"type"`
	Field string         `json:"field"`
}

type ResetExpression struct {
	Type           LeaderboardResetType `json:"reset_type"`
	CronExpression string               `json:"cron, omitempty"`
}

// LeaderboardConfig holds information of a Leaderboard instance
type LeaderboardConfig struct {
	Name            string                        `json:"name"`
	Function        LeaderboardFunctionType       `json:"function"`
	ResetExpression ResetExpression               `json:"reset`
	PrizeTable      LeaderboardPrizeTable         `json:"prizes_table"`
	Scoreboards     []LeaderboardScoreBoardConfig `json:"scoreboards"`
	CronExpression  CronExpression                `json:"-"`
}

// TODO: add a field to represent the metadata to show in the UI
// This may be Avatar, Username, Group Badge etc

// LeaderboardEntry entry data
type LeaderboardEntry struct {
	Metadata string  `json:"metadata"`
	EntryID  string  `json:"entry_id"`
	Score    float64 `json:"score"`
	Rank     int64   `json:"rank"`
}

// LeaderboardScores leaderboard score
type LeaderboardScores struct {
	Name   string             `json:"name"`
	Scores []LeaderboardEntry `json:"scores"`
}

type ScoreUpdate struct {
	Score    float64           `json:"score,omitempty"`
	Done     bool              `json:"done,omitempty"`
	Counter  uint64            `json:"counter,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type ReportScoreOutput struct {
	Update ScoreUpdate
	Epoch  int64
}
