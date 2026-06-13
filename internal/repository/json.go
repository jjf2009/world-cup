package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/h0i5/ipl/internal/domain"
)

const kickoffLayout = "01/02/2006 15:04"

type JSONTeamRepository struct {
	teams []domain.Team
	byID  map[string]domain.Team
}

type JSONStadiumRepository struct {
	stadiums []domain.Stadium
	byID     map[string]domain.Stadium
}

type JSONMatchRepository struct {
	matches []domain.Match
}

type JSONStandingRepository struct {
	standings []domain.Standing
}

type JSONWinnerRepository struct {
	winners []domain.Winner
}

func NewJSONTeamRepository(path string) (*JSONTeamRepository, error) {
	var rows []teamRow
	if err := loadJSON(path, &rows); err != nil {
		return nil, err
	}

	repo := &JSONTeamRepository{byID: make(map[string]domain.Team, len(rows))}
	for _, row := range rows {
		if row.ID == "" {
			return nil, fmt.Errorf("team in %s has empty id", path)
		}
		team := domain.Team{
			ID:       row.ID,
			Name:     row.NameEN,
			FIFACode: row.FIFACode,
			ISO2:     row.ISO2,
			Group:    row.Group,
			FlagURL:  row.Flag,
		}
		if _, exists := repo.byID[team.ID]; exists {
			return nil, fmt.Errorf("duplicate team id %q in %s", team.ID, path)
		}
		repo.teams = append(repo.teams, team)
		repo.byID[team.ID] = team
	}

	sort.Slice(repo.teams, func(i, j int) bool {
		return naturalLess(repo.teams[i].ID, repo.teams[j].ID)
	})
	return repo, nil
}

func (r *JSONTeamRepository) All(context.Context) ([]domain.Team, error) {
	return append([]domain.Team(nil), r.teams...), nil
}

func (r *JSONTeamRepository) ByID(_ context.Context, id string) (domain.Team, error) {
	team, ok := r.byID[id]
	if !ok {
		return domain.Team{}, fmt.Errorf("team %q: %w", id, domain.ErrNotFound)
	}
	return team, nil
}

func NewJSONStadiumRepository(path string) (*JSONStadiumRepository, error) {
	var rows []stadiumRow
	if err := loadJSON(path, &rows); err != nil {
		return nil, err
	}

	repo := &JSONStadiumRepository{byID: make(map[string]domain.Stadium, len(rows))}
	for _, row := range rows {
		if row.ID == "" {
			return nil, fmt.Errorf("stadium in %s has empty id", path)
		}
		stadium := domain.Stadium{
			ID:       row.ID,
			Name:     row.NameEN,
			FIFAName: row.FIFAName,
			City:     row.CityEN,
			Country:  row.CountryEN,
			Capacity: row.Capacity,
			Region:   row.Region,
		}
		if _, exists := repo.byID[stadium.ID]; exists {
			return nil, fmt.Errorf("duplicate stadium id %q in %s", stadium.ID, path)
		}
		repo.stadiums = append(repo.stadiums, stadium)
		repo.byID[stadium.ID] = stadium
	}

	sort.Slice(repo.stadiums, func(i, j int) bool {
		return naturalLess(repo.stadiums[i].ID, repo.stadiums[j].ID)
	})
	return repo, nil
}

func (r *JSONStadiumRepository) All(context.Context) ([]domain.Stadium, error) {
	return append([]domain.Stadium(nil), r.stadiums...), nil
}

func (r *JSONStadiumRepository) ByID(_ context.Context, id string) (domain.Stadium, error) {
	stadium, ok := r.byID[id]
	if !ok {
		return domain.Stadium{}, fmt.Errorf("stadium %q: %w", id, domain.ErrNotFound)
	}
	return stadium, nil
}

