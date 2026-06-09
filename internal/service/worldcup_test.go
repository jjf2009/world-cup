package service

import (
	"context"
	"testing"
	"time"

	"github.com/h0i5/ipl/internal/domain"
)

func TestWorldCupServiceResolvesMatchesAndFallbackSchedule(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 6, 9, 9, 0, 0, 0, time.UTC)

	svc := newTestService()
	svc.now = func() time.Time { return now }

	matches, err := svc.GetAllMatches(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got := matches[0]; got.HomeTeam != "Mexico" || got.AwayTeam != "South Africa" || got.Venue != "Mexico City Stadium" {
		t.Fatalf("match was not resolved: %+v", got)
	}
	if got := matches[1]; got.HomeTeam != "Winner Match 1" || got.AwayTeam != "Canada" {
		t.Fatalf("knockout fallback labels not resolved: %+v", got)
	}

	today, err := svc.GetTodayMatches(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(today) != 3 {
		t.Fatalf("expected fallback to 3 upcoming matches, got %d", len(today))
	}
	if today[0].MatchNumber != "1" || today[2].MatchNumber != "3" {
		t.Fatalf("unexpected fallback order: %+v", today)
	}
}

func TestWorldCupServiceSortsStandingsAndWinners(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()

	standings, err := svc.GetStandings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if standings[0].Team != "Canada" || standings[0].Position != 1 {
		t.Fatalf("expected Canada first on points, got %+v", standings[0])
	}
	if standings[1].Team != "Mexico" || standings[1].Position != 2 {
		t.Fatalf("expected Mexico second, got %+v", standings[1])
	}

	winners, err := svc.GetHistoricalWinners(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if winners[0].Year != 2022 || winners[1].Year != 2018 {
		t.Fatalf("winners not sorted newest first: %+v", winners)
	}
}

func newTestService() *WorldCupService {
	kickoff := func(day int) time.Time {
		return time.Date(2026, 6, day, 13, 0, 0, 0, time.UTC)
	}
	teams := fakeTeams{byID: map[string]domain.Team{
		"1": {ID: "1", Name: "Mexico"},
		"2": {ID: "2", Name: "South Africa"},
		"3": {ID: "3", Name: "Canada"},
	}}
	stadiums := fakeStadiums{byID: map[string]domain.Stadium{
		"1": {ID: "1", Name: "Estadio Azteca", FIFAName: "Mexico City Stadium"},
	}}
	matches := fakeMatches{matches: []domain.Match{
		{ID: "1", MatchNumber: "1", HomeTeamID: "1", AwayTeamID: "2", StadiumID: "1", Group: "A", Kickoff: kickoff(11)},
		{ID: "2", MatchNumber: "2", HomeTeamID: "0", HomeTeamLabel: "Winner Match 1", AwayTeamID: "3", StadiumID: "1", Group: "R32", Kickoff: kickoff(12)},
		{ID: "3", MatchNumber: "3", HomeTeamID: "1", AwayTeamID: "3", StadiumID: "1", Group: "A", Kickoff: kickoff(13)},
	}}
	standings := fakeStandings{standings: []domain.Standing{
		{Group: "A", TeamID: "1", Points: 3, GoalDifference: 1, GoalsFor: 2},
		{Group: "A", TeamID: "3", Points: 6, GoalDifference: 2, GoalsFor: 4},
	}}
	winners := fakeWinners{winners: []domain.Winner{
		{Year: 2018, Winner: "France", RunnerUp: "Croatia"},
		{Year: 2022, Winner: "Argentina", RunnerUp: "France"},
	}}
	live := fakeLive{match: domain.LiveMatch{
		ID:          "1",
		MatchNumber: "1",
		HomeTeam:    "Mexico",
		AwayTeam:    "South Africa",
		Status:      domain.StatusLive,
		Stadium:     "Mexico City Stadium",
		Group:       "A",
		Minute:      12,
	}}

	return NewWorldCupService(teams, stadiums, matches, standings, winners, live)
}

type fakeTeams struct{ byID map[string]domain.Team }

func (f fakeTeams) All(context.Context) ([]domain.Team, error) {
	teams := make([]domain.Team, 0, len(f.byID))
	for _, team := range f.byID {
		teams = append(teams, team)
	}
	return teams, nil
}

func (f fakeTeams) ByID(_ context.Context, id string) (domain.Team, error) {
	return f.byID[id], nil
}

type fakeStadiums struct{ byID map[string]domain.Stadium }

func (f fakeStadiums) All(context.Context) ([]domain.Stadium, error) {
	stadiums := make([]domain.Stadium, 0, len(f.byID))
	for _, stadium := range f.byID {
		stadiums = append(stadiums, stadium)
	}
	return stadiums, nil
}

func (f fakeStadiums) ByID(_ context.Context, id string) (domain.Stadium, error) {
	return f.byID[id], nil
}

type fakeMatches struct{ matches []domain.Match }

func (f fakeMatches) All(context.Context) ([]domain.Match, error) {
	return append([]domain.Match(nil), f.matches...), nil
}

type fakeStandings struct{ standings []domain.Standing }

func (f fakeStandings) All(context.Context) ([]domain.Standing, error) {
	return append([]domain.Standing(nil), f.standings...), nil
}

type fakeWinners struct{ winners []domain.Winner }

func (f fakeWinners) All(context.Context) ([]domain.Winner, error) {
	return append([]domain.Winner(nil), f.winners...), nil
}

type fakeLive struct{ match domain.LiveMatch }

func (f fakeLive) Current(context.Context) (domain.LiveMatch, error) {
	return f.match, nil
}
