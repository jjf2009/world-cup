package main

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/h0i5/ipl/internal/repository"
	"github.com/h0i5/ipl/internal/service"
	"github.com/joho/godotenv"
)

const (
	host = "0.0.0.0"
	port = "6767"
)

func myLoggingMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			start := time.Now()
			slog.Info("connect", "ip", s.RemoteAddr(), "user", s.User())
			next(s)
			slog.Info("disconnect", "ip", s.RemoteAddr(), "user", s.User(), "duration", time.Since(start))
		}
	}
}

func main() {
	godotenv.Load()

	// Ensure logs directory exists
	_ = os.MkdirAll("logs", 0755)

	// log to a file
	logFile, _ := os.OpenFile("logs/server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multi)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	worldCupService, err := buildWorldCupService(appCtx)
	if err != nil {
		log.Fatalf("could not initialize world cup data: %v", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		// Auto-generates host key at this path if missing
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler(worldCupService)),
			myLoggingMiddleware(),
		),
	)
	if err != nil {
		log.Fatalf("could not create server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting SSH server on %s:%s", host, port)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	appCancel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}

func buildWorldCupService(_ context.Context) (*service.WorldCupService, error) {
	// Static data repositories (teams, stadiums, winners — do not change)
	teams, err := repository.NewJSONTeamRepository("data/football.teams.json")
	if err != nil {
		return nil, err
	}
	stadiums, err := repository.NewJSONStadiumRepository("data/football.stadiums.json")
	if err != nil {
		return nil, err
	}
	winners, err := repository.NewJSONWinnerRepository("data/winners.json")
	if err != nil {
		return nil, err
	}

	// Live data repositories — read from fetcher daemon's cache files
	matches := repository.NewCacheFixtureRepository("cache/fixtures.json")
	standings := repository.NewCacheStandingRepository("cache/standings.json")
	live := repository.NewCacheLiveRepository("cache/live_matches.json")

	return service.NewWorldCupService(teams, stadiums, matches, standings, winners, live), nil
}

// teaHandler is called once per SSH session and returns a fresh model.
func teaHandler(worldCupService WorldCupDataService) func(ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		renderer := bubbletea.MakeRenderer(s)

		m := NewModel(renderer, worldCupService)
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
