package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	InitialLoadView = iota
	TabView

	LiveView
	MatchView
	StandingsView
	AboutView

	LastView
)

var tabOrder = []int{
	LiveView,
	MatchView,
	StandingsView,
	AboutView,
}

type Match struct {
	HomeTeam string
	AwayTeam string

	HomeScore int
	AwayScore int

	Status string
	Minute string
	Date   string
}

type Standing struct {
	Position int
	Team     string
	Played   int
	Won      int
	Drawn    int
	Lost     int
	Points   int
}

type Items struct {
	LiveMatch Match
	Matches   []Match
	Standings []Standing
}

type Model struct {
	currentView int
	selectedTab int

	width  int
	height int

	showLoadingCursor bool
	lastUpdated       time.Time

	renderer *lipgloss.Renderer
	styles   Styles

	items Items
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func NewModel(renderer *lipgloss.Renderer) Model {

	return Model{
		currentView: InitialLoadView,
		selectedTab: LiveView,

		renderer: renderer,
		styles:   NewStyles(renderer),

		lastUpdated: time.Now(),

		items: Items{
			LiveMatch: Match{
				HomeTeam:  "Argentina",
				AwayTeam:  "France",
				HomeScore: 2,
				AwayScore: 1,
				Status:    "LIVE",
				Minute:    "78'",
			},

			Matches: []Match{
				{
					HomeTeam: "Brazil",
					AwayTeam: "Germany",
					Date:     "12 Jun",
					Status:   "Upcoming",
				},
				{
					HomeTeam: "Spain",
					AwayTeam: "Japan",
					Date:     "12 Jun",
					Status:   "Upcoming",
				},
				{
					HomeTeam: "England",
					AwayTeam: "USA",
					Date:     "13 Jun",
					Status:   "Upcoming",
				},
			},

			Standings: []Standing{
				{
					Position: 1,
					Team:     "Argentina",
					Played:   3,
					Won:      2,
					Drawn:    1,
					Lost:     0,
					Points:   7,
				},
				{
					Position: 2,
					Team:     "Mexico",
					Played:   3,
					Won:      2,
					Drawn:    0,
					Lost:     1,
					Points:   6,
				},
				{
					Position: 3,
					Team:     "Poland",
					Played:   3,
					Won:      1,
					Drawn:    1,
					Lost:     1,
					Points:   4,
				},
				{
					Position: 4,
					Team:     "Saudi Arabia",
					Played:   3,
					Won:      0,
					Drawn:    0,
					Lost:     3,
					Points:   0,
				},
			},
		},
	}
}