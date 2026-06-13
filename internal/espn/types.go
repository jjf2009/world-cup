package espn

// ESPN API response types.
// Date strings are kept as raw strings because ESPN uses a truncated
// ISO-8601 format ("2026-06-12T19:00Z") that Go's time.Time cannot
// unmarshal natively.

// ScoreboardResponse represents the root structure of the ESPN scoreboard JSON.
type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

// Event represents an individual match event.
type Event struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Date         string        `json:"date"`
	Status       EventStatus   `json:"status"`
	Competitions []Competition `json:"competitions"`
}

// EventStatus represents the status of the match (e.g. clock, finished).
type EventStatus struct {
	Clock        float64    `json:"clock"`
	DisplayClock string     `json:"displayClock"`
	Type         StatusType `json:"type"`
}

// StatusType represents detailed status info.
type StatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"` // "pre", "in", "post"
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
}

// Competition represents match details.
type Competition struct {
	ID          string       `json:"id"`
	Date        string       `json:"date"`
	Competitors []Competitor `json:"competitors"`
	Venue       *Venue       `json:"venue,omitempty"`
}

// Competitor represents a participating team in the match.
type Competitor struct {
	ID       string         `json:"id"`
	HomeAway string         `json:"homeAway"` // "home" or "away"
	Score    string         `json:"score"`
	Winner   bool           `json:"winner"`
	Team     CompetitorTeam `json:"team"`
}

// CompetitorTeam contains details of a competitor team.
type CompetitorTeam struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	DisplayName  string `json:"displayName"`
	Name         string `json:"name"`
	Location     string `json:"location"`
}

// Venue represents the match venue details.
type Venue struct {
	FullName string `json:"fullName"`
}

// StandingsResponse represents the root structure of the ESPN standings JSON.
type StandingsResponse struct {
	Children []StandingChild `json:"children"`
}

// StandingChild represents a tournament stage/group containing standings.
type StandingChild struct {
	Name      string           `json:"name"` // E.g., "Group A"
	Standings StandingWrapper  `json:"standings"`
}

// StandingWrapper wraps standings entries.
type StandingWrapper struct {
	Entries []StandingEntry `json:"entries"`
}

// StandingEntry represents a single team's standing row in a group.
type StandingEntry struct {
	Team  StandingTeam   `json:"team"`
	Stats []StandingStat `json:"stats"`
}

// StandingTeam represents the team inside standings.
type StandingTeam struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	DisplayName  string `json:"displayName"`
}

// StandingStat represents a statistical value inside standings.
type StandingStat struct {
	Name             string  `json:"name"`
	DisplayName      string  `json:"displayName"`
	ShortDisplayName string  `json:"shortDisplayName"`
	Description      string  `json:"description"`
	Abbreviation     string  `json:"abbreviation"`
	Type             string  `json:"type"`
	Value            float64 `json:"value"`
	DisplayValue     string  `json:"displayValue"`
}
