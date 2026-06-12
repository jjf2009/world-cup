package espn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultScoreboardURL = "https://site.api.espn.com/apis/site/v2/sports/soccer/fifa.world/scoreboard?limit=200"
	DefaultStandingsURL  = "https://site.api.espn.com/apis/v2/sports/soccer/fifa.world/standings"
)

type Client struct {
	httpClient    *http.Client
	scoreboardURL string
	standingsURL  string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		scoreboardURL: DefaultScoreboardURL,
		standingsURL:  DefaultStandingsURL,
	}
}

// WithScoreboardURL allows overriding the scoreboard API URL (useful for testing/mocking).
func (c *Client) WithScoreboardURL(url string) *Client {
	c.scoreboardURL = url
	return c
}

// WithStandingsURL allows overriding the standings API URL (useful for testing/mocking).
func (c *Client) WithStandingsURL(url string) *Client {
	c.standingsURL = url
	return c
}

func (c *Client) FetchScoreboard(ctx context.Context) (*ScoreboardResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.scoreboardURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res ScoreboardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &res, nil
}

func (c *Client) FetchStandings(ctx context.Context) (*StandingsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.standingsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res StandingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &res, nil
}
