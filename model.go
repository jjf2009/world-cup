package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
	"github.com/h0i5/ipl/cmd"
)

// all the views im planning to add rn
const (
	InitialLoadView = iota
	TabView
	LiveView
	MatchView
	PointsTableView
	// SquadView
	ScheduleView
	// HistoricalView
	AboutView

	// LastView is always last, never rendered just for limits
	LastView
)

type Items struct {
	liveMatch         cmd.LiveMatchResponse
	matchScores       cmd.MatchScoresResponse
	matchSchedule     cmd.MatchScheduleResponse
	pointsTable       cmd.PointsTableResponse
	historicalWinners cmd.HistoricalWinnersResponse
}

type Model struct {
	currentView int

	// contains the loading state of each view, to be used in view.go
	loadingMap        map[int]bool
	showLoadingCursor bool

	// pty dimensions for rendering
	width  int
	height int

	// this will be where the data lives
	items Items

	// tab view
	selectedTab int

	// match view table
	matchTable       table.Model
	matchTableStyles table.Styles

	// for themes to work correctly over ssh
	renderer *lipgloss.Renderer
	styles   Styles
}

func fetchCmd[T any](fetcher func() (T, error)) tea.Cmd {
	return func() tea.Msg {
		data, err := fetcher()
		if err != nil {
			return err
		}
		return data // T itself is the msg
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		fetchCmd[cmd.MatchScoresResponse](cmd.GetMatchScores),
		fetchCmd[cmd.MatchScheduleResponse](cmd.GetMatchSchedule),
		fetchCmd[cmd.PointsTableResponse](cmd.GetPointsTable),
		fetchCmd[cmd.LiveMatchResponse](cmd.GetLiveMatchScores),
	)
}

func NewModel(renderer *lipgloss.Renderer) Model {

	cols := []table.Column{
		{Title: "Status", Width: 7},
		{Title: "Teams", Width: 14},
		{Title: "Score 1", Width: 16},
		{Title: "Score 2", Width: 16},
		{Title: "Venue", Width: 12},
		{Title: "Result", Width: 26},
		{Title: "Date", Width: 7},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorRule).
		BorderBottom(true).
		Foreground(colorGold).
		Bold(true)
	ts.Selected = ts.Selected.
		Foreground(colorInk).
		Background(colorSienna).
		Bold(true)
	t.SetStyles(ts)

	loadingMap := map[int]bool{
		MatchView:       true,
		ScheduleView:    true,
		PointsTableView: true,
		LiveView:        true,
	}

	return Model{
		currentView:      InitialLoadView,
		loadingMap:       loadingMap,
		renderer:         renderer,
		styles:           NewStyles(renderer),
		selectedTab:      LiveView,
		matchTable:       t,
		matchTableStyles: ts,
	}
}