func NewJSONMatchRepository(path string, loc *time.Location) (*JSONMatchRepository, error) {
	var rows []matchRow
	if err := loadJSON(path, &rows); err != nil {
		return nil, err
	}
	if loc == nil {
		loc = time.Local
	}

	seen := make(map[string]struct{}, len(rows))
	matches := make([]domain.Match, 0, len(rows))
	for _, row := range rows {
		if row.ID == "" {
			return nil, fmt.Errorf("match in %s has empty id", path)
		}
		if _, exists := seen[row.ID]; exists {
			return nil, fmt.Errorf("duplicate match id %q in %s", row.ID, path)
		}
		seen[row.ID] = struct{}{}

		kickoff, err := time.ParseInLocation(kickoffLayout, row.LocalDate, loc)
		if err != nil {
			return nil, fmt.Errorf("match %s kickoff %q: %w", row.ID, row.LocalDate, err)
		}

		matchday, err := atoiField("matchday", row.Matchday)
		if err != nil {
			return nil, fmt.Errorf("match %s: %w", row.ID, err)
		}
		homeScore, err := atoiField("home_score", row.HomeScore)
		if err != nil {
			return nil, fmt.Errorf("match %s: %w", row.ID, err)
		}
		awayScore, err := atoiField("away_score", row.AwayScore)
		if err != nil {
			return nil, fmt.Errorf("match %s: %w", row.ID, err)
		}

		matches = append(matches, domain.Match{
			ID:            row.ID,
			MatchNumber:   row.ID,
			HomeTeamID:    row.HomeTeamID,
			AwayTeamID:    row.AwayTeamID,
			HomeTeamLabel: row.HomeTeamLabel,
			AwayTeamLabel: row.AwayTeamLabel,
			HomeScore:     homeScore,
			AwayScore:     awayScore,
			HomeScorers:   splitScorers(row.HomeScorers),
			AwayScorers:   splitScorers(row.AwayScorers),
			Group:         row.Group,
			Matchday:      matchday,
			Kickoff:       kickoff,
			StadiumID:     row.StadiumID,
			Finished:      parseBool(row.Finished),
			TimeElapsed:   row.TimeElapsed,
			Type:          row.Type,
		})
	}

	sortMatches(matches)
	return &JSONMatchRepository{matches: matches}, nil
}

func (r *JSONMatchRepository) All(context.Context) ([]domain.Match, error) {
	return append([]domain.Match(nil), r.matches...), nil
}

func NewJSONStandingRepository(path string) (*JSONStandingRepository, error) {
	var groups []standingGroupRow
	if err := loadJSON(path, &groups); err != nil {
		return nil, err
	}

	var standings []domain.Standing
	for _, group := range groups {
		for _, row := range group.Teams {
			standing, err := standingFromRow(group.Group, row)
			if err != nil {
				return nil, fmt.Errorf("group %s team %s: %w", group.Group, row.TeamID, err)
			}
			standings = append(standings, standing)
		}
	}
	return &JSONStandingRepository{standings: standings}, nil
}

func (r *JSONStandingRepository) All(context.Context) ([]domain.Standing, error) {
	return append([]domain.Standing(nil), r.standings...), nil
}

func NewJSONWinnerRepository(path string) (*JSONWinnerRepository, error) {
	var rows []winnerRow
	if err := loadJSON(path, &rows); err != nil {
		return nil, err
	}

	winners := make([]domain.Winner, 0, len(rows))
	for _, row := range rows {
		winners = append(winners, domain.Winner{
			Year:     row.Year,
			Winner:   row.Winner,
			RunnerUp: row.RunnerUp,
			Venue:    row.Venue,
		})
	}
	sort.Slice(winners, func(i, j int) bool {
		return winners[i].Year > winners[j].Year
	})
	return &JSONWinnerRepository{winners: winners}, nil
}

func (r *JSONWinnerRepository) All(context.Context) ([]domain.Winner, error) {
	return append([]domain.Winner(nil), r.winners...), nil
}

