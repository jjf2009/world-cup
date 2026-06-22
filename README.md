# World Cup 2026 TUI

A terminal-based FIFA World Cup 2026 experience built with Go, Bubble Tea, Lip Gloss, and Wish.

The application runs as an SSH server, allowing anyone to follow the World Cup directly from their terminal without installing any software.

```bash
ssh worldcup2026.jaredfurtado.tech
```

## Features

### Live Matches

View ongoing World Cup matches with:

* Live scores
* Match status
* Match minute
* Stadium information
* Automatic updates

### Fixtures

Browse upcoming matches including:

* Teams
* Kickoff times
* Venues
* Match status

### Group Standings

Follow tournament standings with:

* Matches played
* Wins
* Draws
* Losses
* Goal difference
* Points

### Historical Winners

Explore previous FIFA World Cup tournaments including:

* Champions
* Runners-up
* Host nations
* Tournament years

### Real Data

The application uses real World Cup data fetched from public ESPN endpoints and cached locally for reliability.

### SSH Access

No installation required.

Simply connect:

```bash
ssh worldcup2026.jaredfurtado.tech
```

---

## Architecture

The TUI never communicates directly with external APIs.

A background fetcher continuously downloads World Cup data and stores it in local cache files.

```text
             ESPN API
                 │
                 ▼
         World Cup Fetcher
                 │
                 ▼
         Local JSON Cache
                 │
                 ▼
         World Cup Service
                 │
                 ▼
          Bubble Tea TUI
                 │
                 ▼
            SSH Server
```

This architecture ensures:

* Fast response times
* No API calls from the UI
* Graceful handling of API failures
* Cached data availability

---

## Project Structure

```text
cmd/
├── fetcher/

internal/
├── espn/
├── repository/
├── service/
├── simulation/

cache/
├── live_matches.json
├── fixtures.json
├── standings.json

data/
logs/
```

---

## Running Locally

### Start the Fetcher

```bash
go run ./cmd/fetcher
```

This downloads and updates:

```text
cache/live_matches.json
cache/fixtures.json
cache/standings.json
```

### Start the SSH Server

```bash
go run .
```

Connect from another terminal:

```bash
ssh localhost -p 6767
```

---

## Deployment

The production version runs on an Ubuntu VPS using systemd.

Services:

```bash
worldcup.service
worldcup-fetcher.service
```

The fetcher continuously refreshes World Cup data while the SSH server serves terminal sessions.

---

## Keyboard Shortcuts

| Key | Action             |
| --- | ------------------ |
| l   | Live Matches       |
| m   | Matches            |
| p   | Standings          |
| s   | Schedule           |
| h   | Historical Winners |
| a   | About              |
| q   | Quit               |

Navigation:

| Key       | Action |
| --------- | ------ |
| ↑ / ↓     | Move   |
| Enter / → | Select |
| ←         | Back   |

---

## Tech Stack

* Go
* Bubble Tea
* Lip Gloss
* Wish
* Bubbles
* Systemd
* Ubuntu VPS
* ESPN Public JSON Endpoints

---

## Inspiration

Inspired by the IPL TUI project by Harsh Iyer.

---

## Try It

```bash
ssh worldcup2026.jaredfurtado.tech
```

No installation required.
