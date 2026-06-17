package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/h0i5/ipl/cmd"
	"github.com/h0i5/ipl/internal/domain"
)

const (
	InitialLoadView = iota
	TabView

	DashboardView
	LiveView
	MatchView
	PointsTableView
	ScheduleView
	HistoricalView
	AboutView

	LastView
)

var tabOrder = []int{
	DashboardView,
	LiveView,
	MatchView,
	PointsTableView,
	ScheduleView,
	HistoricalView,
	AboutView,
}

type Items struct {
	liveMatch         domain.LiveMatchView
	matches           []domain.MatchView
	standings         []domain.StandingView
	matchSchedule     domain.ScheduleView
	historicalWinners []domain.WinnerView
	squads            map[string]cmd.SquadResponse
}

type WorldCupDataService interface {
	GetLiveMatch(context.Context) (domain.LiveMatchView, error)
	GetTodayMatches(context.Context) ([]domain.MatchView, error)
	GetTomorrowMatches(context.Context) ([]domain.MatchView, error)
	GetStandings(context.Context) ([]domain.StandingView, error)
	GetHistoricalWinners(context.Context) ([]domain.WinnerView, error)
	GetAllMatches(context.Context) ([]domain.MatchView, error)
	GetTeamByID(context.Context, string) (domain.Team, error)
	GetStadiumByID(context.Context, string) (domain.Stadium, error)
}

type Model struct {
	currentView int
	selectedTab int

	width  int
	height int

	showLoadingCursor bool
	showLiveCursor    bool
	lastUpdated       time.Time

	loadingMap map[int]bool

	renderer *lipgloss.Renderer
	styles   Styles

	matchTable       table.Model
	matchTableStyles table.Styles

	// standingsVP is a scrollable viewport for the standings table.
	standingsVP    viewport.Model
	standingsReady bool

	service WorldCupDataService
	lastErr string
	items   Items
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), liveTickCmd())
}

func NewModel(renderer *lipgloss.Renderer, service WorldCupDataService) Model {
	columns := []table.Column{
		{Title: "Match #", Width: 7},
		{Title: "Home Team", Width: 14},
		{Title: "Score", Width: 7},
		{Title: "Away Team", Width: 14},
		{Title: "Status", Width: 9},
		{Title: "Date", Width: 8},
		{Title: "Venue", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(nil),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("166")).
		Bold(true)
	t.SetStyles(s)

	styles := NewStyles(renderer)

	// Initial viewport — dimensions will be set once WindowSizeMsg arrives.
	vp := viewport.New(80, 20)
	vp.YPosition = 0

	m := Model{
		currentView:    InitialLoadView,
		selectedTab:    LiveView,
		showLiveCursor: true,

		loadingMap: map[int]bool{
			DashboardView:   false,
			LiveView:        false,
			MatchView:       false,
			PointsTableView: false,
			ScheduleView:    false,
			HistoricalView:  false,
			AboutView:       false,
		},

		renderer:         renderer,
		styles:           styles,
		matchTable:       t,
		matchTableStyles: s,
		standingsVP:      vp,
		standingsReady:   false,
		service:          service,
		lastUpdated:      time.Now(),
		items: Items{
			squads: map[string]cmd.SquadResponse{},
		},
	}

	m.refreshData(context.Background())
	return m
}

func (m *Model) refreshData(ctx context.Context) {
	if m.service == nil {
		m.lastErr = "world cup service is not configured"
		return
	}

	liveMatch, err := m.service.GetLiveMatch(ctx)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		m.lastErr = err.Error()
		return
	}
	matches, err := m.service.GetAllMatches(ctx)
	if err != nil {
		m.lastErr = err.Error()
		return
	}
	today, err := m.service.GetTodayMatches(ctx)
	if err != nil {
		m.lastErr = err.Error()
		return
	}
	tomorrow, err := m.service.GetTomorrowMatches(ctx)
	if err != nil {
		m.lastErr = err.Error()
		return
	}
	standings, err := m.service.GetStandings(ctx)
	if err != nil {
		m.lastErr = err.Error()
		return
	}
	winners, err := m.service.GetHistoricalWinners(ctx)
	if err != nil {
		m.lastErr = err.Error()
		return
	}

	m.items.liveMatch = liveMatch
	m.items.matches = overlayLiveMatch(matches, liveMatch)
	m.items.matchSchedule = domain.ScheduleView{Today: today, Tomorrow: tomorrow}
	m.items.standings = standings
	m.items.historicalWinners = winners
	m.matchTable.SetRows(matchRows(m.items.matches))
	m.lastUpdated = time.Now()
	m.lastErr = ""
}

func overlayLiveMatch(matches []domain.MatchView, live domain.LiveMatchView) []domain.MatchView {
	updated := append([]domain.MatchView(nil), matches...)
	for i := range updated {
		if updated[i].ID != live.ID {
			continue
		}
		updated[i].Status = live.Status
		updated[i].HomeScore = live.HomeScore
		updated[i].AwayScore = live.AwayScore
		updated[i].Scorers = append([]domain.Scorer(nil), live.Scorers...)
		break
	}
	return updated
}

func matchRows(matches []domain.MatchView) []table.Row {
	rows := make([]table.Row, 0, len(matches))
	for _, match := range matches {
		rows = append(rows, table.Row{
			match.MatchNumber,
			truncate(match.HomeTeam, 14),
			fmt.Sprintf("%d - %d", match.HomeScore, match.AwayScore),
			truncate(match.AwayTeam, 14),
			match.Status,
			match.Date,
			truncateVenue(match.Venue, 15),
		})
	}
	return rows
}
