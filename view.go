package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214")).
		Render("⚽ FIFA WORLD CUP 2026")

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		MarginBottom(1)

	liveMatches := card.Render(
		header.Render("LIVE MATCHES") +
			"\n\n🇧🇷 Brazil 2 - 1 Argentina 🇦🇷" +
			"\n72'" +
			"\n\n🇫🇷 France 0 - 0 Germany 🇩🇪" +
			"\n34'",
	)

	groupTable := card.Render(
		header.Render("GROUP A") +
			"\n\n1. Netherlands   6 pts" +
			"\n2. Senegal       4 pts" +
			"\n3. Ecuador       1 pt" +
			"\n4. Qatar         0 pts",
	)

	upcoming := card.Render(
		header.Render("UPCOMING") +
			"\n\nEngland vs USA" +
			"\n20:00 UTC" +
			"\n\nSpain vs Japan" +
			"\n23:00 UTC",
	)

	content := strings.Join([]string{
		title,
		"",
		liveMatches,
		groupTable,
		upcoming,
		"",
		fmt.Sprintf("Terminal Size: %dx%d", m.width, m.height),
		"",
		"[q] quit",
	}, "\n")

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}