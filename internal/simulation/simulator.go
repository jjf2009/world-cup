package simulation

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/h0i5/ipl/internal/domain"
	"github.com/h0i5/ipl/internal/repository"
)

type Simulator struct {
	matches  repository.MatchRepository
	teams    repository.TeamRepository
	stadiums repository.StadiumRepository
	now      func() time.Time
	rand     *rand.Rand

	mu      sync.RWMutex
	current domain.LiveMatch
	used    map[string]bool
}

func NewSimulator(
	matches repository.MatchRepository,
	teams repository.TeamRepository,
	stadiums repository.StadiumRepository,
) *Simulator {
	return &Simulator{
		matches:  matches,
		teams:    teams,
		stadiums: stadiums,
		now:      time.Now,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		used:     make(map[string]bool),
	}
}

func (s *Simulator) Start(ctx context.Context) error {
	if err := s.selectNext(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.tick(ctx)
			}
		}
	}()

	return nil
}

func (s *Simulator) Current(context.Context) (domain.LiveMatch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.current.ID == "" {
		return domain.LiveMatch{}, fmt.Errorf("live match: %w", domain.ErrNotFound)
	}

	match := s.current
	match.Scorers = append([]domain.Scorer(nil), s.current.Scorers...)
	return match, nil
}

func (s *Simulator) tick(ctx context.Context) {
	s.mu.Lock()

	if s.current.ID == "" {
		s.mu.Unlock()
		_ = s.selectNext(ctx)
		return
	}

	if s.current.Status == domain.StatusUpcoming {
		s.current.Status = domain.StatusLive
		s.current.Minute = 1
		s.current.LastUpdated = s.now()
		s.mu.Unlock()
		return
	}

	if s.current.Status == domain.StatusFinished {
		s.mu.Unlock()
		_ = s.selectNext(ctx)
		return
	}

	s.current.Minute++
	if s.current.Minute >= 90 {
		s.current.Minute = 90
		s.current.Status = domain.StatusFinished
		s.current.LastUpdated = s.now()
		s.used[s.current.ID] = true
		s.mu.Unlock()
		_ = s.selectNext(ctx)
		return
	}

	s.maybeScoreLocked()
	s.current.LastUpdated = s.now()
	s.mu.Unlock()
}

func (s *Simulator) maybeScoreLocked() {
	totalGoals := s.current.HomeScore + s.current.AwayScore
	if totalGoals >= 6 {
		return
	}

	minute := s.current.Minute
	probability := 0.025
	if minute >= 75 {
		probability = 0.035
	}
	if s.rand.Float64() > probability {
		return
	}

	homeGoal := s.rand.Intn(2) == 0
	team := s.current.HomeTeam
	if homeGoal {
		s.current.HomeScore++
	} else {
		s.current.AwayScore++
		team = s.current.AwayTeam
	}

	s.current.Scorers = append(s.current.Scorers, domain.Scorer{
		Team:   team,
		Name:   fmt.Sprintf("%s scorer %d", team, countTeamScorers(s.current.Scorers, team)+1),
		Minute: minute,
	})
}

func (s *Simulator) selectNext(ctx context.Context) error {
	matches, err := s.matches.All(ctx)
	if err != nil {
		return err
	}

	now := s.now()
	var selected *domain.Match
	for i := range matches {
		match := matches[i]
		if s.used[match.ID] || match.Finished {
			continue
		}
		if sameDate(match.Kickoff, now) {
			selected = &match
			break
		}
	}
	if selected == nil {
		for i := range matches {
			match := matches[i]
			if s.used[match.ID] || match.Finished || match.Kickoff.Before(now) {
				continue
			}
			selected = &match
			break
		}
	}
	if selected == nil {
		for i := range matches {
			match := matches[i]
			if !s.used[match.ID] && !match.Finished {
				selected = &match
				break
			}
		}
	}
	if selected == nil {
		return fmt.Errorf("select simulated match: %w", domain.ErrNotFound)
	}

	live, err := s.toLiveMatch(ctx, *selected)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.current = live
	s.mu.Unlock()
	return nil
}

func (s *Simulator) toLiveMatch(ctx context.Context, match domain.Match) (domain.LiveMatch, error) {
	home, err := s.teamName(ctx, match.HomeTeamID, match.HomeTeamLabel)
	if err != nil {
		return domain.LiveMatch{}, err
	}
	away, err := s.teamName(ctx, match.AwayTeamID, match.AwayTeamLabel)
	if err != nil {
		return domain.LiveMatch{}, err
	}
	stadium, err := s.stadiums.ByID(ctx, match.StadiumID)
	if err != nil {
		return domain.LiveMatch{}, err
	}

	return domain.LiveMatch{
		ID:          match.ID,
		MatchNumber: match.MatchNumber,
		HomeTeam:    home,
		AwayTeam:    away,
		HomeScore:   match.HomeScore,
		AwayScore:   match.AwayScore,
		Minute:      1,
		Status:      domain.StatusLive,
		Stadium:     stadium.DisplayName(),
		Group:       match.Group,
		Scorers:     nil,
		LastUpdated: s.now(),
	}, nil
}

func (s *Simulator) teamName(ctx context.Context, id, label string) (string, error) {
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

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func countTeamScorers(scorers []domain.Scorer, team string) int {
	count := 0
	for _, scorer := range scorers {
		if scorer.Team == team {
			count++
		}
	}
	return count
}
