package main

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/h0i5/ipl/cmd"
)

const (
	InitialLoadView = iota
	TabView

	LiveView
	MatchView
	PointsTableView
	ScheduleView
	HistoricalView
	AboutView

	LastView
)

var tabOrder = []int{
	LiveView,
	MatchView,
	PointsTableView,
	ScheduleView,
	HistoricalView,
	AboutView,
}

type Items struct {
	liveMatch         cmd.LiveMatchResponse
	matches           []cmd.Match
	pointsTable       cmd.PointsTableResponse
	matchSchedule     cmd.MatchScheduleResponse
	historicalWinners cmd.HistoricalWinnersResponse
	squads            map[string]cmd.SquadResponse
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

	items Items
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), liveTickCmd())
}

func NewModel(renderer *lipgloss.Renderer) Model {
	// Columns for the matches table
	columns := []table.Column{
		{Title: "Match #", Width: 7},
		{Title: "Home Team", Width: 14},
		{Title: "Score", Width: 7},
		{Title: "Away Team", Width: 14},
		{Title: "Status", Width: 9},
		{Title: "Date", Width: 8},
		{Title: "Venue", Width: 15},
	}

	// Rows for the matches table
	rows := []table.Row{
		{"64", "Argentina", "2 - 1", "France", "LIVE", "04 Jun", "Lusail Stadium"},
		{"1", "Qatar", "0 - 2", "Ecuador", "Completed", "20 Nov", "Al Bayt Stadium"},
		{"2", "England", "6 - 2", "Iran", "Completed", "21 Nov", "Khalifa Intl"},
		{"3", "Senegal", "0 - 0", "Netherlands", "Upcoming", "21 Nov", "Al Thumama"},
		{"4", "USA", "0 - 0", "Wales", "Upcoming", "21 Nov", "Ahmad Bin Ali"},
		{"5", "Brazil", "0 - 0", "Germany", "Upcoming", "12 Jun", "Maracanã"},
		{"6", "Spain", "0 - 0", "Japan", "Upcoming", "12 Jun", "Kalyani"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Set custom table styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("166")). // Sienna/Orange highlight matching theme
		Bold(true)
	t.SetStyles(s)

	styles := NewStyles(renderer)

	return Model{
		currentView:    InitialLoadView,
		selectedTab:    LiveView,
		showLiveCursor: true,

		loadingMap: map[int]bool{
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

		lastUpdated: time.Now(),

		items: Items{
			liveMatch: cmd.LiveMatchResponse{
				StatusCode: 200,
				LiveCount:  1,
				Matches: map[string]cmd.LiveMatch{
					"64": {
						ID:          "64",
						MatchNumber: "64",
						Status:      "LIVE",
						HomeTeam:    "Argentina",
						AwayTeam:    "France",
						HomeScore:   2,
						AwayScore:   1,
						Minute:      "78",
						Venue:       "Lusail Stadium",
					},
				},
				LastUpdated: "15:04:05",
			},

			matches: []cmd.Match{
				{ID: "64", MatchNumber: "64", Status: "LIVE", HomeTeam: "Argentina", AwayTeam: "France", HomeScore: 2, AwayScore: 1, Venue: "Lusail Stadium", Date: "04 Jun"},
				{ID: "1", MatchNumber: "1", Status: "Completed", HomeTeam: "Qatar", AwayTeam: "Ecuador", HomeScore: 0, AwayScore: 2, Venue: "Al Bayt Stadium", Date: "20 Nov", Result: "Ecuador won 2-0"},
				{ID: "2", MatchNumber: "2", Status: "Completed", HomeTeam: "England", AwayTeam: "Iran", HomeScore: 6, AwayScore: 2, Venue: "Khalifa Intl", Date: "21 Nov", Result: "England won 6-2"},
				{ID: "3", MatchNumber: "3", Status: "Upcoming", HomeTeam: "Senegal", AwayTeam: "Netherlands", Venue: "Al Thumama", Date: "21 Nov"},
				{ID: "4", MatchNumber: "4", Status: "Upcoming", HomeTeam: "USA", AwayTeam: "Wales", Venue: "Ahmad Bin Ali", Date: "21 Nov"},
				{ID: "5", MatchNumber: "5", Status: "Upcoming", HomeTeam: "Brazil", AwayTeam: "Germany", Venue: "Maracanã", Date: "12 Jun"},
				{ID: "6", MatchNumber: "6", Status: "Upcoming", HomeTeam: "Spain", AwayTeam: "Japan", Venue: "Kalyani", Date: "12 Jun"},
			},

			pointsTable: cmd.PointsTableResponse{
				StatusCode: 200,
				Standings: map[string]cmd.TeamStanding{
					"ARG": {Position: 1, Team: "Argentina", Played: 3, Won: 2, Drawn: 1, Lost: 0, GoalsFor: 5, GoalsAgainst: 2, GoalDifference: 3, Points: 7},
					"POL": {Position: 2, Team: "Poland", Played: 3, Won: 1, Drawn: 1, Lost: 1, GoalsFor: 2, GoalsAgainst: 2, GoalDifference: 0, Points: 4},
					"MEX": {Position: 3, Team: "Mexico", Played: 3, Won: 1, Drawn: 1, Lost: 1, GoalsFor: 2, GoalsAgainst: 3, GoalDifference: -1, Points: 4},
					"KSA": {Position: 4, Team: "Saudi Arabia", Played: 3, Won: 1, Drawn: 0, Lost: 2, GoalsFor: 3, GoalsAgainst: 5, GoalDifference: -2, Points: 3},
				},
			},

			matchSchedule: cmd.MatchScheduleResponse{
				StatusCode: 200,
				Schedule: map[string]cmd.ScheduledMatch{
					"1": {MatchNumber: "1", HomeTeam: "Qatar", AwayTeam: "Ecuador", Venue: "Al Bayt Stadium", Date: "20 Nov 2022", Time: "19:00"},
					"2": {MatchNumber: "2", HomeTeam: "England", AwayTeam: "Iran", Venue: "Khalifa Intl", Date: "21 Nov 2022", Time: "16:00"},
					"3": {MatchNumber: "3", HomeTeam: "Senegal", AwayTeam: "Netherlands", Venue: "Al Thumama", Date: "21 Nov 2022", Time: "19:00"},
					"4": {MatchNumber: "4", HomeTeam: "USA", AwayTeam: "Wales", Venue: "Ahmad Bin Ali", Date: "21 Nov 2022", Time: "22:00"},
				},
			},

			historicalWinners: cmd.HistoricalWinnersResponse{
				StatusCode: 200,
				Winners: map[string]cmd.YearWinner{
					"2022": {Winner: "Argentina", RunnerUp: "France", Venue: "Lusail Stadium, Qatar"},
					"2018": {Winner: "France", RunnerUp: "Croatia", Venue: "Luzhniki Stadium, Russia"},
					"2014": {Winner: "Germany", RunnerUp: "Argentina", Venue: "Maracanã, Brazil"},
					"2010": {Winner: "Spain", RunnerUp: "Netherlands", Venue: "Soccer City, South Africa"},
				},
			},

			squads: map[string]cmd.SquadResponse{
				"argentina": {
					StatusCode: 200,
					Team:       "Argentina",
					Squad: map[string]cmd.SquadPlayer{
						"1": {Name: "E. Martínez", Nationality: "ARG", Position: "GK", JerseyNumber: 23},
						"2": {Name: "N. Otamendi", Nationality: "ARG", Position: "DEF", JerseyNumber: 19},
						"3": {Name: "C. Romero", Nationality: "ARG", Position: "DEF", JerseyNumber: 13},
						"4": {Name: "R. De Paul", Nationality: "ARG", Position: "MID", JerseyNumber: 7},
						"5": {Name: "E. Fernández", Nationality: "ARG", Position: "MID", JerseyNumber: 24},
						"6": {Name: "A. Mac Allister", Nationality: "ARG", Position: "MID", JerseyNumber: 20},
						"7": {Name: "L. Messi", Nationality: "ARG", Position: "FWD", JerseyNumber: 10},
						"8": {Name: "J. Álvarez", Nationality: "ARG", Position: "FWD", JerseyNumber: 9},
					},
				},
				"france": {
					StatusCode: 200,
					Team:       "France",
					Squad: map[string]cmd.SquadPlayer{
						"1": {Name: "H. Lloris", Nationality: "FRA", Position: "GK", JerseyNumber: 1},
						"2": {Name: "R. Varane", Nationality: "FRA", Position: "DEF", JerseyNumber: 4},
						"3": {Name: "D. Upamecano", Nationality: "FRA", Position: "DEF", JerseyNumber: 18},
						"4": {Name: "T. Hernandez", Nationality: "FRA", Position: "DEF", JerseyNumber: 22},
						"5": {Name: "A. Tchouaméni", Nationality: "FRA", Position: "MID", JerseyNumber: 8},
						"6": {Name: "A. Rabiot", Nationality: "FRA", Position: "MID", JerseyNumber: 14},
						"7": {Name: "A. Griezmann", Nationality: "FRA", Position: "MID", JerseyNumber: 7},
						"8": {Name: "K. Mbappé", Nationality: "FRA", Position: "FWD", JerseyNumber: 10},
					},
				},
			},
		},
	}
}