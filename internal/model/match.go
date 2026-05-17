package model

// Match — унифицированная запись матча для JSONL и последующего анализа.
type Match struct {
	MatchID      int     `json:"match_id"`
	League       string  `json:"league"`
	Season       string  `json:"season"`
	MatchDate    string  `json:"match_date"`
	HomeTeam     string  `json:"home_team"`
	AwayTeam     string  `json:"away_team"`
	HomeGoals    *int    `json:"home_goals"`
	AwayGoals    *int    `json:"away_goals"`
	Status       string  `json:"status"`
	Matchday     int     `json:"matchday"`
	Location     string  `json:"location,omitempty"`
	GroupName    string  `json:"group_name,omitempty"`
	SourceAPI    string  `json:"source_api"`
	CollectedAt  string  `json:"collected_at"`
}
