package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/h0i5/ipl/internal/domain"
)

// CacheStandingRepository implements StandingRepository by reading from the
// standings cache file produced by the fetcher daemon.
type CacheStandingRepository struct {
	cachePath string
}

func NewCacheStandingRepository(cachePath string) *CacheStandingRepository {
	return &CacheStandingRepository{cachePath: cachePath}
}

func (r *CacheStandingRepository) All(ctx context.Context) ([]domain.Standing, error) {
	file, err := os.Open(r.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open standings cache: %w", err)
	}
	defer file.Close()

	var wrapper CacheWrapper[Standing]
	if err := json.NewDecoder(file).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode standings cache: %w", err)
	}

	results := make([]domain.Standing, 0, len(wrapper.Data))
	for _, s := range wrapper.Data {
		results = append(results, domain.Standing{
			Group:          s.Group,
			TeamID:         s.Team, // We store name here; service resolves by name
			Played:         s.Played,
			Won:            s.Won,
			Drawn:          s.Drawn,
			Lost:           s.Lost,
			GoalsFor:       s.GoalsFor,
			GoalsAgainst:   s.GoalsAgainst,
			GoalDifference: s.GoalDifference,
			Points:         s.Points,
		})
	}

	return results, nil
}
