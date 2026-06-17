package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/h0i5/ipl/internal/domain"
	"github.com/h0i5/ipl/internal/espn"
	"github.com/h0i5/ipl/internal/repository"
	"github.com/joho/godotenv"
)

type Config struct {
	LiveCachePath      string
	FixturesCachePath  string
	StandingsCachePath string
	TeamsPath          string
	LogPath            string
}

func main() {
	_ = godotenv.Load()

	config := Config{
		LiveCachePath:      getEnv("LIVE_CACHE_PATH", "cache/live_matches.json"),
		FixturesCachePath:  getEnv("FIXTURES_CACHE_PATH", "cache/fixtures.json"),
		StandingsCachePath: getEnv("STANDINGS_CACHE_PATH", "cache/standings.json"),
		TeamsPath:          getEnv("TEAMS_PATH", "data/football.teams.json"),
		LogPath:            getEnv("LOG_PATH", "logs/fetcher.log"),
	}

	// Ensure directories exist
	_ = os.MkdirAll(filepath.Dir(config.LiveCachePath), 0755)
	_ = os.MkdirAll(filepath.Dir(config.LogPath), 0755)

	// Setup logging
	logFile, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(mw, "[fetcher] ", log.LstdFlags)

	logger.Println("starting World Cup fetcher daemon")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Printf("received signal %v, shutting down gracefully", sig)
		cancel()
	}()

	// Load Teams Reference
	teamsRepo, err := repository.NewJSONTeamRepository(config.TeamsPath)
	if err != nil {
		logger.Fatalf("failed to load team repository: %v", err)
	}
	teams, err := teamsRepo.All(ctx)
	if err != nil {
		logger.Fatalf("failed to list teams: %v", err)
	}

	// Build mapping lookups
	abbrevMap := make(map[string]domain.Team)
	nameMap := make(map[string]domain.Team)
	for _, t := range teams {
		abbrevMap[strings.ToUpper(t.FIFACode)] = t
		nameMap[strings.ToUpper(t.Name)] = t
	}

	client := espn.NewClient(10 * time.Second)

	for {
		start := time.Now()
		logger.Println("starting update cycle")

		hasLiveMatches, err := runUpdate(ctx, client, config, abbrevMap, nameMap, logger)
		duration := time.Since(start)

		if err != nil {
			logger.Printf("update cycle completed with errors in %v: %v", duration, err)
		} else {
			logger.Printf("update cycle completed successfully in %v", duration)
		}

		sleepDuration := 5 * time.Second
		if hasLiveMatches {
			logger.Println("live match active: sleeping for 5 seconds")
		} else {
			logger.Println("no live matches: sleeping for 5 seconds")
		}

		select {
		case <-ctx.Done():
			logger.Println("fetcher loop stopped")
			return
		case <-time.After(sleepDuration):
		}
	}
}

