package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func GetApiData[T any](route string) (T, error) {
	var result T

	apiBase := os.Getenv("API_URL")

	resp, err := http.Get(apiBase + "/" + route)
	if err != nil {
		return result, fmt.Errorf("error fetching %s: %w", route, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("error decoding %s: %w", route, err)
	}

	return result, nil
}

// ========================================
// Matches Screen
// ========================================

func GetMatchScores() (MatchScoresResponse, error) {
	return GetApiData[MatchScoresResponse]("matches")
}

// ========================================
// Live Screen
// ========================================

func GetLiveMatches() (LiveMatchResponse, error) {
	return GetApiData[LiveMatchResponse]("live")
}

// ========================================
// Points Screen
// ========================================

func GetPointsTable() (PointsTableResponse, error) {
	return GetApiData[PointsTableResponse]("standings")
}

// ========================================
// Schedule Screen
// ========================================

func GetMatchSchedule() (MatchScheduleResponse, error) {
	return GetApiData[MatchScheduleResponse]("schedule")
}

// ========================================
// Historical Winners (optional)
// ========================================

func GetHistoricalWinners() (HistoricalWinnersResponse, error) {
	return GetApiData[HistoricalWinnersResponse]("history")
}

func TeamToSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}