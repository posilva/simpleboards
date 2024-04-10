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
	refreshIntervalSecs = 60
)

type SimpleConfigProvider[T domain.LeaderboardsConfigMap] struct {
	lock          sync.RWMutex
	configGetter  ports.ConfigGetter
	currentConfig domain.LeaderboardsConfigMap
	err           error
	scheduler     *services.Scheduler
	logger        ports.Logger
}

func NewSimpleConfigProvider(configGetter ports.ConfigGetter, logger ports.Logger) *SimpleConfigProvider[domain.LeaderboardsConfigMap] {
	cp := SimpleConfigProvider[domain.LeaderboardsConfigMap]{
		err:           errors.New("config not initiatiazed"),
		currentConfig: nil,
		logger:        logger,
		configGetter:  configGetter,
	}

	cp.scheduler = services.NewScheduler(refreshIntervalSecs, cp.Refresh)
	return &cp
}

func (cp *SimpleConfigProvider[T]) Refresh() {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	cfgMap, err := cp.configGetter.GetConfig()
	if err != nil {
		cp.logger.Error("failed to get configuration: %v", err)
		return
	}
	cp.currentConfig = cfgMap
}

// Provide configurations for the leaderboard:s
func (cp *SimpleConfigProvider[T]) Provide() (domain.LeaderboardsConfigMap, error) {
	cp.lock.RLock()
	defer cp.lock.RUnlock()

	if cp.currentConfig != nil {
		return cp.currentConfig, nil
	}
	return cp.currentConfig, fmt.Errorf("configuration is not valid: %v", cp.err)
}
