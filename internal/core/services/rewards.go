package services

import (
	"github.com/posilva/simpleboards/internal/core/ports"
)

// RewardWatcher represents the data
type RewardWatcher struct {
	logger     ports.Logger
	scheduler  *Scheduler
	config     ports.ConfigProvider
	repository ports.Repository
}

// NewRewardWatcher creates a new reward watcher
func NewRewardWatcher(intervalSecs int, repo ports.Repository, config ports.ConfigProvider, logger ports.Logger) *RewardWatcher {
	r := &RewardWatcher{
		logger:     logger,
		config:     config,
		repository: repo,
	}

	r.scheduler = NewScheduler(intervalSecs, r.Check)
	return r
}

// Check if is time for reward participants
func (r *RewardWatcher) Check() {
	// fetch configuration if exists already

	cfgs, err := r.config.Provide()
	if err != nil {
		r.logger.Error("failed to read configs: %v", err)
		return
	}

	// for each leaderboard check if it's time to reset
	for lb, cfg := range cfgs {
		_ = lb

		// calculate the epoch
		lbNameWithEpoch, _, err := GetLeaderboardNameWithEpoch(lb, cfg.Reset)
		if err != nil {
			r.logger.Error("failed to calculate epoch: %v", err)
		}
		_ = lbNameWithEpoch
	}
}
