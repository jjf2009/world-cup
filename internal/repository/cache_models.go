package repository

import "time"

// CacheWrapper adds update metadata to the cached data arrays.
type CacheWrapper[T any] struct {
	UpdatedAt time.Time `json:"updated_at"`
	Data      []T       `json:"data"`
}

type LiveMatch struct {
	ID        string `json:"id"`
	HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
	Minute    string `json:"minute"`
	Status    string `json:"status"`
	Venue     string `json:"venue"`
	Group     string `json:"group"`
}

type Fixture struct {
	ID          string `json:"id"`
	HomeTeam    string `json:"home_team"`
	AwayTeam    string `json:"away_team"`
	KickoffTime string `json:"kickoff_time"`
	Status      string `json:"status"`
	Venue       string `json:"venue"`
	HomeScore   int    `json:"home_score"`
	AwayScore   int    `json:"away_score"`
}

type Standing struct {
	Team           string `json:"team"`
	Played         int    `json:"played"`
	Won            int    `json:"won"`
	Drawn          int    `json:"drawn"`
	Lost           int    `json:"lost"`
	GoalsFor       int    `json:"goals_for"`
	GoalsAgainst   int    `json:"goals_against"`
	GoalDifference int    `json:"goal_difference"`
	Points         int    `json:"points"`
	Group          string `json:"group"`
}
