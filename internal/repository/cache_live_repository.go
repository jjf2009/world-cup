package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/h0i5/ipl/internal/domain"
)

type CacheLiveRepository struct {
	cachePath string
}

func NewCacheLiveRepository(cachePath string) *CacheLiveRepository {
	return &CacheLiveRepository{cachePath: cachePath}
}

func (r *CacheLiveRepository) Current(ctx context.Context) (domain.LiveMatch, error) {
	file, err := os.Open(r.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.LiveMatch{}, fmt.Errorf("live match: %w", domain.ErrNotFound)
		}
		return domain.LiveMatch{}, fmt.Errorf("failed to open live cache: %w", err)
	}
	defer file.Close()

	var wrapper CacheWrapper[LiveMatch]
	if err := json.NewDecoder(file).Decode(&wrapper); err != nil {
		return domain.LiveMatch{}, fmt.Errorf("failed to decode live cache: %w", err)
	}

	if len(wrapper.Data) == 0 {
		return domain.LiveMatch{}, fmt.Errorf("live match: %w", domain.ErrNotFound)
	}

	lm := wrapper.Data[0]

	minuteVal := 0
	cleanMin := strings.TrimSuffix(lm.Minute, "'")
	if val, err := strconv.Atoi(cleanMin); err == nil {
		minuteVal = val
	}

	return domain.LiveMatch{
		ID:          lm.ID,
		MatchNumber: lm.ID,
		HomeTeam:    lm.HomeTeam,
		AwayTeam:    lm.AwayTeam,
		HomeScore:   lm.HomeScore,
		AwayScore:   lm.AwayScore,
		Minute:      minuteVal,
		Status:      lm.Status,
		Stadium:     lm.Venue,
		Group:       lm.Group,
		Scorers:     nil,
		LastUpdated: wrapper.UpdatedAt,
	}, nil
}
