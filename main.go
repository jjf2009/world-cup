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

	// log to a file
	logFile, _ := os.OpenFile("logs/server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	multi := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multi)

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		// Auto-generates host key at this path if missing
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}

// teaHandler is called once per SSH session — return your model here.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	//pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	m := NewModel(renderer)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
