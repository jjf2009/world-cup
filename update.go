package main

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleTabNavigation(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "l":
		m.selectedTab = LiveView
	case "m":
		m.selectedTab = MatchView
	case "p":
		m.selectedTab = PointsTableView
	case "s":
		m.selectedTab = ScheduleView
	case "h":
		m.selectedTab = HistoricalView
	case "a":
		m.selectedTab = AboutView
	}

	m.currentView = m.selectedTab
	return m, nil
}

func (m Model) handleNavToTabView() (tea.Model, tea.Cmd) {
	m.selectedTab = m.currentView
	m.currentView = TabView
	return m, nil
}

func (m Model) tabCursorIndex() int {
	for i, v := range tabOrder {
		if v == m.selectedTab {
			return i
		}
	}
	return 0
}

func (m Model) handleTabCursor(key string) (tea.Model, tea.Cmd) {
	idx := m.tabCursorIndex()

	switch key {
	case "up":
		idx--

	case "down":
		idx++

	case "enter", "right":
		m.currentView = m.selectedTab
		return m, nil
	}

	if idx < 0 {
		idx = len(tabOrder) - 1
	}

	if idx >= len(tabOrder) {
		idx = 0
	}

	m.selectedTab = tabOrder[idx]

	return m, nil
}

func (m Model) handleQuit(key string) (tea.Model, tea.Cmd) {
	if key == "q" || key == "ctrl+c" {
		return m, tea.Quit
	}

	return m, nil
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(
		500*time.Millisecond,
		func(t time.Time) tea.Msg {
			return tickMsg(t)
		},
	)
}

type liveTickMsg time.Time

func liveTickCmd() tea.Cmd {
	return tea.Tick(
		1000*time.Millisecond,
		func(t time.Time) tea.Msg {
			return liveTickMsg(t)
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tickMsg:
		m.showLoadingCursor = !m.showLoadingCursor
		return m, tickCmd()

	case liveTickMsg:
		m.showLiveCursor = !m.showLiveCursor
		m.refreshData(context.Background())
		// Refresh standings viewport content after data update
		if m.standingsReady {
			oldYOffset := m.standingsVP.YOffset
			m.standingsVP.SetContent(m.buildStandingsContent(m.standingsVP.Width))
			m.standingsVP.SetYOffset(oldYOffset)
		}
		return m, liveTickCmd()

	case tea.KeyMsg:
		key := msg.String()

		if model, cmd := m.handleQuit(key); cmd != nil {
			return model, cmd
		}

		switch key {
		case "l", "m", "p", "s", "h", "a":
			return m.handleTabNavigation(key)
		}

		switch m.currentView {

		case LiveView, ScheduleView, HistoricalView, AboutView:
			if key == "left" {
				return m.handleNavToTabView()
			}

		case PointsTableView:
			if key == "left" {
				return m.handleNavToTabView()
			}
			// Forward scroll keys to the standings viewport
			var cmd tea.Cmd
			m.standingsVP, cmd = m.standingsVP.Update(msg)
			return m, cmd

		case MatchView:
			if key == "left" {
				return m.handleNavToTabView()
			}
			// Forward keyboard events to the table to allow scrolling
			var cmd tea.Cmd
			m.matchTable, cmd = m.matchTable.Update(msg)
			return m, cmd

		case TabView:
			return m.handleTabCursor(key)

		case InitialLoadView:
			// ignore input
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Compute the available body area for the standings viewport.
		// Mirrors the sizing logic in View(): totalW capped 80-160, bodyW = totalW - sidebarW - 2
		// Body style has Padding(1,2) = 4 chars width, 2 rows height.
		totalW := clamp(m.width-4, 80, 160)
		sidebarW := 24
		bodyW := totalW - sidebarW - 2
		vpW := bodyW - 4 // body padding 2 each side

		// totalH = 90% of terminal height; subtract outer border (2) + body padding (2) + header rows (~4)
		totalH := int(float64(m.height) * 0.9)
		vpH := totalH - 8
		if vpH < 5 {
			vpH = 5
		}

		if !m.standingsReady {
			m.standingsVP.Width = vpW
			m.standingsVP.Height = vpH
			m.standingsVP.SetContent(m.buildStandingsContent(vpW))
			m.standingsReady = true
		} else {
			m.standingsVP.Width = vpW
			m.standingsVP.Height = vpH
			// re-render content in case column widths changed
			oldYOffset := m.standingsVP.YOffset
			m.standingsVP.SetContent(m.buildStandingsContent(vpW))
			m.standingsVP.SetYOffset(oldYOffset)
		}

		// Auto-enter live view once terminal size is known
		if m.currentView == InitialLoadView {
			m.currentView = LiveView
			m.selectedTab = LiveView
		}
	}

	return m, tea.Batch(cmds...)
}
