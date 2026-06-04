package main

import (
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

		case LiveView, PointsTableView, ScheduleView, HistoricalView, AboutView:
			if key == "left" {
				return m.handleNavToTabView()
			}

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

		// Auto-enter live view once terminal size is known
		if m.currentView == InitialLoadView {
			m.currentView = LiveView
			m.selectedTab = LiveView
		}
	}

	return m, tea.Batch(cmds...)
}