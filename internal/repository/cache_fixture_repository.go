package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/h0i5/ipl/internal/domain"
)

// espnDateFormats lists the date layouts ESPN may return.
// The primary format omits seconds: "2006-01-02T15:04Z"
var espnDateFormats = []string{
	"2006-01-02T15:04Z",
	time.RFC3339,
	"2006-01-02T15:04:05Z",
}

func parseESPNDate(s string) time.Time {
	for _, layout := range espnDateFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// CacheFixtureRepository implements MatchRepository by reading from the
// fixtures cache file produced by the fetcher daemon.
type CacheFixtureRepository struct {
	cachePath string
}

func NewCacheFixtureRepository(cachePath string) *CacheFixtureRepository {
	return &CacheFixtureRepository{cachePath: cachePath}
}

func (r *CacheFixtureRepository) All(ctx context.Context) ([]domain.Match, error) {
	file, err := os.Open(r.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Cache not yet populated — return empty slice so TUI shows no data
			return nil, nil
		}
		return nil, fmt.Errorf("open fixtures cache: %w", err)
	}
	defer file.Close()

	var wrapper CacheWrapper[Fixture]
	if err := json.NewDecoder(file).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode fixtures cache: %w", err)
	}

	matches := make([]domain.Match, 0, len(wrapper.Data))
	for i, f := range wrapper.Data {
		// Map ESPN status string to domain status
		var state string
		switch f.Status {
		case "in":
			state = domain.StatusLive
		case "post":
			state = domain.StatusFinished
		default:
			state = domain.StatusUpcoming
		}

		// Determine finished flag
		finished := f.Status == "post"

		// Use index as a stable match number when no other source is available
		matchNumber := fmt.Sprintf("%d", i+1)

		kickoff := parseESPNDate(f.KickoffTime)

		matches = append(matches, domain.Match{
			ID:          f.ID,
			MatchNumber: matchNumber,
			HomeTeamID:  "",
			AwayTeamID:  "",
			// Labels are used when IDs are unknown (our cache-driven path)
			HomeTeamLabel: f.HomeTeam,
			AwayTeamLabel: f.AwayTeam,
			HomeScore:     f.HomeScore,
			AwayScore:     f.AwayScore,
			Kickoff:       kickoff,
			Finished:      finished,
			TimeElapsed:   state,
			Type:          "group",
		})
	}

	return matches, nil
}
