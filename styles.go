package main

import "github.com/charmbracelet/lipgloss"

// Palette
var (
	colorInk      = lipgloss.Color("#EAE4DA")
	colorInkMuted = lipgloss.Color("#A89F8E")
	colorRule     = lipgloss.Color("#2A2622")
	colorSienna   = lipgloss.Color("#C2410C")
	colorGold     = lipgloss.Color("#F0D9B5")
	colorDim      = lipgloss.Color("#6B6560")
	colorLive     = lipgloss.Color("#EF4444")
	colorGreen    = lipgloss.Color("#4ADE80")
)

type Styles struct {
	// Loading
	cursorStyle lipgloss.Style

	// Sidebar
	tabActive   lipgloss.Style
	tabInactive lipgloss.Style
	tabKey      lipgloss.Style
	tabKeyDim   lipgloss.Style
	sidebar     lipgloss.Style

	// Outer frame
	outerBorder lipgloss.Style

	// Body
	body   lipgloss.Style
	header lipgloss.Style

	// Match card
	teamName  lipgloss.Style
	score     lipgloss.Style
	liveDot   lipgloss.Style
	venue     lipgloss.Style
	result    lipgloss.Style
	matchCard lipgloss.Style

	// Points table
	tableHeader lipgloss.Style
	tableRow    lipgloss.Style
	tableRowAlt lipgloss.Style
	highlight   lipgloss.Style

	// General
	faint   lipgloss.Style
	muted   lipgloss.Style
	gold    lipgloss.Style
	loading lipgloss.Style
}

func NewStyles(r *lipgloss.Renderer) Styles {
	return Styles{
		// Loading cursor
		cursorStyle: r.NewStyle().Foreground(colorSienna),

		// Sidebar tabs
		tabActive: r.NewStyle().
			Foreground(colorInk).
			Padding(0, 2).
			Bold(true),

		tabInactive: r.NewStyle().
			Foreground(colorInkMuted).
			Padding(0, 2),

		tabKey: r.NewStyle().
			Foreground(colorInkMuted).
			PaddingRight(1),

		tabKeyDim: r.NewStyle().
			Foreground(colorDim).
			PaddingRight(1),

		// Sidebar: only a right border — the outer frame handles top/bottom/left
		sidebar: r.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(colorRule).
			PaddingTop(1),

		// Outer frame wrapping sidebar + body together
		outerBorder: r.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorRule),

		// Body: no border — outer frame handles it
		body: r.NewStyle().
			Padding(1, 2),

		header: r.NewStyle().
			Foreground(colorSienna).
			Bold(true),

		// Match card
		teamName: r.NewStyle().
			Foreground(colorGold).
			Bold(true).
			Width(6),

		score: r.NewStyle().
			Foreground(colorInk).
			Bold(true),

		liveDot: r.NewStyle().
			Foreground(colorLive).
			Bold(true),

		venue: r.NewStyle().
			Foreground(colorDim).
			Faint(true),

		result: r.NewStyle().
			Foreground(colorGreen),

		matchCard: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorRule).
			Padding(0, 2),

		// Points table
		tableHeader: r.NewStyle().
			Foreground(colorGold).
			Bold(true),

		tableRow: r.NewStyle().
			Foreground(colorInk),

		tableRowAlt: r.NewStyle().
			Foreground(colorInkMuted),

		highlight: r.NewStyle().
			Foreground(colorSienna).
			Bold(true),

		// General
		faint: r.NewStyle().
			Foreground(colorInkMuted).
			Faint(true),

		muted: r.NewStyle().
			Foreground(colorInkMuted),

		gold: r.NewStyle().
			Foreground(colorGold),

		loading: r.NewStyle().
			Foreground(colorSienna).
			Faint(true),
	}
}
