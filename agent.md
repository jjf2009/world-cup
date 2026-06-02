# ipl Repo Guide

This repository is a terminal-based IPL 2026 viewer built with Go, Bubble Tea, Lip Gloss, and Wish.
It runs as an SSH TUI server, fetches live cricket data from an API, and renders multiple views such as live scores, all matches, the points table, the schedule, historical winners, and an about screen.

## How The App Works

1. `main.go` starts an SSH server on `0.0.0.0:6767`.
2. Each SSH session gets a Bubble Tea model from `teaHandler`.
3. The model initializes by fetching match scores, schedule, points table, live scores, and winners in parallel.
4. The `Update` loop reacts to API responses, timers, key presses, and window resize events.
5. The `View` layer renders the active screen using Lip Gloss styles and table widgets.
6. Live data is refreshed every 5 seconds, and the live cursor blinks every second.

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
├── table.go
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
- fetched API data
- the match table widget and its styles

It also defines `Init`, which kicks off the first batch of API requests and animation timers.

### `update.go`

Contains the Bubble Tea state machine.
This file decides what happens when the app receives:

- API responses
- tick messages for loading cursor animation
- blinking updates for live scores
- refresh timers
- keyboard input
- terminal resize events

It also contains the navigation helpers for switching tabs and going back to the tab selector.

### `view.go`

Contains all rendering logic.
This file turns the current model state into terminal UI output.

Main sections:

- root layout with sidebar and body
- sidebar tab navigation
- live match view
- all matches table
- points table
- upcoming schedule
- historical winners
- about screen

It also renders detailed live-match information such as batters, bowler, partnership, wickets, and recent overs, plus squad columns when live match team details are available.

### `table.go`

Builds the all-matches table from API data.
It sorts matches into:

- live first
- upcoming next, sorted by time
- completed last, sorted by most recent

It also normalizes score and date display values and formats the table rows.

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
- squad data

It also contains team-name-to-slug normalization so squad data can be loaded from live match names.

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
This repo uses Go 1.26.2 and depends mainly on Bubble Tea, Bubbles, Lip Gloss, Wish, SSH, and Godotenv.

## Runtime Flow In Detail

### Startup

`main.go` loads environment variables, configures logging, and starts a Wish SSH server.
When a client connects, `teaHandler` creates a new `Model` with a Lip Gloss renderer.

### Initial Data Load

`Model.Init()` launches multiple fetch commands in parallel.
Each fetch command calls into the `cmd` package and returns a typed response message back into `Update()`.

### State Updates

`Update()` stores each response into `m.items` and clears the matching loading flag.
Once all core data has loaded, the app switches to the live view by default.

### Navigation

The app supports both single-key shortcuts and a tab selector view:

- `l` live
- `m` matches
- `p` points
- `s` schedule
- `h` historical
- `a` about
- `q` quit

In the tab selector, arrow keys move the selection and Enter or Right opens the selected tab.

### Live Refresh Loop

The app uses three tick loops:

- loading cursor blink while initial requests are pending
- live cursor blink every second
- live score refresh every 5 seconds

When live scores refresh, the app may also fetch squad data for the teams currently playing.

## Data And Environment

The HTTP layer reads `API_URL` from the environment.
The server also expects a writable `logs/` directory and uses `.ssh/term_info_ed25519` for its host key if one does not already exist.

## Notes For Future Changes

- Add new views by extending the view constants in `model.go`, the navigation maps in `view.go`, and the message handling in `update.go`.
- Add new API data by defining response types in `cmd/types.go` and fetch helpers in `cmd/http.go`.
- Keep table formatting changes in `table.go` and visual tweaks in `styles.go`.
- If you change rendering width behavior, check both the live view and the match table because they use different width rules.
