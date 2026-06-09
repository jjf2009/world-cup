package domain

import (
	"errors"
	"time"
)

const (
	StatusUpcoming = "UPCOMING"
	StatusLive     = "LIVE"
	StatusFinished = "FINISHED"
)

var ErrNotFound = errors.New("not found")

type Team struct {
	ID       string
	Name     string
	FIFACode string
	ISO2     string
	Group    string
	FlagURL  string
}

type Stadium struct {
	ID       string
	Name     string
	FIFAName string
	City     string
	Country  string
	Capacity int
	Region   string
}

func (s Stadium) DisplayName() string {
	if s.FIFAName != "" {
		return s.FIFAName
	}
	return s.Name
}

type Match struct {
	ID            string
	MatchNumber   string
	HomeTeamID    string
	AwayTeamID    string
	HomeTeamLabel string
	AwayTeamLabel string
	HomeScore     int
	AwayScore     int
	HomeScorers   []string
	AwayScorers   []string
	Group         string
	Matchday      int
	Kickoff       time.Time
	StadiumID     string
	Finished      bool
	TimeElapsed   string
	Type          string
}

type Standing struct {
	Group          string
	TeamID         string
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	Points         int
}

type Winner struct {
	Year     int
	Winner   string
	RunnerUp string
	Venue    string
}

type Scorer struct {
	Team   string
	Name   string
	Minute int
}

type LiveMatch struct {
	ID          string
	MatchNumber string
	HomeTeam    string
	AwayTeam    string
	HomeScore   int
	AwayScore   int
	Minute      int
	Status      string
	Stadium     string
	Group       string
	Scorers     []Scorer
	LastUpdated time.Time
}

type MatchView struct {
	ID          string
	MatchNumber string
	Status      string
	HomeTeam    string
	AwayTeam    string
	HomeScore   int
	AwayScore   int
	Venue       string
	Group       string
	Date        string
	Time        string
	Kickoff     time.Time
	Result      string
	Scorers     []Scorer
}

type StandingView struct {
	Group          string
	Position       int
	Team           string
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	Points         int
}

type WinnerView struct {
	Year     int
	Winner   string
	RunnerUp string
	Venue    string
}

type LiveMatchView struct {
	ID          string
	MatchNumber string
	Status      string
	HomeTeam    string
	AwayTeam    string
	HomeScore   int
	AwayScore   int
	Minute      int
	Venue       string
	Group       string
	Scorers     []Scorer
	LastUpdated time.Time
}

type ScheduleView struct {
	Today    []MatchView
	Tomorrow []MatchView
}
