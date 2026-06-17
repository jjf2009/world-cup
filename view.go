package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/h0i5/ipl/cmd"
	"github.com/h0i5/ipl/internal/domain"
)

var tabLabels = map[int]struct{ key, label string }{
	DashboardView:   {"d", "dashboard"},
	LiveView:        {"l", "live"},
	MatchView:       {"m", "matches"},
	PointsTableView: {"p", "standings"},
	ScheduleView:    {"s", "schedule"},
	HistoricalView:  {"h", "historical"},
	AboutView:       {"a", "about"},
}

// clamp returns v clamped to [lo, hi]
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ── Root ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	var cursor string
	if m.showLoadingCursor {
		cursor = m.styles.cursorStyle.Render("█")
	} else {
		cursor = " "
	}

	if m.currentView == InitialLoadView {
		content := lipgloss.JoinHorizontal(
			lipgloss.Center,
			"World Cup 2026",
			cursor,
		)

		return m.styles.loading.
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(content)
	}

	// Cap the total UI size so it doesn't sprawl on huge terminals
	totalW := clamp(m.width-4, 80, 160)
	totalH := int(float64(m.height) * 0.9)

	sidebarW := 24
	// Body gets the rest; -2 accounts for the sidebar's right border
	bodyW := totalW - sidebarW - 2

	sidebar := m.renderSidebar(sidebarW, totalH)
	body := m.renderBody(bodyW, totalH)

	inner := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, body)

	// Single outer border around the whole thing
	framed := m.styles.outerBorder.
		Width(totalW).
		Render(inner)

	// Center in the actual terminal window
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}

// ── Sidebar ───────────────────────────────────────────────────────────────────

func (m Model) renderSidebar(width, height int) string {
	s := m.styles

	innerW := width - 1 // -1 for the right border

	titleText := lipgloss.JoinHorizontal(lipgloss.Center,
		s.gold.Bold(true).Render("World Cup 2026"),
		s.muted.Render(" · jjf2009"),
	)
	title := lipgloss.NewStyle().Width(innerW).Align(lipgloss.Center).Render(titleText)

	divider := s.faint.Render(strings.Repeat("─", innerW))

	inTabView := m.currentView == TabView

	var rows []string
	rows = append(rows, title, divider, "")

	for _, view := range tabOrder {
		info := tabLabels[view]
		keyStr := fmt.Sprintf("[%s]", info.key)

		isActive := m.currentView == view              // currently open view
		isCursor := inTabView && view == m.selectedTab // hovered in tab nav

		var row string
		switch {
		case isActive:
			key := s.tabKey.Render(keyStr)
			label := s.tabActive.Render(info.label)
			row = lipgloss.JoinHorizontal(lipgloss.Center, " ", key, label)

		case isCursor:
			key := s.tabKey.Render(keyStr)
			label := s.highlight.Padding(0, 2).Render(info.label)
			row = lipgloss.JoinHorizontal(lipgloss.Center, " ", key, label)

		default:
			key := s.tabKeyDim.Render(keyStr)
			label := s.tabInactive.Render(info.label)
			row = lipgloss.JoinHorizontal(lipgloss.Center, " ", key, label)
		}

		rows = append(rows, row)
	}

	rows = append(rows, "", divider, "")
	rows = append(rows, s.faint.Render("  ↑↓  navigate"))
	if m.currentView != TabView {
		rows = append(rows, s.faint.Render("  ←   focus"))
	} else {
		rows = append(rows, s.faint.Render("  →   move"))
	}
	rows = append(rows, s.faint.Render("  q   quit"))

	content := strings.Join(rows, "\n")

	return s.sidebar.
		Width(width).
		Height(height).
		Render(content)
}

// ── Body ──────────────────────────────────────────────────────────────────────

