package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/posilva/simpleboards/internal/testutil"
	"github.com/stretchr/testify/suite"
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
	Done  bool    `json:"done"`
	Epoch int     `json:"epoch"`
	Count int     `json:"count"`
}

type putScoreRequest struct {
	Entry    string          `json:"entry"`
	Score    float64         `json:"score"`
	Metadata domain.Metadata `json:"metadata"`
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

func (suite *E2ETestSuite) TestSingleEntrySumMultipleScoreWithMetadataTest() {
	entryID := testutil.NewID()

	lbName := defaultLbNameSum
	resp, err := reportScoreWithMetadata(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScoreWithMetadata(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(float64(25), resp.Score)
	_, err = listScores(lbName)
	suite.NoError(err)
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

func (suite *E2ETestSuite) TestSingleEntryMaxMultipleScoreWithMetadataTest() {
	lbName := defaultLbNameMax
	entryID := testutil.NewID()
	resp, err := reportScoreWithMetadata(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScoreWithMetadata(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(float64(15), resp.Score)
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

func (suite *E2ETestSuite) TestSingleEntryMinMultipleScoreWithMetadataTest() {
	lbName := defaultLbNameMin
	entryID := testutil.NewID()
	resp, err := reportScoreWithMetadata(lbName, entryID, 10)
	suite.NoError(err)
	suite.Equal(int(1), resp.Count)
	resp, err = reportScoreWithMetadata(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(false, resp.Done)
	suite.Equal(int(0), resp.Count)
	suite.Equal(float64(0), resp.Score)
	resp, err = reportScoreWithMetadata(lbName, entryID, 5)
	suite.NoError(err)
	suite.Equal(float64(5), resp.Score)

	suite.Equal(int(2), resp.Count)
	s, err := listScoresWithMetadata(lbName, metadataDefault)
	suite.Len(s.Scores, 3) // this will return all scoreboards
	fmt.Println("List scores result", s)
	suite.NoError(err)
}

func (suite *E2ETestSuite) TestSingleEntryMinMultipleScoreWithMetadataWrongCountryTest() {
	lbName := defaultLbNameMin
	entryID := testutil.NewID()
	resp, err := reportScoreWithMetadata(lbName, entryID, 10)
	suite.NoError(err)
	suite.Equal(int(1), resp.Count)
	// in this test the country is different from the existing so it will not be accepted once already exists
	resp, err = reportScoreWithMetadataCountryLeage(lbName, entryID, 5, "UK", "gold")
	suite.NoError(err)
	suite.Equal(false, resp.Done)
	suite.Equal(float64(0), resp.Score)
	suite.Equal(int(0), resp.Count)
	s, err := listScoresWithMetadata(lbName, metadataDefault)
	suite.Len(s.Scores, 3) // this will return all scoreboards
	fmt.Println("List scores result", s)
	suite.NoError(err)
}

func (suite *E2ETestSuite) TestSingleEntryMinMultipleScoreTest() {
	lbName := defaultLbNameMin
	entryID := testutil.NewID()
	resp, err := reportScore(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScore(lbName, entryID, 15)
	suite.NoError(err)
	suite.Equal(false, resp.Done)
	suite.Equal(float64(0), resp.Score)
	_, err = listScores(lbName)
	suite.NoError(err)
}
func (suite *E2ETestSuite) TestSingleEntryLastMultipleScoreWithMetadataTest() {
	lbName := defaultLbNameLast
	entryID := testutil.NewID()
	resp, err := reportScoreWithMetadata(lbName, entryID, 10)
	suite.NoError(err)
	resp, err = reportScoreWithMetadata(lbName, entryID, 15)
	suite.NoError(err)
	resp, err = reportScoreWithMetadata(lbName, entryID, 5)
	suite.NoError(err)
	suite.Equal(float64(5), resp.Score)
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
	r, err := listScores(lbName)
	fmt.Println(r)

	suite.NoError(err)
}

func (suite *E2ETestSuite) TestFullSizeLeaderboards() {
	country := "PT"
	league := "gold"
	lbName := defaultLbNameSumMultiple
	for i := 0; i < 100; i++ {
		entryID := testutil.NewID()
		if i%2 == 0 {
			country = "PT"
			league = "silver"
		} else {
			country = "UK"
			league = "gold"
		}
		score := float64((i * 10) + 10)
		_, err := reportScoreWithMetadataCountryLeage(lbName, entryID, score, country, league)
		suite.NoError(err)
	}
	list, err := listScoresWithMetadata(lbName, metadataDefault)
	suite.NoError(err)
	suite.Len(list.Scores, 3)
	suite.Len(list.Scores[0].Scores, 51)
	suite.Len(list.Scores[1].Scores, 50)
	suite.Len(list.Scores[2].Scores, 50)
	list, err = listScoresWithMetadata(lbName, map[string]string{
		"country": "UK",
		"league":  "silver",
	})
	suite.NoError(err)
	suite.Len(list.Scores, 3)
	suite.Len(list.Scores[0].Scores, 51)
	suite.Len(list.Scores[1].Scores, 50)
	suite.Len(list.Scores[2].Scores, 50)
	list, err = listScoresWithMetadata(lbName, map[string]string{
		"country": "PT",
		"league":  "silver",
	})
	suite.NoError(err)
	suite.Len(list.Scores, 3)
	suite.Len(list.Scores[0].Scores, 51)
	suite.Len(list.Scores[1].Scores, 50)
	suite.Len(list.Scores[2].Scores, 50)
	list, err = listScoresWithMetadata(lbName, map[string]string{
		"country": "PT",
		"league":  "gold",
	})
	suite.NoError(err)
	suite.Len(list.Scores, 3)
	suite.Len(list.Scores[0].Scores, 51)
	suite.Len(list.Scores[1].Scores, 50)
	suite.Len(list.Scores[2].Scores, 50)

	for _, v := range list.Scores[0].Scores {
		fmt.Println(v)
	}
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

func reportScoreWithMetadataCountryLeage(ldbName string, entry string, score float64, country string, league string) (response putScoreResponse, err error) {
	path := fmt.Sprintf("/api/v1/score/%s", ldbName)

	req := putScoreRequest{
		Entry: entry,
		Score: score,
		Metadata: domain.Metadata{
			"country": country,
			"league":  league,
		},
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
func reportScoreWithMetadata(ldbName string, entry string, score float64) (response putScoreResponse, err error) {
	path := fmt.Sprintf("/api/v1/score/%s", ldbName)

	req := putScoreRequest{
		Entry: entry,
		Score: score,
		Metadata: domain.Metadata{
			"country": "PT",
			"league":  "gold",
		},
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

// TODO: we need to pass metadata as query parameters
func listScoresWithMetadata(lbname string, meta domain.Metadata) (response listScoresResponse, err error) {
	var qry url.Values = make(map[string][]string)
	// add metadata to query parameters
	for k, v := range meta {
		qry.Add(fmt.Sprintf("meta_%s", k), v)
	}
	encode := qry.Encode()

	path := fmt.Sprintf("/api/v1/scores/%s?%s", lbname, encode)
	err = requests.
		URL(path).
		Host(baseURL).
		Scheme(defaultScheme).
		CheckStatus(http.StatusOK).
		ToJSON(&response).
		Fetch(context.Background())

	return response, err

}
