# World Cup TUI Repo Guide

This repository is a terminal-based FIFA World Cup 2026 TUI viewer built with Go, Bubble Tea, Lip Gloss, and Wish.
It runs as an SSH TUI server, providing a terminal interface to view matches, standings, schedule, historical winners, and live match details.

## How The App Works

1. `main.go` starts an SSH server on `0.0.0.0:6767`.
2. Each SSH session gets a Bubble Tea model from `teaHandler`.
3. The model initializes with populated mock data for live scores, matches, standings, schedule, and historical winners.
4. The `Update` loop reacts to timers, key presses, and window resize events.
5. The `View` layer renders the active screen using Lip Gloss styles and table widgets.
6. The live cursor blinks periodically to indicate active operation.

## Folder Structure

```text
.
├── ascii.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── main.go
├── model.go
├── README.md
├── styles.go
├── update.go
├── view.go
└── cmd/
    ├── http.go
    └── types.go
```

## File By File

### `main.go`

Starts the SSH server and wires the Bubble Tea middleware.
It also loads environment variables from `.env`, writes logs to `logs/server.log`, and shuts down cleanly on SIGINT or SIGTERM.

### `model.go`

Defines the main Bubble Tea `Model` and the data it stores.
This is where the app’s state lives:

- the active view
- loading flags for each section
- window size
- last update time
- fetched/mocked API data
- the match table widget and its styles

It also defines `Init`, which kicks off the animation/refresh timers.

### `update.go`

Contains the Bubble Tea state machine.
This file decides what happens when the app receives:

- tick messages for loading cursor animation
- blinking updates for live scores
- keyboard input
- terminal resize events

It also contains the navigation helpers for switching tabs and going back to the tab selector.

### `view.go`

Contains all rendering logic.
This file turns the current model state into terminal UI output.

Main sections:

- root layout with sidebar and body
- sidebar tab navigation
- live match view (rendering match details, status, minute, score, and venue)
- all matches table
- group standings (points table)
- upcoming schedule
- historical winners
- about screen

It also contains `renderSquads` to display squad columns (GK, DEF, MID, FWD) when live match team details are available.

### `table.go` (Not present or integrated into main logic)

*(Note: Match table is built using Bubble Tea's built-in table component inside `model.go` and `view.go`)*

### `styles.go`

Defines the visual theme for the TUI.
This includes the color palette, sidebar styling, borders, table styles, match card styles, and general text styles.

### `ascii.go`

Stores ASCII art for team names and score digits.
The live view uses these helpers to render larger visual team and score blocks when there is enough width.

### `cmd/http.go`

Handles HTTP fetching from the backend API.
It exposes typed helper functions for:

- match scores
- match schedule
- points table
- live match scores
- historical winners

It also contains a team-to-slug helper for matching team names to their slugs.

### `cmd/types.go`

Defines the JSON response shapes used by the API layer.
These types match the data returned by the backend endpoints and are consumed by the model and renderers.

### `README.md`

Contains a short project overview and screenshots of the rendered screens.

### `Dockerfile`

Container build definition for running the app in a containerized environment.

### `docker-compose.yml`

Local compose setup for running the service with its supporting environment.

### `go.mod`

Declares the module path and dependencies.
This repo uses Go 1.25.2 and depends mainly on Bubble Tea, Bubbles, Lip Gloss, Wish, SSH, and Godotenv.

## Runtime Flow In Detail

### Startup

`main.go` loads environment variables, configures logging, and starts a Wish SSH server.
When a client connects, `teaHandler` creates a new `Model` with a Lip Gloss renderer.

### Initial Data Load

`Model.Init()` launches timers. The `Model` is initialized with mock data (e.g. Argentina vs France) in `NewModel` for a static-view demonstration. The HTTP client functions in `cmd/http.go` are prepared to query the API endpoints when dynamic loading is enabled.

### State Updates

`Update()` stores responses or handles user events. Keyboard controls change the views and allow scrolling.

### Navigation

The app supports both single-key shortcuts and a tab selector view:

- `l` live
- `m` matches
- `p` standings (points table)
- `s` schedule
- `h` historical
- `a` about
- `q` quit

In the tab selector, arrow keys move the selection and Enter or Right opens the selected tab.

### Live Refresh Loop

The app uses two tick loops:

- loading cursor blink while initial requests are pending
- live cursor blink every second

## Data And Environment

The HTTP layer reads `API_URL` from the environment if configured.
The server also expects a writable `logs/` directory and uses `.ssh/term_info_ed25519` for its host key if one does not already exist.

## Notes For Future Changes

- Add new views by extending the view constants in `model.go` and the navigation maps in `view.go` / `update.go`.
- Add new API data by defining response types in `cmd/types.go` and fetch helpers in `cmd/http.go`.
- Keep visual tweaks and colors in `styles.go`.
