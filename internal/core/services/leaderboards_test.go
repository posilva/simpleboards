package services

import (
	"fmt"
	"strings"
	"testing"

	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/posilva/simpleboards/internal/core/ports/mocks"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TODO: missing implementation
func TestGetConfig(t *testing.T) {
	assert.Nil(t, nil)
}

func TestReportScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()
	value := 100.0

	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)

	configProvider := defaultConfigProviderMock(ctrl, lbName)
	nameEpoch, _, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo.EXPECT().AddWithMetadata(entryID, nameEpoch, value, nil).Return(domain.ScoreUpdate{Score: value, Done: true}, nil)
	scoreboard.EXPECT().AddScore(entryID, nameEpoch, value).Return(nil)
	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, err := lbSrv.ReportScore(entryID, lbName, value)
	assert.NoError(t, err)
	assert.Nil(t, nil)
	assert.Equal(t, value, v.Update.Score)
}

func TestReportScoreWithScoreboards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	entryID := testutil.NewID()
	value := 100.0

	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)

	configProvider := defaultConfigProviderMockWithScoreboards(ctrl, lbName)
	_, _, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo.EXPECT().AddWithMetadata(entryID, gomock.Any(), value, nil).Return(domain.ScoreUpdate{Score: value, Done: true}, nil).AnyTimes()
	scoreboard.EXPECT().AddScore(entryID, gomock.Any(), value).Return(nil).AnyTimes()
	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, err := lbSrv.ReportScore(entryID, lbName, value)
	assert.NoError(t, err)
	assert.Nil(t, nil)
	assert.Equal(t, value, v.Update.Score)
}

func TestListScores(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	nameEpoch, _, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)
	configProvider := defaultConfigProviderMock(ctrl, lbName)

	scoreboard.EXPECT().Get(nameEpoch).Return([]domain.ScoreboardResult{}, nil)

	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, _, err := lbSrv.ListScores(lbName)
	assert.NoError(t, err)
	assert.Len(t, v, 1)
	assert.True(t, strings.Contains(v[0].Name, lbName))
}

func TestListScoresWithMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	nameEpoch, _, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)
	configProvider := defaultConfigProviderMock(ctrl, lbName)

	scoreboard.EXPECT().Get(nameEpoch).Return([]domain.ScoreboardResult{}, nil)

	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, _, err := lbSrv.ListScoresWithMetadata(lbName, domain.Metadata{
		"country": "PT",
		"league":  "gold",
	})
	assert.NoError(t, err)
	assert.Len(t, v, 1)
	assert.True(t, strings.Contains(v[0].Name, lbName))
}

func TestGetResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	nameEpoch, epoch, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)
	configProvider := defaultConfigProviderMock(ctrl, lbName)

	scoreboard.EXPECT().Get(nameEpoch).Return([]domain.ScoreboardResult{}, nil)

	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, err := lbSrv.GetResults(lbName, epoch)
	fmt.Println(v)
	assert.NoError(t, err)
	assert.Len(t, v, 1)
	assert.True(t, strings.Contains(v[0].Name, lbName))
}
func TestGetResultsWithMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	nameEpoch, epoch, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)
	configProvider := defaultConfigProviderMock(ctrl, lbName)

	scoreboard.EXPECT().Get(nameEpoch).Return([]domain.ScoreboardResult{}, nil)

	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, err := lbSrv.GetResultsWithMetadata(lbName, epoch, domain.Metadata{
		"country": "PT",
		"league":  "gold",
	})
	fmt.Println(v)
	assert.NoError(t, err)
	assert.Len(t, v, 1)
	assert.True(t, strings.Contains(v[0].Name, lbName))
}

func TestGetResultsWithScoreboards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lbName := testutil.NewUnique(testutil.Name(t))
	_, epoch, err := GetLeaderboardNameWithEpoch(lbName, domain.Hourly)
	assert.NoError(t, err)
	repo := mocks.NewMockRepository(ctrl)
	scoreboard := mocks.NewMockScoreboard(ctrl)
	configProvider := defaultConfigProviderMockWithScoreboards(ctrl, lbName)

	scoreboard.EXPECT().Get(gomock.Any()).Return([]domain.ScoreboardResult{}, nil).AnyTimes()

	lbSrv := NewLeaderboardsService(repo, scoreboard, configProvider)

	v, err := lbSrv.GetResults(lbName, epoch)
	assert.NoError(t, err)
	assert.Len(t, v, 3)
	assert.True(t, strings.Contains(v[0].Name, lbName))
}

func defaultConfigProviderMock(ctrl *gomock.Controller, lbName string) *mocks.MockConfigProvider {
	cp := mocks.NewMockConfigProvider(ctrl)

	configMap := make(map[string]domain.LeaderboardConfig)
	configMap[lbName] = testutil.NewLeaderboardConfig(lbName, 1, 1, "reward_test")
	cp.EXPECT().Provide().Return(configMap, nil)
	return cp
}

func defaultConfigProviderMockWithScoreboards(ctrl *gomock.Controller, lbName string) *mocks.MockConfigProvider {
	cp := mocks.NewMockConfigProvider(ctrl)

	configMap := make(map[string]domain.LeaderboardConfig)
	configMap[lbName] = testutil.NewLeaderboardConfigWithScoreboards(lbName, domain.Hourly, domain.Sum)
	cp.EXPECT().Provide().Return(configMap, nil)
	return cp
}