func runUpdate(
	ctx context.Context,
	client *espn.Client,
	config Config,
	abbrevMap map[string]domain.Team,
	nameMap map[string]domain.Team,
	logger *log.Logger,
) (bool, error) {
	var errs []error

	// 1. Fetch Scoreboard
	var hasLiveMatches bool
	scoreboard, err := client.FetchScoreboard(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to fetch scoreboard: %w", err))
	} else {
		// Process Live Matches & Fixtures
		var liveMatches []repository.LiveMatch
		var fixtures []repository.Fixture

		for _, event := range scoreboard.Events {
			if len(event.Competitions) == 0 {
				continue
			}
			comp := event.Competitions[0]
			if len(comp.Competitors) < 2 {
				continue
			}

			// Map competitors
			var homeComp, awayComp espn.Competitor
			if comp.Competitors[0].HomeAway == "home" {
				homeComp = comp.Competitors[0]
				awayComp = comp.Competitors[1]
			} else {
				homeComp = comp.Competitors[1]
				awayComp = comp.Competitors[0]
			}

			// Resolve team names
			homeTeamObj, homeOk := resolveTeam(homeComp.Team.DisplayName, homeComp.Team.Abbreviation, abbrevMap, nameMap)
			awayTeamObj, awayOk := resolveTeam(awayComp.Team.DisplayName, awayComp.Team.Abbreviation, abbrevMap, nameMap)

			homeName := homeComp.Team.DisplayName
			if homeOk {
				homeName = homeTeamObj.Name
			}
			awayName := awayComp.Team.DisplayName
			if awayOk {
				awayName = awayTeamObj.Name
			}

			homeScore, _ := strconv.Atoi(homeComp.Score)
			awayScore, _ := strconv.Atoi(awayComp.Score)

			venueName := ""
			if comp.Venue != nil {
				venueName = comp.Venue.FullName
			}

			state := event.Status.Type.State // "pre", "in", "post"
			statusText := event.Status.Type.Description
			if state == "in" {
				statusText = event.Status.DisplayClock
			}

			// 1.a Fixture cache item
			fixtures = append(fixtures, repository.Fixture{
				ID:          event.ID,
				HomeTeam:    homeName,
				AwayTeam:    awayName,
				KickoffTime: event.Date,
				Status:      state,
				Venue:       venueName,
				HomeScore:   homeScore,
				AwayScore:   awayScore,
			})

			// 1.b Live match cache item
			if state == "in" {
				hasLiveMatches = true

				group := ""
				if homeOk && awayOk && homeTeamObj.Group == awayTeamObj.Group {
					group = homeTeamObj.Group
				}

				liveMatches = append(liveMatches, repository.LiveMatch{
					ID:        event.ID,
					HomeTeam:  homeName,
					AwayTeam:  awayName,
					HomeScore: homeScore,
					AwayScore: awayScore,
					Minute:    event.Status.DisplayClock,
					Status:    statusText,
					Venue:     venueName,
					Group:     group,
				})
			}
		}

		// Write Live Matches Cache
		liveWrapper := repository.CacheWrapper[repository.LiveMatch]{
			UpdatedAt: time.Now(),
			Data:      liveMatches,
		}
		if err := writeAtomic(config.LiveCachePath, liveWrapper); err != nil {
			errs = append(errs, fmt.Errorf("failed to write live matches cache: %w", err))
		} else {
			logger.Printf("cached %d live matches to %s", len(liveMatches), config.LiveCachePath)
		}

		// Write Fixtures Cache
		fixturesWrapper := repository.CacheWrapper[repository.Fixture]{
			UpdatedAt: time.Now(),
			Data:      fixtures,
		}
		if err := writeAtomic(config.FixturesCachePath, fixturesWrapper); err != nil {
			errs = append(errs, fmt.Errorf("failed to write fixtures cache: %w", err))
		} else {
			logger.Printf("cached %d fixtures to %s", len(fixtures), config.FixturesCachePath)
		}
	}

	// 2. Fetch Standings
	standings, err := client.FetchStandings(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to fetch standings: %w", err))
	} else {
		var standingsList []repository.Standing

		for _, child := range standings.Children {
			groupName := strings.TrimPrefix(child.Name, "FIFA World Cup ") // Clean group names
			groupName = strings.TrimSpace(groupName)

			for _, entry := range child.Standings.Entries {
				teamObj, ok := resolveTeam(entry.Team.DisplayName, entry.Team.Abbreviation, abbrevMap, nameMap)
				teamName := entry.Team.DisplayName
				if ok {
					teamName = teamObj.Name
				}

				// Extract stats
				statsMap := make(map[string]int)
				for _, stat := range entry.Stats {
					statsMap[stat.Name] = int(stat.Value)
				}

				standingsList = append(standingsList, repository.Standing{
					Team:           teamName,
					Played:         statsMap["gamesPlayed"],
					Won:            statsMap["wins"],
					Drawn:          statsMap["ties"],
					Lost:           statsMap["losses"],
					GoalsFor:       statsMap["pointsFor"],
					GoalsAgainst:   statsMap["pointsAgainst"],
					GoalDifference: statsMap["pointDifferential"],
					Points:         statsMap["points"],
					Group:          groupName,
				})
			}
		}

		standingsWrapper := repository.CacheWrapper[repository.Standing]{
			UpdatedAt: time.Now(),
			Data:      standingsList,
		}
		if err := writeAtomic(config.StandingsCachePath, standingsWrapper); err != nil {
			errs = append(errs, fmt.Errorf("failed to write standings cache: %w", err))
		} else {
			logger.Printf("cached %d standings to %s", len(standingsList), config.StandingsCachePath)
		}
	}

	if len(errs) > 0 {
		return hasLiveMatches, errors.Join(errs...)
	}

	return hasLiveMatches, nil
}

func resolveTeam(
	displayName, abbreviation string,
	abbrevMap map[string]domain.Team,
	nameMap map[string]domain.Team,
) (domain.Team, bool) {
	if abbreviation != "" {
		if t, ok := abbrevMap[strings.ToUpper(abbreviation)]; ok {
			return t, true
		}
	}
	if displayName != "" {
		if t, ok := nameMap[strings.ToUpper(displayName)]; ok {
			return t, true
		}
		cleanDisp := cleanName(displayName)
		for name, t := range nameMap {
			if cleanName(name) == cleanDisp {
				return t, true
			}
		}
	}
	return domain.Team{}, false
}

func cleanName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "and", "")
	s = strings.ReplaceAll(s, "the", "")
	s = strings.ReplaceAll(s, "republicof", "")
	s = strings.ReplaceAll(s, "democratic", "")
	s = strings.ReplaceAll(s, "drcongo", "congo")
	s = strings.ReplaceAll(s, "dem.rep.", "congo")
	return s
}

func writeAtomic[T any](path string, data repository.CacheWrapper[T]) error {
	tempFile := path + ".tmp"

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := os.WriteFile(tempFile, bytes, 0644); err != nil {
		return fmt.Errorf("write temp failed: %w", err)
	}

	tempBytes, err := os.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("read temp failed: %w", err)
	}
	var dummy repository.CacheWrapper[T]
	if err := json.Unmarshal(tempBytes, &dummy); err != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("validate temp failed: JSON is invalid: %w", err)
	}

	if err := os.Rename(tempFile, path); err != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("rename failed: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
