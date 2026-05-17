package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sports-stats-pipeline/internal/model"
)

const openLigaBase = "https://api.openligadb.com"

// LeagueConfig — лига OpenLigaDB (футбол).
type LeagueConfig struct {
	ShortName string // bl1, bl2, ...
	Label     string
	Season    string // 2024, 2023/2024
}

var DefaultLeagues = []LeagueConfig{
	{ShortName: "bl1", Label: "Bundesliga", Season: "2024"},
	{ShortName: "bl2", Label: "2. Bundesliga", Season: "2024"},
}

type openLigaMatch struct {
	MatchID       int    `json:"matchID"`
	MatchDateTime string `json:"matchDateTime"`
	MatchIsFinished bool `json:"matchIsFinished"`
	Group         *struct {
		GroupName string `json:"groupName"`
	} `json:"group"`
	Matchday      int `json:"matchday"`
	Location      *struct {
		LocationStadium string `json:"locationStadium"`
	} `json:"location"`
	Team1         *teamScore `json:"team1"`
	Team2         *teamScore `json:"team2"`
}

type teamScore struct {
	TeamName string `json:"teamName"`
	Goals    *int   `json:"goals"`
}

// LeagueMatchesURL формирует URL эндпоинта (base — корень API, например httptest.Server.URL).
func LeagueMatchesURL(base string, league LeagueConfig) string {
	if base == "" {
		base = openLigaBase
	}
	return fmt.Sprintf("%s/getmatchdata/%s/%s", base, league.ShortName, league.Season)
}

// FetchLeagueMatches загружает матчи лиги через HTTP API.
func FetchLeagueMatches(ctx context.Context, client *http.Client, league LeagueConfig) ([]model.Match, error) {
	return FetchLeagueMatchesAt(ctx, client, league, openLigaBase)
}

// FetchLeagueMatchesAt то же, с настраиваемым base URL (для тестов).
func FetchLeagueMatchesAt(
	ctx context.Context,
	client *http.Client,
	league LeagueConfig,
	baseURL string,
) ([]model.Match, error) {
	url := LeagueMatchesURL(baseURL, league)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sports-stats-collector/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openligadb %s: HTTP %d", league.ShortName, resp.StatusCode)
	}

	var raw []openLigaMatch
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return MapOpenLigaMatches(raw, league, time.Now().UTC().Format(time.RFC3339)), nil
}

// MapOpenLigaMatches преобразует ответ API в модель (тестируемая логика).
func MapOpenLigaMatches(raw []openLigaMatch, league LeagueConfig, collectedAt string) []model.Match {
	out := make([]model.Match, 0, len(raw))
	for _, m := range raw {
		rec := model.Match{
			MatchID:     m.MatchID,
			League:      league.Label,
			Season:      league.Season,
			MatchDate:   m.MatchDateTime,
			Matchday:    m.Matchday,
			Status:      "scheduled",
			SourceAPI:   "openligadb",
			CollectedAt: collectedAt,
		}
		if m.MatchIsFinished {
			rec.Status = "finished"
		}
		if m.Group != nil {
			rec.GroupName = m.Group.GroupName
		}
		if m.Location != nil {
			rec.Location = m.Location.LocationStadium
		}
		if m.Team1 != nil {
			rec.HomeTeam = m.Team1.TeamName
			rec.HomeGoals = m.Team1.Goals
		}
		if m.Team2 != nil {
			rec.AwayTeam = m.Team2.TeamName
			rec.AwayGoals = m.Team2.Goals
		}
		out = append(out, rec)
	}
	return out
}
