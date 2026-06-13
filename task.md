````markdown
# World Cup 2026 TUI - Live Data Integration Tasks

## Goal

Replace the current simulated data with real World Cup 2026 data from ESPN's public JSON endpoints while preserving the existing TUI architecture.

The final user experience should remain:

```bash
ssh worldcup2026.jaredfurtado.tech
```

with live scores, real fixtures, and real standings.

---

# Current Architecture

Current flow:

```text
JSON Repositories
      ↓
WorldCupService
      ↓
Bubble Tea TUI
      ↓
Wish SSH Server
```

The architecture should remain unchanged.

The TUI must never call external APIs directly.

All external data must be fetched separately and written to local cache files.

---

# Data Sources

## Scoreboard Endpoint

Used for:

- Live Matches
- Fixtures
- Match Results

Endpoint:

https://site.api.espn.com/apis/site/v2/sports/soccer/fifa.world/scoreboard

Example query:

https://site.api.espn.com/apis/site/v2/sports/soccer/fifa.world/scoreboard?limit=200

Important fields:

```go
event.id

event.shortName

event.date

event.status.displayClock

event.status.type.state

event.status.type.description

event.status.type.completed

event.competitions[0].venue.fullName

event.competitions[0].competitors[0].team.displayName

event.competitions[0].competitors[1].team.displayName

event.competitions[0].competitors[0].score

event.competitions[0].competitors[1].score
```

State values:

```text
pre
in
post
```

---

## Standings Endpoint

Used for:

- Group Standings

Endpoint:

https://site.api.espn.com/apis/site/v2/sports/soccer/fifa.world/standings

---

# Required Architecture

Create:

```text
cmd/
├── fetcher/

cache/
├── live_matches.json
├── fixtures.json
├── standings.json
```

---

# Task 1

Create ESPN client package.

Example:

```text
internal/espn/
├── client.go
├── scoreboard.go
├── standings.go
├── types.go
```

Responsibilities:

- Fetch scoreboard
- Fetch standings
- Parse ESPN responses
- Return typed structs

---

# Task 2

Create Cache Models

Add cache models for:

```go
type LiveMatch struct {
    ID          string
    HomeTeam    string
    AwayTeam    string
    HomeScore   int
    AwayScore   int
    Minute      string
    Status      string
    Venue       string
}

type Fixture struct {
    ID          string
    HomeTeam    string
    AwayTeam    string
    KickoffTime time.Time
    Status      string
    Venue       string
}

type Standing struct {
    Team             string
    Played           int
    Won              int
    Drawn            int
    Lost             int
    GoalsFor         int
    GoalsAgainst     int
    GoalDifference   int
    Points           int
}
```

---

# Task 3

Build Fetcher Command

Create:

```text
cmd/fetcher/main.go
```

Responsibilities:

1. Fetch scoreboard
2. Fetch standings
3. Transform ESPN response
4. Write cache files

Output:

```text
cache/live_matches.json

cache/fixtures.json

cache/standings.json
```

---

# Task 4

Implement Atomic Cache Writes

Do NOT write directly.

Required flow:

```text
Fetch Data
      ↓
Validate
      ↓
Write temp file
      ↓
Rename temp file
      ↓
Replace cache
```

The TUI should never encounter partially written JSON.

---

# Task 5

Add Metadata

Each cache file should contain:

```json
{
  "updated_at": "2026-06-12T15:00:00Z",
  "data": []
}
```

This timestamp will be displayed in the UI later.

---

# Task 6

Create Cache Repositories

Create:

```text
internal/repository/
├── cache_live_repository.go
├── cache_fixture_repository.go
├── cache_standings_repository.go
```

Responsibilities:

- Read cache JSON
- Return domain objects
- No API calls

---

# Task 7

Replace Simulated Data Sources

Current TUI views should read:

```text
Live View
    ↓
live_matches.json

Fixtures View
    ↓
fixtures.json

Standings View
    ↓
standings.json
```

without changing the UI.

The goal is to swap data sources while preserving rendering logic.

---

# Task 8

Build Background Fetch Loop

Fetcher should run continuously.

Example:

```go
for {
    updateData()

    if liveMatchExists {
        sleep(60 seconds)
    } else {
        sleep(15 minutes)
    }
}
```

Detection logic:

```go
event.status.type.state == "in"
```

means a match is currently live.

---

# Task 9

Create Systemd Service

Create:

```text
/etc/systemd/system/worldcup-fetcher.service
```

Requirements:

- Auto start on boot
- Restart on failure
- Run independently of SSH sessions

Example:

```ini
[Unit]
Description=World Cup Fetcher
After=network.target

[Service]
WorkingDirectory=/opt/worldcup
ExecStart=/opt/worldcup/fetcher
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

---

# Task 10

Graceful Failure Handling

Requirements:

If ESPN:

- returns invalid JSON
- is unreachable
- times out

Then:

- Keep existing cache
- Log the error
- Continue running

Never:

- Delete cache
- Write empty cache
- Crash fetcher

The TUI must always have valid data available.

---

# Task 11

Logging

Create:

```text
logs/
├── fetcher.log
```

Log:

- fetch start
- fetch success
- fetch duration
- cache updates
- errors

---

# Task 12

Deployment

Final binaries:

```text
/opt/worldcup/

world-cup-tui
fetcher

data/
cache/
logs/
```

Services:

```bash
systemctl status worldcup
systemctl status worldcup-fetcher
```

Both services must survive reboot.

---

# Success Criteria

After deployment:

```bash
ssh worldcup2026.jaredfurtado.tech
```

should display:

- Real World Cup 2026 fixtures
- Real World Cup 2026 standings
- Live scores when matches are active

No simulator data should be visible.

The TUI must continue functioning even if ESPN becomes temporarily unavailable.
````
