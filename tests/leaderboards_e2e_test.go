package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/posilva/simpleboards/internal/adapters/output/repository"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	baseURL       = "localhost:8808"
	defaultScheme = "http"
)

var defaultLbName = "integration_lb_tests"

type putScoreResponse struct {
	Score float64 `json:"new_score"`
}
type putScoreRequest struct {
	Entry string  `json:"entry"`
	Score float64 `json:"score"`
}

func configLeaderboard() {
	settings := testutil.NewDefaultDynamoDBSettings()
	repo, err := repository.NewDynamoDBRepository(settings)
	if err != nil {
		panic(err)
	}
	lbName := defaultLbName
	lbConfig := testutil.NewLeaderboardConfig(lbName, 1, 1, "reward 1")
	err = repo.Update(lbName, lbConfig)
	_, err = repo.GetConfig()
	if err != nil {
		panic(err)
	}
}

func TestConfigLeaderboard(t *testing.T) {
	assert.Equal(t, true, true)
}

func TestHTTPReportScore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	configLeaderboard()
	resp, err := reportScore(defaultLbName, "123", 10)
	if err != nil {
		log.Fatalf("failed to execute GetAuthConfig: %v", err)
	}

	assert.Equal(t, float64(10), resp.Score)
}

func reportScore(ldbName string, entry string, score float64) (response putScoreResponse, err error) {
	path := fmt.Sprintf("/api/v1/score/%s", ldbName)

	req := putScoreRequest{
		Entry: entry,
		Score: score,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return response, err
	}
	err = requests.
		URL(path).
		Put().
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		BodyReader(strings.NewReader(string(data))).
		ToJSON(&response).
		Fetch(context.Background())
	return response, err
}
