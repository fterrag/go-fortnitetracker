package fortnitetracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseURL = "https://api.fortnitetracker.com/v1"

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type StatsResponse struct {
	AccountID        string         `json:"accountId"`
	PlatformID       int            `json:"platformId"`
	PlatformName     string         `json:"platformName"`
	PlatformNameLong string         `json:"platformNameLong"`
	EpicUserHandle   string         `json:"epicUserHandle"`
	Stats            Stats          `json:"stats"`
	RecentMatches    []Match        `json:"recentMatches"`
	LifetimeStats    []LifetimeStat `json:"lifetimeStats"`
}

type Stats struct {
	LifetimeSolo  StatsGroup `json:"p2"`
	LifetimeDuo   StatsGroup `json:"p10"`
	LifetimeSquad StatsGroup `json:"p9"`

	CurrentSeasonSolo  StatsGroup `json:"curr_p2"`
	CurrentSeasonDuo   StatsGroup `json:"curr_p10"`
	CurrentSeasonSquad StatsGroup `json:"curr_p9"`
}

type StatsGroup struct {
	TRNRating     Metric `json:"trnRating"`
	Score         Metric `json:"score"`
	Top1          Metric `json:"top1"`
	Top3          Metric `json:"top3"`
	Top5          Metric `json:"top5"`
	Top6          Metric `json:"top6"`
	Top10         Metric `json:"top10"`
	Top12         Metric `json:"top12"`
	Top25         Metric `json:"top25"`
	KD            Metric `json:"kd"`
	WinRatio      Metric `json:"wins"`
	Matches       Metric `json:"matches"`
	Kills         Metric `json:"kills"`
	KPG           Metric `json:kpg`
	ScorePerMatch Metric `json:"scorePerMatch"`
}

type Metric struct {
	Label        string  `json:"label"`
	Field        string  `json:"field"`
	Category     string  `json:"category"`
	ValueInt     int     `json:"valueInt"`
	Value        string  `json:"value"`
	Rank         int     `json:"rank"`
	Percentile   float64 `json:"percentile"`
	DisplayValue string  `json:"displayValue"`
}

type Match struct {
	ID            int    `json:"id"`
	AccountID     string `json:"accountId"`
	Playlist      string `json:"playlist"`
	Kills         int    `json:"kills"`
	MinutesPlayed int    `json:"minutesPlayed"`
	Top1          int    `json:"top1"`
	Matches       int    `json:"matches"`
	DateCollected string `json:"dateCollected"`
	Score         int    `json:"score"`
	Platform      int    `json:"platform"`
}

func (m Match) GetPlaylistName() string {
	switch m.Playlist {
	case "p2":
		return "Solo"
	case "p10":
		return "Duo"
	case "p9":
		return "Squad"
	default:
		return "Unknown"
	}
}

type LifetimeStat struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type response struct {
	httpResponse *http.Response
	body         []byte
}

type FortniteTracker struct {
	httpClient HTTPClient
	key        string
}

func NewFortniteTracker(httpClient HTTPClient, key string) *FortniteTracker {
	return &FortniteTracker{
		httpClient: httpClient,
		key:        key,
	}
}

func (f *FortniteTracker) request(url string) (*response, error) {
	url = strings.Join([]string{baseURL, url}, "")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while creting HTTP request: %s", err)
	}

	req.Header.Add("TRN-Api-Key", f.key)

	res, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while making HTTP request: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading HTTP response body: %s", err)
	}

	return &response{
		httpResponse: res,
		body:         body,
	}, nil
}

func (f *FortniteTracker) GetStats(platform string, displayName string) (*StatsResponse, error) {
	if platform != "pc" && platform != "xbl" && platform != "psn" {
		return nil, errors.New("platform must contain a value of pc, xbl, or psn")
	}

	url := "/profile/{platform}/{displayName}"
	url = strings.Replace(url, "{platform}", platform, 1)
	url = strings.Replace(url, "{displayName}", displayName, 1)

	res, err := f.request(url)
	if err != nil {
		return nil, fmt.Errorf("Error while making HTTP request: %s", err)
	}

	if res.httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP response returned non-200 status code: %s", res.httpResponse.Status)
	}

	stats := &StatsResponse{}

	err = json.Unmarshal(res.body, stats)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshaling HTTP body: %s", err)
	}

	return stats, nil
}
