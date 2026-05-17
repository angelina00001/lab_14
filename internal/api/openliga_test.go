package api

import (
	"testing"
)

func TestMapOpenLigaMatches(t *testing.T) {
	league := LeagueConfig{ShortName: "bl1", Label: "Bundesliga", Season: "2024"}
	g4, g0 := 4, 0

	tests := []struct {
		name     string
		raw      []openLigaMatch
		wantLen  int
		wantStat string
		wantHome string
		wantHG   *int
	}{
		{
			name: "finished match with score",
			raw: []openLigaMatch{{
				MatchID:         1,
				MatchDateTime:   "2024-08-23T18:30:00",
				MatchIsFinished: true,
				Matchday:        1,
				Team1:           &teamScore{TeamName: "FC Bayern", Goals: &g4},
				Team2:           &teamScore{TeamName: "Werder", Goals: &g0},
			}},
			wantLen:  1,
			wantStat: "finished",
			wantHome: "FC Bayern",
			wantHG:   &g4,
		},
		{
			name: "scheduled without goals",
			raw: []openLigaMatch{{
				MatchID:         2,
				MatchDateTime:   "2025-05-17T15:30:00",
				MatchIsFinished: false,
				Matchday:        34,
				Team1:           &teamScore{TeamName: "Dortmund", Goals: nil},
				Team2:           &teamScore{TeamName: "Frankfurt", Goals: nil},
			}},
			wantLen:  1,
			wantStat: "scheduled",
			wantHome: "Dortmund",
			wantHG:   nil,
		},
		{
			name:    "empty payload",
			raw:     nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapOpenLigaMatches(tt.raw, league, "2026-01-01T00:00:00Z")
			if len(got) != tt.wantLen {
				t.Fatalf("len=%d want %d", len(got), tt.wantLen)
			}
			if tt.wantLen == 0 {
				return
			}
			m := got[0]
			if m.Status != tt.wantStat {
				t.Errorf("status=%q want %q", m.Status, tt.wantStat)
			}
			if m.HomeTeam != tt.wantHome {
				t.Errorf("home=%q want %q", m.HomeTeam, tt.wantHome)
			}
			if !ptrEqual(m.HomeGoals, tt.wantHG) {
				t.Errorf("home_goals=%v want %v", m.HomeGoals, tt.wantHG)
			}
			if m.League != league.Label {
				t.Errorf("league=%q", m.League)
			}
		})
	}
}

func ptrEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func TestMapOpenLigaMatches_groupAndLocation(t *testing.T) {
	league := LeagueConfig{Label: "BL", Season: "2024"}
	raw := []openLigaMatch{{
		MatchID: 10,
		Group: &struct {
			GroupName string `json:"groupName"`
		}{GroupName: "Round 1"},
		Location: &struct {
			LocationStadium string `json:"locationStadium"`
		}{LocationStadium: "Arena"},
		Team1: &teamScore{TeamName: "A"},
		Team2: &teamScore{TeamName: "B"},
	}}
	got := MapOpenLigaMatches(raw, league, "ts")
	if got[0].GroupName != "Round 1" || got[0].Location != "Arena" {
		t.Fatalf("got %+v", got[0])
	}
	if got[0].MatchID != 10 {
		t.Fatalf("match id")
	}
}
