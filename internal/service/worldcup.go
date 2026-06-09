package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/h0i5/ipl/internal/domain"
	"github.com/h0i5/ipl/internal/repository"
)

type WorldCupService struct {
	teams     repository.TeamRepository
	stadiums  repository.StadiumRepository
	matches   repository.MatchRepository
	standings repository.StandingRepository
	winners   repository.WinnerRepository
	live      repository.LiveMatchProvider
	now       func() time.Time
}

func NewWorldCupService(
	teams repository.TeamRepository,
	stadiums repository.StadiumRepository,
	matches repository.MatchRepository,
	standings repository.StandingRepository,
	winners repository.WinnerRepository,
	live repository.LiveMatchProvider,
) *WorldCupService {
	return &WorldCupService{
		teams:     teams,
		stadiums:  stadiums,
		matches:   matches,
		standings: standings,
		winners:   winners,
		live:      live,
		now:       time.Now,
	}
}

func (s *WorldCupService) GetLiveMatch(ctx context.Context) (domain.LiveMatchView, error) {
	match, err := s.live.Current(ctx)
	if err != nil {
		return domain.LiveMatchView{}, err
	}

	return domain.LiveMatchView{
		ID:          match.ID,
		MatchNumber: match.MatchNumber,
		Status:      match.Status,
		HomeTeam:    match.HomeTeam,
		AwayTeam:    match.AwayTeam,
		HomeScore:   match.HomeScore,
		AwayScore:   match.AwayScore,
		Minute:      match.Minute,
		Venue:       match.Stadium,
		Group:       match.Group,
		Scorers:     append([]domain.Scorer(nil), match.Scorers...),
		LastUpdated: match.LastUpdated,
	}, nil
}

func (s *WorldCupService) GetTodayMatches(ctx context.Context) ([]domain.MatchView, error) {
	return s.matchesForDate(ctx, dateOnly(s.now()), 3)
}

func (s *WorldCupService) GetTomorrowMatches(ctx context.Context) ([]domain.MatchView, error) {
	return s.matchesForDate(ctx, dateOnly(s.now().AddDate(0, 0, 1)), 3)
}

func (s *WorldCupService) GetStandings(ctx context.Context) ([]domain.StandingView, error) {
	standings, err := s.standings.All(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]domain.StandingView, 0, len(standings))
	for _, standing := range standings {
		team, err := s.teams.ByID(ctx, standing.TeamID)
		if err != nil {
			return nil, err
		}
		views = append(views, domain.StandingView{
			Group:          standing.Group,
			Team:           team.Name,
			Played:         standing.Played,
			Won:            standing.Won,
			Drawn:          standing.Drawn,
			Lost:           standing.Lost,
			GoalsFor:       standing.GoalsFor,
			GoalsAgainst:   standing.GoalsAgainst,
			GoalDifference: standing.GoalDifference,
			Points:         standing.Points,
		})
	}

	sort.SliceStable(views, func(i, j int) bool {
		a, b := views[i], views[j]
		if a.Group != b.Group {
			return a.Group < b.Group
		}
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GoalDifference != b.GoalDifference {
			return a.GoalDifference > b.GoalDifference
		}
		if a.GoalsFor != b.GoalsFor {
			return a.GoalsFor > b.GoalsFor
		}
		return a.Team < b.Team
	})

	groupPosition := map[string]int{}
	for i := range views {
		groupPosition[views[i].Group]++
		views[i].Position = groupPosition[views[i].Group]
	}

	return views, nil
}

func (s *WorldCupService) GetHistoricalWinners(ctx context.Context) ([]domain.WinnerView, error) {
	winners, err := s.winners.All(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]domain.WinnerView, 0, len(winners))
	for _, winner := range winners {
		views = append(views, domain.WinnerView{
			Year:     winner.Year,
			Winner:   winner.Winner,
			RunnerUp: winner.RunnerUp,
			Venue:    winner.Venue,
		})
	}
	sort.Slice(views, func(i, j int) bool {
		return views[i].Year > views[j].Year
	})
	return views, nil
}

func (s *WorldCupService) GetAllMatches(ctx context.Context) ([]domain.MatchView, error) {
	matches, err := s.matches.All(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]domain.MatchView, 0, len(matches))
	for _, match := range matches {
		view, err := s.buildMatchView(ctx, match)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}

	sortMatchViews(views)
	return views, nil
}

