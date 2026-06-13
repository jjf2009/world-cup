package repository

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/h0i5/ipl/internal/domain"
)

func TestJSONRepositoriesNormalizeData(t *testing.T) {
	dir := t.TempDir()

	teamsPath := writeFixture(t, dir, "teams.json", `[
		{"id":"1","name_en":"Mexico","flag":"flag","fifa_code":"MEX","iso2":"MX","groups":"A"}
	]`)
	stadiumsPath := writeFixture(t, dir, "stadiums.json", `[
		{"id":"1","name_en":"Estadio Azteca","fifa_name":"Mexico City Stadium","city_en":"Mexico City","country_en":"Mexico","capacity":83000,"region":"Central"}
	]`)
	matchesPath := writeFixture(t, dir, "matches.json", `[
		{"id":"1","home_team_id":"1","away_team_id":"0","home_score":"2","away_score":"1","home_scorers":"A, B","away_scorers":"null","group":"A","matchday":"1","local_date":"06/11/2026 13:00","stadium_id":"1","finished":"TRUE","time_elapsed":"90","type":"group","away_team_label":"Winner Match 2"}
	]`)
	standingsPath := writeFixture(t, dir, "standings.json", `[
		{"group":"A","teams":[{"team_id":"1","mp":"1","w":"1","l":"0","d":"0","pts":"3","gf":"2","ga":"1","gd":"1"}]}
	]`)
	winnersPath := writeFixture(t, dir, "winners.json", `[
		{"year":2022,"winner":"Argentina","runner_up":"France"}
	]`)

	teams, err := NewJSONTeamRepository(teamsPath)
	if err != nil {
		t.Fatal(err)
	}
	stadiums, err := NewJSONStadiumRepository(stadiumsPath)
	if err != nil {
		t.Fatal(err)
	}
	matches, err := NewJSONMatchRepository(matchesPath, time.UTC)
	if err != nil {
		t.Fatal(err)
	}
	standings, err := NewJSONStandingRepository(standingsPath)
	if err != nil {
		t.Fatal(err)
	}
	winners, err := NewJSONWinnerRepository(winnersPath)
	if err != nil {
		t.Fatal(err)
	}

	team, err := teams.ByID(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "Mexico" || team.Group != "A" {
		t.Fatalf("unexpected team: %+v", team)
	}

	stadium, err := stadiums.ByID(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if stadium.DisplayName() != "Mexico City Stadium" {
		t.Fatalf("unexpected stadium: %+v", stadium)
	}

	allMatches, err := matches.All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := allMatches[0]; got.HomeScore != 2 || got.AwayScore != 1 || !got.Finished || len(got.HomeScorers) != 2 || len(got.AwayScorers) != 0 {
		t.Fatalf("unexpected match normalization: %+v", got)
	}

	allStandings, err := standings.All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if allStandings[0].Points != 3 || allStandings[0].GoalDifference != 1 {
		t.Fatalf("unexpected standing normalization: %+v", allStandings[0])
	}

	allWinners, err := winners.All(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if allWinners[0].Year != 2022 || allWinners[0].Venue != "" {
		t.Fatalf("unexpected winner normalization: %+v", allWinners[0])
	}
}

func TestJSONTeamRepositoryDuplicateAndNotFound(t *testing.T) {
	dir := t.TempDir()
	path := writeFixture(t, dir, "teams.json", `[
		{"id":"1","name_en":"Mexico"},
		{"id":"1","name_en":"Canada"}
	]`)

	if _, err := NewJSONTeamRepository(path); err == nil {
		t.Fatal("expected duplicate id error")
	}

	path = writeFixture(t, dir, "teams_ok.json", `[{"id":"1","name_en":"Mexico"}]`)
	repo, err := NewJSONTeamRepository(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.ByID(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func writeFixture(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}