func loadJSON(path string, dst any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(dst); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func standingFromRow(group string, row standingTeamRow) (domain.Standing, error) {
	mp, err := atoiField("mp", row.Played)
	if err != nil {
		return domain.Standing{}, err
	}
	won, err := atoiField("w", row.Won)
	if err != nil {
		return domain.Standing{}, err
	}
	lost, err := atoiField("l", row.Lost)
	if err != nil {
		return domain.Standing{}, err
	}
	drawn, err := atoiField("d", row.Drawn)
	if err != nil {
		return domain.Standing{}, err
	}
	points, err := atoiField("pts", row.Points)
	if err != nil {
		return domain.Standing{}, err
	}
	gf, err := atoiField("gf", row.GoalsFor)
	if err != nil {
		return domain.Standing{}, err
	}
	ga, err := atoiField("ga", row.GoalsAgainst)
	if err != nil {
		return domain.Standing{}, err
	}
	gd, err := atoiField("gd", row.GoalDifference)
	if err != nil {
		return domain.Standing{}, err
	}

	return domain.Standing{
		Group:          group,
		TeamID:         row.TeamID,
		Played:         mp,
		Won:            won,
		Lost:           lost,
		Drawn:          drawn,
		Points:         points,
		GoalsFor:       gf,
		GoalsAgainst:   ga,
		GoalDifference: gd,
	}, nil
}

func atoiField(name, value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("%s is empty", name)
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s %q: %w", name, value, err)
	}
	return n, nil
}

func parseBool(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "TRUE", "YES", "1":
		return true
	default:
		return false
	}
}

func splitScorers(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "null") {
		return nil
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';'
	})
	scorers := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			scorers = append(scorers, part)
		}
	}
	return scorers
}

func sortMatches(matches []domain.Match) {
	sort.Slice(matches, func(i, j int) bool {
		if !matches[i].Kickoff.Equal(matches[j].Kickoff) {
			return matches[i].Kickoff.Before(matches[j].Kickoff)
		}
		return naturalLess(matches[i].ID, matches[j].ID)
	})
}

func naturalLess(a, b string) bool {
	ai, aerr := strconv.Atoi(a)
	bi, berr := strconv.Atoi(b)
	if aerr == nil && berr == nil {
		return ai < bi
	}
	return a < b
}

func IsNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

type teamRow struct {
	ID       string `json:"id"`
	NameEN   string `json:"name_en"`
	Flag     string `json:"flag"`
	FIFACode string `json:"fifa_code"`
	ISO2     string `json:"iso2"`
	Group    string `json:"groups"`
}

type stadiumRow struct {
	ID        string `json:"id"`
	NameEN    string `json:"name_en"`
	FIFAName  string `json:"fifa_name"`
	CityEN    string `json:"city_en"`
	CountryEN string `json:"country_en"`
	Capacity  int    `json:"capacity"`
	Region    string `json:"region"`
}

type matchRow struct {
	ID            string `json:"id"`
	HomeTeamID    string `json:"home_team_id"`
	AwayTeamID    string `json:"away_team_id"`
	HomeScore     string `json:"home_score"`
	AwayScore     string `json:"away_score"`
	HomeScorers   string `json:"home_scorers"`
	AwayScorers   string `json:"away_scorers"`
	Group         string `json:"group"`
	Matchday      string `json:"matchday"`
	LocalDate     string `json:"local_date"`
	StadiumID     string `json:"stadium_id"`
	Finished      string `json:"finished"`
	TimeElapsed   string `json:"time_elapsed"`
	Type          string `json:"type"`
	HomeTeamLabel string `json:"home_team_label"`
	AwayTeamLabel string `json:"away_team_label"`
}

type standingGroupRow struct {
	Group string            `json:"group"`
	Teams []standingTeamRow `json:"teams"`
}

type standingTeamRow struct {
	TeamID         string `json:"team_id"`
	Played         string `json:"mp"`
	Won            string `json:"w"`
	Lost           string `json:"l"`
	Drawn          string `json:"d"`
	Points         string `json:"pts"`
	GoalsFor       string `json:"gf"`
	GoalsAgainst   string `json:"ga"`
	GoalDifference string `json:"gd"`
}

type winnerRow struct {
	Year     int    `json:"year"`
	Winner   string `json:"winner"`
	RunnerUp string `json:"runner_up"`
	Venue    string `json:"venue"`
}
