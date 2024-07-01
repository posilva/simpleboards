package configprovider

import (
	"errors"
	"fmt"
	"sync"

	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/posilva/simpleboards/internal/core/ports"
	"github.com/posilva/simpleboards/internal/core/services"
)

// TODO: We can generalise the refresh component
const (
	refreshIntervalSecs = 5
)

type DynamoConfigProvider[T domain.LeaderboardsConfigMap] struct {
	lock          sync.RWMutex
	configGetter  ports.ConfigGetter
	currentConfig domain.LeaderboardsConfigMap
	err           error
	scheduler     *services.Scheduler
	logger        ports.Logger
}

func NewDynamoConfigProvider(configGetter ports.ConfigGetter, logger ports.Logger) *DynamoConfigProvider[domain.LeaderboardsConfigMap] {
	cp := DynamoConfigProvider[domain.LeaderboardsConfigMap]{
		err:           errors.New("config not initialized"),
		currentConfig: nil,
		logger:        logger,
		configGetter:  configGetter,
	}

	cp.scheduler = services.NewScheduler(refreshIntervalSecs, cp.Refresh)
	return &cp
}

func (cp *DynamoConfigProvider[T]) Refresh() {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	cfgMap, err := cp.configGetter.GetConfig()
	cp.logger.Debug("Refreshing configuration: %v", cfgMap)
	if err != nil {
		cp.logger.Error("failed to get configuration: %v", err)
		return
	}

	for name, config := range cfgMap {
		ce, err := domain.NewCronExpression(config.ResetExpression)
		if err != nil {
			cp.logger.Error("failed to update cron expression (%v) in configuration '%v': %v", config.ResetExpression, name, err)
			return
		}
		config.CronExpression = ce
		cfgMap[name] = config
	}

	cp.currentConfig = cfgMap
}

// Provide configurations for the leaderboard:s
func (cp *DynamoConfigProvider[T]) Provide() (domain.LeaderboardsConfigMap, error) {
	cp.lock.RLock()
	defer cp.lock.RUnlock()

	if cp.currentConfig != nil {
		return cp.currentConfig, nil
	}
	return cp.currentConfig, fmt.Errorf("configuration is not valid: %v", cp.err)
}
