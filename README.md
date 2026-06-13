# World Cup 2026 TUI

A terminal-based FIFA World Cup 2026 viewer built using **Go**, **Bubble Tea**, **Lip Gloss**, and **Wish**. The application runs as a secure SSH server, allowing users to connect and check match details, standings, and historical data directly from their terminals.

Inspired by the IPL TUI by [harshiyer.in](https://ipl.harshiyer.in).

## Features

- **Live Matches**: Keep track of ongoing games, including scorelines, match minutes, and venue.
- **Match Results (All Matches)**: Interactive table of all matches categorized by status (Live, Upcoming, Completed).
- **Group Standings**: View real-time group stages standings with stats like Played, Won, Drawn, Lost, Goal Difference, and Points.
- **Fixtures & Schedule**: Look ahead at upcoming match dates, kickoff times, and host venues.
- **Historical Winners**: Browse past tournament winners, runners-up, and host nations dating back through World Cup history.
- **Squads**: Deep dive into team squad details structured by positions (Goalkeepers, Defenders, Midfielders, Forwards) and jersey numbers.

## How It Works

The app runs as a standalone SSH server. When a client connects via SSH, the server spawns an interactive Bubble Tea terminal application inside the user's terminal window.

```
                  ┌──────────────────────┐
                  │   World Cup 2026     │
                  │   TUI SSH Server     │
                  └──────────▲───────────┘
                             │ SSH Port 6767
                             │
  ┌──────────────────────────┴──────────────────────────┐
  │                                                     │
┌─┴───────────────────┐                               ┌─┴───────────────────┐
│     User Term       │                               │     User Term       │
│  ssh localhost -p   │                               │  ssh localhost -p   │
└─────────────────────┘                               └─────────────────────┘
```

## Getting Started

### Prerequisites

- Go 1.25.2 or later
- An SSH client (built into Linux/macOS and Windows PowerShell/Command Prompt)

### Running Locally

1. **Start the SSH server**:
   ```bash
   go run .
   ```
   This will start the server on `0.0.0.0:6767` by default.

2. **Connect using SSH**:
   Open a new terminal window and run:
   ```bash
   ssh localhost -p 6767
   ```

### Running with Docker

You can also run the TUI inside a container:

```bash
docker-compose up --build -d
```

## Keyboard Shortcuts

Navigate through the screens using the keyboard:

- `l` - Go to **Live View**
- `m` - Go to **Matches**
- `p` - Go to **Standings (Points Table)**
- `s` - Go to **Schedule**
- `h` - Go to **Historical Winners**
- `a` - Go to **About**
- `q` / `Ctrl+C` - Quit the application

In the tab sidebar navigation view:
- `↑` / `↓` - Move cursor selection
- `→` / `Enter` - Select and enter the view
- `←` - Return focus to the sidebar navigation

## Configuration

You can configure the base URL for the live matches API via the `.env` file:

```env
API_URL=http://localhost:8080
```

## Tech Stack

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Style definitions and layouts
- **[Wish](https://github.com/charmbracelet/wish)** - SSH server library
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - Interactive components (tables)