func (s *WorldCupService) GetTeamByID(ctx context.Context, id string) (domain.Team, error) {
	return s.teams.ByID(ctx, id)
}

func (s *WorldCupService) GetStadiumByID(ctx context.Context, id string) (domain.Stadium, error) {
	return s.stadiums.ByID(ctx, id)
}

func (s *WorldCupService) matchesForDate(ctx context.Context, target time.Time, limit int) ([]domain.MatchView, error) {
	all, err := s.GetAllMatches(ctx)
	if err != nil {
		return nil, err
	}

	var selected []domain.MatchView
	for _, match := range all {
		if sameDate(match.Kickoff, target) {
			selected = append(selected, match)
		}
	}

	if len(selected) == 0 {
		now := s.now()
		for _, match := range all {
			if !match.Kickoff.Before(now) && match.Status != domain.StatusFinished {
				selected = append(selected, match)
			}
			if len(selected) >= limit {
				break
			}
		}
	}

	if len(selected) > limit {
		selected = selected[:limit]
	}
	return selected, nil
}

func (s *WorldCupService) buildMatchView(ctx context.Context, match domain.Match) (domain.MatchView, error) {
	stadiumName := "TBD"
	if match.StadiumID != "" && match.StadiumID != "0" {
		stadium, err := s.stadiums.ByID(ctx, match.StadiumID)
		if err != nil {
			return domain.MatchView{}, err
		}
		stadiumName = stadium.DisplayName()
	}

	homeTeam, err := s.resolveTeamName(ctx, match.HomeTeamID, match.HomeTeamLabel)
	if err != nil {
		return domain.MatchView{}, err
	}
	awayTeam, err := s.resolveTeamName(ctx, match.AwayTeamID, match.AwayTeamLabel)
	if err != nil {
		return domain.MatchView{}, err
	}

	status := domain.StatusUpcoming
	if match.Finished {
		status = domain.StatusFinished
	}
	if match.TimeElapsed != "" && match.TimeElapsed != "notstarted" && !match.Finished {
		status = domain.StatusLive
	}

	result := ""
	if status == domain.StatusFinished {
		switch {
		case match.HomeScore > match.AwayScore:
			result = fmt.Sprintf("%s won %d-%d", homeTeam, match.HomeScore, match.AwayScore)
		case match.AwayScore > match.HomeScore:
			result = fmt.Sprintf("%s won %d-%d", awayTeam, match.AwayScore, match.HomeScore)
		default:
			result = fmt.Sprintf("Draw %d-%d", match.HomeScore, match.AwayScore)
		}
	}

	return domain.MatchView{
		ID:          match.ID,
		MatchNumber: match.MatchNumber,
		Status:      status,
		HomeTeam:    homeTeam,
		AwayTeam:    awayTeam,
		HomeScore:   match.HomeScore,
		AwayScore:   match.AwayScore,
		Venue:       stadiumName,
		Group:       match.Group,
		Date:        match.Kickoff.Format("02 Jan"),
		Time:        match.Kickoff.Format("15:04"),
		Kickoff:     match.Kickoff,
		Result:      result,
		Scorers:     scorersFromMatch(homeTeam, awayTeam, match),
	}, nil
}

func (s *WorldCupService) resolveTeamName(ctx context.Context, id, label string) (string, error) {
	if id == "" || id == "0" {
		if label != "" {
			return label, nil
		}
		return "TBD", nil
	}
	team, err := s.teams.ByID(ctx, id)
	if err != nil {
		return "", err
	}
	return team.Name, nil
}

func scorersFromMatch(homeTeam, awayTeam string, match domain.Match) []domain.Scorer {
	var scorers []domain.Scorer
	for _, name := range match.HomeScorers {
		scorers = append(scorers, domain.Scorer{Team: homeTeam, Name: name})
	}
	for _, name := range match.AwayScorers {
		scorers = append(scorers, domain.Scorer{Team: awayTeam, Name: name})
	}
	return scorers
}

func sortMatchViews(matches []domain.MatchView) {
	sort.Slice(matches, func(i, j int) bool {
		if !matches[i].Kickoff.Equal(matches[j].Kickoff) {
			return matches[i].Kickoff.Before(matches[j].Kickoff)
		}
		return matches[i].MatchNumber < matches[j].MatchNumber
	})
}

func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
