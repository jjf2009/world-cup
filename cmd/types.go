package cmd

// ========================================
// Match Results
// ========================================

type MatchScoresResponse struct {
	StatusCode int              `json:"status_code"`
	Matches    map[string]Match `json:"matches"`
	LastUpdated string          `json:"last_updated"`
}

type Match struct {
	ID string `json:"id"`

	MatchNumber string `json:"match_number"`

	Status string `json:"status"`

	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`

	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`

	Venue string `json:"venue"`

	Date string `json:"date"`

	Result string `json:"result,omitempty"`
}

// ========================================
// Live Matches
// ========================================

type LiveMatchResponse struct {
	StatusCode int `json:"status_code"`

	LiveCount int `json:"live_count"`

	Matches map[string]LiveMatch `json:"matches"`

	LastUpdated string `json:"last_updated"`
}

type LiveMatch struct {
	ID string `json:"id"`

	MatchNumber string `json:"match_number"`

	Status string `json:"status"`

	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`

	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`

	Minute string `json:"minute"`

	Venue string `json:"venue"`
}

// ========================================
// Schedule
// ========================================

type MatchScheduleResponse struct {
	StatusCode int `json:"status_code"`

	Schedule map[string]ScheduledMatch `json:"schedule"`
}

type ScheduledMatch struct {
	MatchNumber string `json:"match_number"`

	HomeTeam string `json:"home_team"`
	AwayTeam string `json:"away_team"`

	Venue string `json:"venue"`

	Date string `json:"date"`
	Time string `json:"time"`
}

// ========================================
// Points Table
// ========================================

type PointsTableResponse struct {
	StatusCode int `json:"status_code"`

	Standings map[string]TeamStanding `json:"standings"`
}

type TeamStanding struct {
	Position int `json:"position"`

	Team string `json:"team"`

	Played int `json:"played"`
	Won    int `json:"won"`
	Drawn  int `json:"drawn"`
	Lost   int `json:"lost"`

	GoalsFor        int `json:"goals_for"`
	GoalsAgainst    int `json:"goals_against"`
	GoalDifference  int `json:"goal_difference"`

	Points int `json:"points"`
}

// ========================================
// Historical Winners
// ========================================

type HistoricalWinnersResponse struct {
	StatusCode int `json:"status_code"`

	Winners map[string]YearWinner `json:"winners"`
}

type YearWinner struct {
	Winner string `json:"winner"`

	RunnerUp string `json:"runner_up"`

	Venue string `json:"venue"`
}

// ========================================
// Squads
// ========================================

type SquadResponse struct {
	StatusCode int `json:"status_code"`

	Team string `json:"team"`

	Squad map[string]SquadPlayer `json:"squad"`
}

type SquadPlayer struct {
	Name string `json:"name"`

	Nationality string `json:"nationality"`

	Position string `json:"position"`

	JerseyNumber int `json:"jersey_number"`
}