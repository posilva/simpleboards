package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
	"testing"
)

var (
	baseURL       = "localhost:8808"
	defaultScheme = "http"
)

type listScoresResponse struct {
	Scores []struct {
		Name   string `json:"name"`
		Scores []struct {
			Metadata string  `json:"metadata"`
			Entry    string  `json:"entry_id"`
			Score    float64 `json:"score"`
			Rank     int     `json:"rank"`
		} `json:"scores"`
	} `json:"scores"`
	Epoch int `json:"epoch"`
}

type putScoreResponse struct {
	Score float64 `json:"new_score"`
}
type putScoreRequest struct {
	Entry string  `json:"entry"`
	Score float64 `json:"score"`
}

type E2ETestSuite struct {
	BaseTestSuite
}

func (suite *E2ETestSuite) SetupSuite() {
	setup(&suite.BaseTestSuite)
	baseURL = suite.ServiceEndpoint
}

func (suite *E2ETestSuite) TearDownSuite() {
	fmt.Println("Running teardown suite")
	teardown(&suite.BaseTestSuite)
}

func (suite *E2ETestSuite) TestSingleEntrySumMultipleScoreTest() {
	entryID := testutil.NewID()

	lbName := defaultLbNameSum
	resp, err := reportScore(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScore(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(float64(25), resp.Score)
	_, err = listScores(lbName)
	suite.NoError(err)
}

func (suite *E2ETestSuite) TestSingleEntryMaxMultipleScoreTest() {
	lbName := defaultLbNameMax
	entryID := testutil.NewID()
	resp, err := reportScore(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScore(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(float64(15), resp.Score)
	_, err = listScores(lbName)
	suite.NoError(err)
}

func (suite *E2ETestSuite) TestSingleEntryMinMultipleScoreTest() {
	lbName := defaultLbNameMin
	entryID := testutil.NewID()
	resp, err := reportScore(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScore(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(float64(0), resp.Score)
	_, err = listScores(lbName)
	suite.NoError(err)
}

func (suite *E2ETestSuite) TestSingleEntryLastMultipleScoreTest() {
	lbName := defaultLbNameLast
	entryID := testutil.NewID()
	resp, err := reportScore(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScore(lbName, entryID, 15)
	suite.NoError(err)
	resp, err = reportScore(lbName, entryID, 5)
	suite.NoError(err)
	suite.Equal(float64(5), resp.Score)
	_, err = listScores(lbName)
	suite.NoError(err)
}

func (suite *E2ETestSuite) TestMultipleSumEntry() {
	entryID := testutil.NewID()
	entryID2 := testutil.NewID()
	lbName := defaultLbNameSum
	resp, err := reportScore(lbName, entryID, 10)

	suite.NoError(err)
	suite.Equal(float64(10), resp.Score)

	resp, err = reportScore(lbName, entryID2, 15)
	suite.NoError(err)
	suite.Equal(float64(15), resp.Score)

	list, err := listScores(lbName)
	suite.NoError(err)
	suite.Len(list.Scores[0].Scores, 2)
	suite.Equal(list.Scores[0].Scores[0].Entry, entryID2)
	suite.Equal(list.Scores[0].Scores[0].Rank, 1)
	suite.Equal(list.Scores[0].Scores[1].Entry, entryID)
	suite.Equal(list.Scores[0].Scores[1].Rank, 2)
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
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

func listScores(lbname string) (response listScoresResponse, err error) {
	path := fmt.Sprintf("/api/v1/scores/%s", lbname)
	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToJSON(&response).
		Fetch(context.Background())

	return response, err

}
