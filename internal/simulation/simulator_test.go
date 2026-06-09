package simulation

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/h0i5/ipl/internal/domain"
)

func TestSimulatorProducesLiveSnapshotsAndKeepsScoresRealistic(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 6, 11, 9, 0, 0, 0, time.UTC)

	sim := NewSimulator(
		simMatches{matches: []domain.Match{
			{ID: "1", MatchNumber: "1", HomeTeamID: "1", AwayTeamID: "2", StadiumID: "1", Group: "A", Kickoff: now.Add(2 * time.Hour)},
			{ID: "2", MatchNumber: "2", HomeTeamID: "3", AwayTeamID: "4", StadiumID: "1", Group: "A", Kickoff: now.Add(4 * time.Hour)},
		}},
		simTeams{byID: map[string]domain.Team{
			"1": {ID: "1", Name: "Mexico"},
			"2": {ID: "2", Name: "South Africa"},
			"3": {ID: "3", Name: "Canada"},
			"4": {ID: "4", Name: "Brazil"},
		}},
		simStadiums{byID: map[string]domain.Stadium{
			"1": {ID: "1", FIFAName: "Mexico City Stadium"},
		}},
	)
	sim.now = func() time.Time { return now }
	sim.rand = rand.New(rand.NewSource(1))

	if err := sim.selectNext(ctx); err != nil {
		t.Fatal(err)
	}

	live, err := sim.Current(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if live.Status != domain.StatusLive || live.Minute != 1 || live.HomeTeam != "Mexico" || live.Stadium != "Mexico City Stadium" {
		t.Fatalf("unexpected initial live snapshot: %+v", live)
	}

	for i := 0; i < 40; i++ {
		sim.tick(ctx)
	}
	live, err = sim.Current(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if live.HomeScore+live.AwayScore > 6 {
		t.Fatalf("score is not realistic: %+v", live)
	}
	if live.Minute <= 1 || live.Minute > 90 {
		t.Fatalf("minute did not advance realistically: %+v", live)
	}
}

func TestSimulatorMovesToAnotherLiveMatchAfterFinish(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 6, 11, 9, 0, 0, 0, time.UTC)

	sim := NewSimulator(
		simMatches{matches: []domain.Match{
			{ID: "1", MatchNumber: "1", HomeTeamID: "1", AwayTeamID: "2", StadiumID: "1", Group: "A", Kickoff: now.Add(2 * time.Hour)},
			{ID: "2", MatchNumber: "2", HomeTeamID: "3", AwayTeamID: "4", StadiumID: "1", Group: "A", Kickoff: now.Add(4 * time.Hour)},
		}},
		simTeams{byID: map[string]domain.Team{
			"1": {ID: "1", Name: "Mexico"},
			"2": {ID: "2", Name: "South Africa"},
			"3": {ID: "3", Name: "Canada"},
			"4": {ID: "4", Name: "Brazil"},
		}},
		simStadiums{byID: map[string]domain.Stadium{
			"1": {ID: "1", FIFAName: "Mexico City Stadium"},
		}},
	)
	sim.now = func() time.Time { return now }
	sim.rand = rand.New(rand.NewSource(2))

	if err := sim.selectNext(ctx); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 89; i++ {
		sim.tick(ctx)
	}

	live, err := sim.Current(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if live.ID != "2" || live.Status != domain.StatusLive {
		t.Fatalf("expected second match to become live, got %+v", live)
	}
}

type simMatches struct{ matches []domain.Match }

func (s simMatches) All(context.Context) ([]domain.Match, error) {
	return append([]domain.Match(nil), s.matches...), nil
}

type simTeams struct{ byID map[string]domain.Team }

func (s simTeams) All(context.Context) ([]domain.Team, error) {
	teams := make([]domain.Team, 0, len(s.byID))
	for _, team := range s.byID {
		teams = append(teams, team)
	}
	return teams, nil
}

func (s simTeams) ByID(_ context.Context, id string) (domain.Team, error) {
	return s.byID[id], nil
}

type simStadiums struct{ byID map[string]domain.Stadium }

func (s simStadiums) All(context.Context) ([]domain.Stadium, error) {
	stadiums := make([]domain.Stadium, 0, len(s.byID))
	for _, stadium := range s.byID {
		stadiums = append(stadiums, stadium)
	}
	return stadiums, nil
}

func (s simStadiums) ByID(_ context.Context, id string) (domain.Stadium, error) {
	return s.byID[id], nil
}
