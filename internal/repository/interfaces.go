package repository

import (
	"context"

	"github.com/h0i5/ipl/internal/domain"
)

type TeamRepository interface {
	All(context.Context) ([]domain.Team, error)
	ByID(context.Context, string) (domain.Team, error)
}

type StadiumRepository interface {
	All(context.Context) ([]domain.Stadium, error)
	ByID(context.Context, string) (domain.Stadium, error)
}

type MatchRepository interface {
	All(context.Context) ([]domain.Match, error)
}

type StandingRepository interface {
	All(context.Context) ([]domain.Standing, error)
}

type WinnerRepository interface {
	All(context.Context) ([]domain.Winner, error)
}

type LiveMatchProvider interface {
	Current(context.Context) (domain.LiveMatch, error)
}