func (m Model) renderBody(width, height int) string {
	s := m.styles

	innerW := width - 4 // subtract body padding (2 each side)

	activeView := m.currentView
	if m.currentView == TabView {
		activeView = m.selectedTab
	}

	var content string
	switch activeView {
	case DashboardView:
		content = m.renderDashboard(innerW)

	case LiveView:
		content = m.renderLive(innerW)

	case MatchView:
		content = m.renderMatches(innerW)

	case PointsTableView:
		content = m.renderStandings(innerW)

	case ScheduleView:
		content = m.renderSchedule(innerW)

	case HistoricalView:
		content = m.renderHistorical(innerW)

	case AboutView:
		content = m.renderAbout(innerW)
	}

	return s.body.
		Width(width).
		Height(height).
		Render(content)
}

// ── Live ──────────────────────────────────────────────────────────────────────

func (m Model) renderLive(width int) string {
	s := m.styles

	if m.loadingMap[LiveView] {
		return s.loading.Render("Fetching live scores...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder
	sb.WriteString(heading.Render("live matches"))
	sb.WriteString("\n\n")

	// Prefer the cache's own LastUpdated field; fall back to m.lastUpdated
	cacheTime := m.items.liveMatch.LastUpdated
	if cacheTime.IsZero() {
		cacheTime = m.lastUpdated
	}
	lastUpdated := "never"
	if !cacheTime.IsZero() {
		lastUpdated = cacheTime.UTC().Format("15:04 UTC")
	}

	match := m.items.liveMatch
	if match.ID == "" {
		sb.WriteString(s.muted.Render("No live matches currently in progress."))
		sb.WriteString("\n\n")

		next := nextUpcomingMatch(m.items.matches)
		if next != nil {
			now := time.Now().UTC()
			timeUntil := next.Kickoff.UTC().Sub(now)

			sb.WriteString(s.gold.Render("Next Match"))
			sb.WriteString("\n")
			sb.WriteString(s.teamName.Render(next.HomeTeam + "  vs  " + next.AwayTeam))
			sb.WriteString("\n")
			kickoffStr := next.Kickoff.UTC().Format("02 Jan · 15:04 UTC")
			sb.WriteString(s.muted.Render(kickoffStr))
			sb.WriteString("\n")
			if next.Venue != "" && next.Venue != "TBD" {
				sb.WriteString(s.venue.Render("📍 " + next.Venue))
				sb.WriteString("\n")
			}
			if timeUntil > 0 {
				sb.WriteString("\n")
				sb.WriteString(s.faint.Render("Next Kickoff: "))
				sb.WriteString(s.highlight.Render(formatCountdown(timeUntil)))
			}
		} else {
			sb.WriteString(s.faint.Render("No upcoming matches scheduled."))
		}
		return sb.String()
	}

	updatedLine := s.faint.Render(
		"last updated " + lastUpdated + " • auto updates every 1s",
	)
	sb.WriteString(updatedLine)
	sb.WriteString("\n\n")

	dot := "●"
	if !m.showLiveCursor {
		dot = " "
	}
	badge := s.liveDot.Render(dot + " " + match.Status)
	if match.Minute > 0 {
		badge = lipgloss.JoinHorizontal(lipgloss.Left,
			badge,
			s.faint.Render(fmt.Sprintf("  •  %d'", match.Minute)),
		)
	}
	sb.WriteString(badge + "\n\n")

	homeTeam := s.teamName.Render(truncate(match.HomeTeam, 24))
	awayTeam := s.teamName.Render(truncate(match.AwayTeam, 24))
	homeScore := s.score.Render(fmt.Sprintf("%d", match.HomeScore))
	awayScore := s.score.Render(fmt.Sprintf("%d", match.AwayScore))

	hArt := teamArt(match.HomeTeam)
	aArt := teamArt(match.AwayTeam)
	sLine := fmt.Sprintf("%d-%d", match.HomeScore, match.AwayScore)
	sArt := scoreArt(sLine)

	if hArt != "" && aArt != "" && sArt != "" && width >= 72 {
		// Count lines in each art block so we can vertically center the
		// score (which is much shorter) relative to the team banners.
		hLines := strings.Count(hArt, "\n")
		aLines := strings.Count(aArt, "\n")
		sLines := strings.Count(sArt, "\n")
		tallest := hLines
		if aLines > tallest {
			tallest = aLines
		}
		scorePad := (tallest - sLines) / 2
		if scorePad < 0 {
			scorePad = 0
		}
		scorePadStr := strings.Repeat("\n", scorePad)

		homeArtStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("166")).Render(hArt)
		scoreStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Padding(0, 3).Render(scorePadStr + sArt)
		awayArtStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("166")).Render(aArt)

		banner := lipgloss.JoinHorizontal(lipgloss.Top, homeArtStyled, scoreStyled, awayArtStyled)
		sb.WriteString(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(banner))
		sb.WriteString("\n")

		// Label row under the banner so it is always clear which team is which
		homeLabelW := lipgloss.Width(homeArtStyled)
		scoreW := lipgloss.Width(scoreStyled)
		awayLabelW := lipgloss.Width(awayArtStyled)
		homeLabel := lipgloss.NewStyle().Width(homeLabelW).Align(lipgloss.Center).
			Render(s.teamName.Render(truncate(match.HomeTeam, 24)))
		scoreSpacer := lipgloss.NewStyle().Width(scoreW).Render("")
		awayLabel := lipgloss.NewStyle().Width(awayLabelW).Align(lipgloss.Center).
			Render(s.teamName.Render(truncate(match.AwayTeam, 24)))
		labelRow := lipgloss.JoinHorizontal(lipgloss.Top, homeLabel, scoreSpacer, awayLabel)
		sb.WriteString(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(labelRow))
		sb.WriteString("\n\n")
	} else if width >= 80 {
		colW := width / 3
		midW := width - 2*colW

		scoreLine := fmt.Sprintf("%s  -  %s", homeScore, awayScore)

		leftCol := lipgloss.NewStyle().Width(colW).Align(lipgloss.Right).Render(homeTeam)
		midCol := lipgloss.NewStyle().Width(midW).Align(lipgloss.Center).Render(scoreLine)
		rightCol := lipgloss.NewStyle().Width(colW).Align(lipgloss.Left).Render(awayTeam)

		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftCol, midCol, rightCol))
		sb.WriteString("\n")
	} else {
		row1 := lipgloss.JoinHorizontal(lipgloss.Left, homeTeam, "  ", homeScore)
		row2 := lipgloss.JoinHorizontal(lipgloss.Left, awayTeam, "  ", awayScore)
		narrowW := width - 6
		if narrowW > 60 {
			narrowW = 60
		}
		sb.WriteString(s.matchCard.Width(narrowW).Render(strings.Join([]string{"", row1, row2}, "\n")))
	}

	if match.Venue != "" {
		sb.WriteString(s.venue.Render("Venue: "+match.Venue) + "\n")
	}
	if match.Group != "" {
		sb.WriteString(s.faint.Render("Group "+match.Group) + "\n")
	}
	if match.MatchNumber != "" {
		sb.WriteString(s.faint.Render("Match "+match.MatchNumber) + "\n")
	}
	if len(match.Scorers) > 0 {
		sb.WriteString("\n")
		sb.WriteString(s.tableHeader.Render("scorers"))
		sb.WriteString("\n")
		for _, scorer := range match.Scorers {
			minute := ""
			if scorer.Minute > 0 {
				minute = fmt.Sprintf(" %d'", scorer.Minute)
			}
			sb.WriteString(s.muted.Render(fmt.Sprintf("%s%s - %s", scorer.Name, minute, scorer.Team)))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ── Matches ───────────────────────────────────────────────────────────────────

func (m Model) renderMatches(width int) string {
	s := m.styles

	if m.loadingMap[MatchView] {
		return s.loading.Render("fetching matches...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder
	sb.WriteString(heading.Render("all matches"))
	sb.WriteString("\n\n")

	if len(m.items.matches) == 0 {
		sb.WriteString(s.muted.Render("No matches available"))
		sb.WriteString("\n\n")
		sb.WriteString(s.faint.Render("Press 'r' to refresh"))
		return sb.String()
	}

	tbl := m.matchTable
	if m.currentView == TabView {
		ts := m.matchTableStyles
		ts.Selected = lipgloss.NewStyle()
		tbl.SetStyles(ts)
	}

	sb.WriteString(tbl.View())
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render("  ↑↓ scroll  ← back"))

	return sb.String()
}

func (m Model) renderMatchCard(match domain.MatchView, width int, live bool) string {
	s := m.styles

	var badge string
	if live {
		badge = s.liveDot.Render("● LIVE")
	} else {
		badge = s.faint.Render("✓ " + match.Status)
	}

	homeTeam := s.teamName.Render(match.HomeTeam)
	awayTeam := s.teamName.Render(match.AwayTeam)
	homeScore := s.score.Render(fmt.Sprintf("%d", match.HomeScore))
	awayScore := s.score.Render(fmt.Sprintf("%d", match.AwayScore))

	teamsRow := fmt.Sprintf(
		"%s  %s    –    %s  %s",
		homeTeam, homeScore,
		awayScore, awayTeam,
	)

	matchNumber := ""
	if match.MatchNumber != "" {
		matchNumber = s.faint.Render("Match " + match.MatchNumber)
	}

	venueStr := s.venue.Render("📍 " + match.Venue)

	dateStr := ""
	if match.Date != "" {
		dateStr = s.faint.Render("🗓 " + match.Date)
	}

	var resultStr string
	if match.Result != "" {
		resultStr = s.result.Render("→ " + match.Result)
	}

	lines := []string{}
	if matchNumber != "" {
		lines = append(lines, matchNumber)
	}
	lines = append(lines, badge, teamsRow, venueStr)
	if dateStr != "" {
		lines = append(lines, dateStr)
	}
	if resultStr != "" {
		lines = append(lines, resultStr)
	}

	cardW := width - 6
	if cardW > 60 {
		cardW = 60
	}

	return s.matchCard.
		Width(cardW).
		Render(strings.Join(lines, "\n"))
}

// ── Standings ─────────────────────────────────────────────────────────────────

// buildStandingsContent builds the raw standings text used both by
// renderStandings (for the viewport) and by the resize handler in Update.
func (m Model) buildStandingsContent(width int) string {
	s := m.styles

	var sb strings.Builder

	// Header row — football columns: P W D L GF GA GD Pts
	header := fmt.Sprintf("%-3s %-28s %3s %3s %3s %3s %3s %3s %4s %4s",
		"#", "Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts",
	)
	sb.WriteString(s.tableHeader.Render(header))
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render(strings.Repeat("─", width)))
	sb.WriteString("\n")

	lastGroup := ""
	groupRow := 0
	for i, team := range m.items.standings {
		if team.Group != lastGroup {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(s.gold.Render("Group " + team.Group))
			sb.WriteString("\n")
			lastGroup = team.Group
			groupRow = 0
		}
		groupRow++

		row := fmt.Sprintf("%-3d %-28s %3d %3d %3d %3d %3d %3d %+4d %4d",
			team.Position,
			truncate(team.Team, 27),
			team.Played,
			team.Won,
			team.Drawn,
			team.Lost,
			team.GoalsFor,
			team.GoalsAgainst,
			team.GoalDifference,
			team.Points,
		)

		if groupRow <= 2 {
			sb.WriteString(s.tableRow.Render(row))
		} else {
			sb.WriteString(s.tableRowAlt.Render(row))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(s.faint.Render("top teams advance from group stage"))

	return sb.String()
}

func (m Model) renderStandings(width int) string {
	s := m.styles

	if m.loadingMap[PointsTableView] {
		return s.loading.Render("fetching standings...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder

	// Title row — always visible above the scrollable viewport
	sb.WriteString(heading.Render("group standings"))
	sb.WriteString("\n\n")

	if !m.standingsReady {
		// Viewport not yet sized (first render before WindowSizeMsg) — render inline
		sb.WriteString(m.buildStandingsContent(width))
		return sb.String()
	}

	// Render the viewport (scrollable standings table)
	sb.WriteString(m.standingsVP.View())
	sb.WriteString("\n")

	// Scroll-position hint at the bottom
	pct := m.standingsVP.ScrollPercent()
	var hint string
	if pct <= 0 {
		hint = "↓ scroll for more  •  ↑↓ / PgUp PgDn"
	} else if pct >= 1 {
		hint = "↑ top reached"
	} else {
		hint = fmt.Sprintf("↑↓ scroll  •  %.0f%%  •  ← back", pct*100)
	}
	sb.WriteString(s.faint.Render(hint))

	return sb.String()
}

// ── Schedule ──────────────────────────────────────────────────────────────────

func (m Model) renderSchedule(width int) string {
	s := m.styles

	if m.loadingMap[ScheduleView] {
		return s.loading.Render("fetching schedule...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder

	sb.WriteString(heading.Render("upcoming matches"))
	sb.WriteString("\n\n")

	header := fmt.Sprintf("%-3s  %-13s  %-13s  %-11s  %-5s  %s",
		"#", "Home", "Away", "Date", "Time", "Venue")

	sb.WriteString(renderScheduleSection(s, "Today", header, m.items.matchSchedule.Today, width))
	sb.WriteString("\n")
	sb.WriteString(renderScheduleSection(s, "Tomorrow", header, m.items.matchSchedule.Tomorrow, width))

	return sb.String()
}

func renderScheduleSection(s Styles, title, header string, matches []domain.MatchView, width int) string {
	var sb strings.Builder
	sb.WriteString(s.gold.Render(title))
	sb.WriteString("\n")
	sb.WriteString(s.tableHeader.Render(header))
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render(strings.Repeat("─", width)))
	sb.WriteString("\n")

	if len(matches) == 0 {
		sb.WriteString(s.muted.Render("No matches available"))
		sb.WriteString("\n")
		return sb.String()
	}

	for i, match := range matches {
		row := fmt.Sprintf("%-3s  %-13s  %-13s  %-11s  %-5s  %s",
			match.MatchNumber,
			truncate(match.HomeTeam, 13),
			truncate(match.AwayTeam, 13),
			truncate(match.Date, 11),
			match.Time,
			truncateVenue(match.Venue, 20),
		)
		if i%2 == 0 {
			sb.WriteString(s.tableRow.Render(row))
		} else {
			sb.WriteString(s.tableRowAlt.Render(row))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ── Historical ────────────────────────────────────────────────────────────────

func (m Model) renderHistorical(width int) string {
	s := m.styles

	if m.loadingMap[HistoricalView] {
		return s.loading.Render("fetching winners...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder
	sb.WriteString(heading.Render("world cup winners"))
	sb.WriteString("\n\n")

	header := fmt.Sprintf("%-6s  %-15s  %-15s  %s", "Year", "Winner", "Runner Up", "Venue")
	sb.WriteString(s.tableHeader.Render(header))
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render(strings.Repeat("─", width)))
	sb.WriteString("\n")

	for i, w := range m.items.historicalWinners {
		venue := w.Venue
		if venue == "" {
			venue = "-"
		}
		row := fmt.Sprintf("%-6s  %-15s  %-15s  %s",
			fmt.Sprintf("%d", w.Year),
			truncate(w.Winner, 15),
			truncate(w.RunnerUp, 15),
			truncate(venue, 30),
		)
		if i%2 == 0 {
			sb.WriteString(s.tableRow.Render(row))
		} else {
			sb.WriteString(s.tableRowAlt.Render(row))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	if m.currentView == HistoricalView {
		sb.WriteString(s.faint.Render("← back"))
	}

	return sb.String()
}

// ── Squads ────────────────────────────────────────────────────────────────────

func (m Model) renderSquads(slug1, slug2, name1, name2 string, width int) string {
	s := m.styles
	sq1, ok1 := m.items.squads[cmd.TeamToSlug(slug1)]
	sq2, ok2 := m.items.squads[cmd.TeamToSlug(slug2)]

	if !ok1 && !ok2 {
		return s.faint.Render("  loading squads...")
	}

	colW := (width - 4) / 2
	if colW < 20 {
		colW = 20
	}

	left := buildSquadColumn(s, name1, sq1, ok1, colW)
	right := buildSquadColumn(s, name2, sq2, ok2, colW)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(colW).Render(left),
		"    ",
		lipgloss.NewStyle().Width(colW).Render(right),
	)
}

func buildSquadColumn(s Styles, teamName string, sq cmd.SquadResponse, loaded bool, width int) string {
	divider := s.faint.Render(strings.Repeat("─", width-2))
	header := s.gold.Bold(true).Render(teamName)

	if !loaded {
		return strings.Join([]string{header, divider, s.faint.Render("  loading...")}, "\n")
	}

	type entry struct {
		key string
		p   cmd.SquadPlayer
	}
	seen := map[string]bool{}
	players := make([]entry, 0, len(sq.Squad))
	for k, p := range sq.Squad {
		if seen[p.Name] {
			continue
		}
		seen[p.Name] = true
		players = append(players, entry{k, p})
	}
	sort.Slice(players, func(i, j int) bool {
		ri, rj := squadPlayerRank(players[i].p), squadPlayerRank(players[j].p)
		if ri != rj {
			return ri < rj
		}
		return players[i].p.Name < players[j].p.Name
	})

	lines := []string{header, divider}
	for _, e := range players {
		p := e.p
		name := p.Name
		if len([]rune(name)) > 22 {
			name = string([]rune(name)[:21]) + "…"
		}
		nameW := lipgloss.NewStyle().Width(24).Render(s.muted.Render(name))
		tag := squadPlayerTag(p)
		numStr := ""
		if p.JerseyNumber > 0 {
			numStr = s.faint.Render(fmt.Sprintf("#%-3d", p.JerseyNumber))
		}
		lines = append(lines, "  "+numStr+" "+nameW+s.faint.Render(tag))
	}

	return strings.Join(lines, "\n")
}

// squadPlayerRank orders by position: GK → DEF → MID → FWD
func squadPlayerRank(p cmd.SquadPlayer) int {
	switch strings.ToUpper(p.Position) {
	case "GK", "GOALKEEPER":
		return 0
	case "DEF", "DEFENDER":
		return 1
	case "MID", "MIDFIELDER":
		return 2
	case "FWD", "FORWARD":
		return 3
	default:
		return 4
	}
}

func squadPlayerTag(p cmd.SquadPlayer) string {
	pos := strings.ToUpper(p.Position)
	var tag string
	switch {
	case strings.HasPrefix(pos, "GK") || pos == "GOALKEEPER":
		tag = "GK  "
	case strings.HasPrefix(pos, "DEF") || pos == "DEFENDER":
		tag = "DEF "
	case strings.HasPrefix(pos, "MID") || pos == "MIDFIELDER":
		tag = "MID "
	case strings.HasPrefix(pos, "FWD") || pos == "FORWARD":
		tag = "FWD "
	default:
		tag = strings.ToUpper(truncate(p.Position, 4))
	}
	if p.Nationality != "" {
		tag += " " + truncate(p.Nationality, 3)
	}
	return tag
}

// ── About ─────────────────────────────────────────────────────────────────────

func (m Model) renderAbout(width int) string {
	s := m.styles

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	lines := []string{
		heading.Bold(true).Render("world cup 2026 tui"),
		"",
		s.muted.Render("a terminal viewer for the FIFA World Cup 2026."),
		"",
		s.faint.Render("why?"),
		s.faint.Render("inspired by the IPL TUI built by harshiyer.in — go check it out"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			s.faint.Render("also check me out "),
			s.highlight.Render("jaredfurtado.tech"),
			s.faint.Render(" :)"),
		),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			s.faint.Render("built with bubbletea • code @ "),
			s.highlight.Render("https://github.com/h0i5/ipl"),
		),
		"",
	}

	return strings.Join(lines, "\n")
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

func truncateVenue(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	words := strings.Fields(s)
	if len(words) <= 1 {
		return truncate(s, max)
	}
	last := words[len(words)-1]
	prefix := words[:len(words)-1]
	for len(prefix) > 0 {
		candidate := strings.Join(prefix, " ") + "… " + last
		if len([]rune(candidate)) <= max {
			return candidate
		}
		prefix = prefix[:len(prefix)-1]
	}
	return "… " + last
}

// ── Dashboard ─────────────────────────────────────────────────────────────────

func (m Model) renderDashboard(width int) string {
	s := m.styles

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder
	sb.WriteString(heading.Bold(true).Render("World Cup 2026"))
	sb.WriteString("\n\n")

	// ── Tournament Stage ──────────────────────────────────────────────
	stage := deriveTournamentStage(m.items.standings)
	sb.WriteString(s.faint.Render("Current Stage"))
	sb.WriteString("\n")
	sb.WriteString(s.gold.Bold(true).Render(stage))
	sb.WriteString("\n\n")

	// ── Match Counts ──────────────────────────────────────────────────
	counts := todayMatchCounts(m.items.matches)
	divider := s.faint.Render(strings.Repeat("─", min(width, 40)))
	sb.WriteString(divider)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%-22s %s\n",
		s.muted.Render("Today's Matches"),
		s.score.Render(fmt.Sprintf("%d", counts.total)),
	))
	sb.WriteString(fmt.Sprintf("%-22s %s\n",
		s.liveDot.Render("● Live"),
		s.score.Render(fmt.Sprintf("%d", counts.live)),
	))
	sb.WriteString(fmt.Sprintf("%-22s %s\n",
		s.faint.Render("  Completed"),
		s.faint.Render(fmt.Sprintf("%d", counts.completed)),
	))
	sb.WriteString(fmt.Sprintf("%-22s %s\n",
		s.faint.Render("  Upcoming"),
		s.faint.Render(fmt.Sprintf("%d", counts.upcoming)),
	))
	sb.WriteString(divider)
	sb.WriteString("\n\n")

	// ── Live Match or Next Match ──────────────────────────────────────
	live := m.items.liveMatch
	if live.ID != "" {
		dot := "●"
		if !m.showLiveCursor {
			dot = " "
		}
		sb.WriteString(s.liveDot.Render(dot + " LIVE NOW"))
		sb.WriteString("\n")
		sb.WriteString(s.teamName.Bold(true).Render(live.HomeTeam + "  vs  " + live.AwayTeam))
		sb.WriteString("\n")
		scoreStr := fmt.Sprintf("%d – %d", live.HomeScore, live.AwayScore)
		if live.Minute > 0 {
			scoreStr += fmt.Sprintf("  (%d')", live.Minute)
		}
		sb.WriteString(s.score.Render(scoreStr))
		sb.WriteString("\n")
		if live.Venue != "" {
			sb.WriteString(s.venue.Render("📍 " + live.Venue))
			sb.WriteString("\n")
		}
	} else {
		next := nextUpcomingMatch(m.items.matches)
		if next != nil {
			now := time.Now().UTC()
			timeUntil := next.Kickoff.UTC().Sub(now)
			sb.WriteString(s.gold.Render("Next Match"))
			sb.WriteString("\n")
			sb.WriteString(s.teamName.Render(next.HomeTeam + "  vs  " + next.AwayTeam))
			sb.WriteString("\n")
			sb.WriteString(s.muted.Render(next.Kickoff.UTC().Format("02 Jan · 15:04 UTC")))
			sb.WriteString("\n")
			if next.Venue != "" && next.Venue != "TBD" {
				sb.WriteString(s.venue.Render("📍 " + next.Venue))
				sb.WriteString("\n")
			}
			if timeUntil > 0 {
				sb.WriteString(s.faint.Render("Next Kickoff: "))
				sb.WriteString(s.highlight.Render(formatCountdown(timeUntil)))
				sb.WriteString("\n")
			}
		} else {
			sb.WriteString(s.muted.Render("No live matches currently in progress."))
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")

	// ── Cache Timestamp ───────────────────────────────────────────────
	cacheTime := m.items.liveMatch.LastUpdated
	if cacheTime.IsZero() {
		cacheTime = m.lastUpdated
	}
	if !cacheTime.IsZero() {
		sb.WriteString(s.faint.Render("Last Updated: " + cacheTime.UTC().Format("15:04 UTC")))
		sb.WriteString("\n\n")
	}

	// ── Top Teams ─────────────────────────────────────────────────────
	topTeams := topStandingTeams(m.items.standings, 3)
	if len(topTeams) > 0 {
		sb.WriteString(s.tableHeader.Render("Top Teams"))
		sb.WriteString("\n")
		medals := []string{"🥇", "🥈", "🥉"}
		for i, t := range topTeams {
			medal := medals[i]
			row := fmt.Sprintf("%s  %-22s %s",
				medal,
				truncate(t.Team, 20),
				s.faint.Render(fmt.Sprintf("%dpts", t.Points)),
			)
			sb.WriteString(s.tableRow.Render(row))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// ── Navigation hint ───────────────────────────────────────────────
	sb.WriteString(s.faint.Render("[l] Live  [m] Matches  [p] Standings"))

	return sb.String()
}

// ── Shared helpers ────────────────────────────────────────────────────────────

type matchCounts struct {
	total     int
	live      int
	completed int
	upcoming  int
}

// todayMatchCounts counts matches scheduled for today (UTC).
func todayMatchCounts(matches []domain.MatchView) matchCounts {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	tomorrow := today.AddDate(0, 0, 1)

	var c matchCounts
	for _, m := range matches {
		k := m.Kickoff.UTC()
		if k.Before(today) || !k.Before(tomorrow) {
			continue
		}
		c.total++
		switch m.Status {
		case domain.StatusLive:
			c.live++
		case domain.StatusFinished:
			c.completed++
		default:
			c.upcoming++
		}
	}
	return c
}

// nextUpcomingMatch returns the earliest future match that hasn't started yet.
func nextUpcomingMatch(matches []domain.MatchView) *domain.MatchView {
	now := time.Now().UTC()
	for i := range matches {
		m := matches[i]
		if m.Kickoff.UTC().After(now) && m.Status != domain.StatusFinished {
			return &matches[i]
		}
	}
	return nil
}

// formatCountdown formats a duration as "2h 14m" or "45m".
func formatCountdown(d time.Duration) string {
	if d <= 0 {
		return "Starting soon"
	}
	h := int(d.Hours())
	mins := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm", h, mins)
	}
	return fmt.Sprintf("%dm", int(d.Minutes()))
}

// deriveTournamentStage inspects standings group names to determine the current stage.
func deriveTournamentStage(standings []domain.StandingView) string {
	for _, sv := range standings {
		g := strings.ToLower(sv.Group)
		switch {
		case strings.Contains(g, "round of 32"):
			return "Round of 32"
		case strings.Contains(g, "round of 16"):
			return "Round of 16"
		case strings.Contains(g, "quarter"):
			return "Quarter Finals"
		case strings.Contains(g, "semi"):
			return "Semi Finals"
		case strings.Contains(g, "final"):
			return "Final"
		}
	}
	return "Group Stage"
}

// topStandingTeams returns the top n teams ranked by points across all groups.
func topStandingTeams(standings []domain.StandingView, n int) []domain.StandingView {
	sorted := append([]domain.StandingView(nil), standings...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Points != sorted[j].Points {
			return sorted[i].Points > sorted[j].Points
		}
		if sorted[i].GoalDifference != sorted[j].GoalDifference {
			return sorted[i].GoalDifference > sorted[j].GoalDifference
		}
		return sorted[i].GoalsFor > sorted[j].GoalsFor
	})
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	return sorted
}

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

