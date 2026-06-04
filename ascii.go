package main

// this file stores the ascii arts for World Cup 2026 team names and score digits
// generated from https://patorjk.com/software/taag (Font: Big Money-nw)

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Teams ─────────────────────────────────────────────────────────────────────

const teamBRA string = `
 /$$$$$$$  /$$$$$$$   /$$$$$$  
| $$__  $$| $$__  $$ /$$__  $$ 
| $$  \ $$| $$  \ $$| $$  \ $$ 
| $$$$$$$ | $$$$$$$/| $$$$$$$$ 
| $$__  $$| $$__  $$| $$__  $$ 
| $$  \ $$| $$  \ $$| $$  \ $$ 
| $$$$$$$/| $$  | $$| $$  | $$ 
|_______/ |__/  |__/|__/  |__/ 
`

const teamARG string = `
  /$$$$$$  /$$$$$$$   /$$$$$$ 
 /$$__  $$| $$__  $$ /$$__  $$
| $$  \ $$| $$  \ $$| $$  \__/
| $$$$$$$$| $$$$$$$/| $$ /$$$$
| $$__  $$| $$__  $$| $$|_  $$
| $$  | $$| $$  \ $$| $$  \ $$
| $$  | $$| $$  | $$|  $$$$$$/
|__/  |__/|__/  |__/ \______/ 
`

const teamFRA string = `
 /$$$$$$$$ /$$$$$$$   /$$$$$$  
| $$_____/| $$__  $$ /$$__  $$ 
| $$      | $$  \ $$| $$  \ $$ 
| $$$$$   | $$$$$$$/| $$$$$$$$ 
| $$__/   | $$__  $$| $$__  $$ 
| $$      | $$  \ $$| $$  \ $$ 
| $$      | $$  | $$| $$  | $$ 
|__/      |__/  |__/|__/  |__/ 
`

const teamENG string = `
 /$$$$$$$$ /$$   /$$  /$$$$$$ 
| $$_____/| $$$ | $$ /$$__  $$
| $$      | $$$$| $$| $$  \__/
| $$$$$   | $$ $$ $$| $$ /$$$$
| $$__/   | $$  $$$$| $$|_  $$
| $$      | $$\  $$$| $$  \ $$
| $$$$$$$$| $$ \  $$|  $$$$$$/
|________/|__/  \__/ \______/ 
`

const teamESP string = `
 /$$$$$$$$ /$$$$$$  /$$$$$$$  
| $$_____//$$__  $$| $$__  $$ 
| $$     | $$  \__/| $$  \ $$ 
| $$$$$  |  $$$$$$ | $$$$$$$/
| $$__/   \____  $$| $$____/ 
| $$      /$$  \ $$| $$      
| $$$$$$$$|  $$$$$$/| $$      
|________/ \______/ |__/      
`

const teamGER string = `
  /$$$$$$  /$$$$$$$$ /$$$$$$$  
 /$$__  $$| $$_____/| $$__  $$ 
| $$  \__/| $$      | $$  \ $$ 
| $$ /$$$$| $$$$$   | $$$$$$$/
| $$|_  $$| $$__/   | $$__  $$ 
| $$  \ $$| $$      | $$  \ $$ 
|  $$$$$$/| $$$$$$$$| $$  | $$ 
 \______/ |________/|__/  |__/ 
`

const teamPOR string = `
 /$$$$$$$   /$$$$$$  /$$$$$$$  
| $$__  $$ /$$__  $$| $$__  $$ 
| $$  \ $$| $$  \ $$| $$  \ $$ 
| $$$$$$$/| $$  | $$| $$$$$$$/
| $$____/ | $$  | $$| $$__  $$ 
| $$      | $$  | $$| $$  \ $$ 
| $$      |  $$$$$$/| $$  | $$ 
|__/       \______/ |__/  |__/ 
`

const teamNED string = `
 /$$   /$$ /$$$$$$$$ /$$$$$$$ 
| $$$ | $$| $$_____/| $$__  $$
| $$$$| $$| $$      | $$  \ $$
| $$ $$ $$| $$$$$   | $$  | $$
| $$  $$$$| $$__/   | $$  | $$
| $$\  $$$| $$      | $$  | $$
| $$ \  $$| $$$$$$$$| $$$$$$$/
|__/  \__/|________/|_______/ 
`

const teamUSA string = `
 /$$   /$$  /$$$$$$   /$$$$$$  
| $$  | $$ /$$__  $$ /$$__  $$ 
| $$  | $$| $$  \__/| $$  \ $$ 
| $$  | $$|  $$$$$$ | $$$$$$$$ 
| $$  | $$ \____  $$| $$__  $$ 
| $$  | $$ /$$  \ $$| $$  | $$ 
|  $$$$$$/|  $$$$$$/| $$  | $$ 
 \______/  \______/ |__/  |__/ 
`

const teamMEX string = `
 /$$      /$$ /$$$$$$$$ /$$   /$$
| $$$    /$$$ | $$_____/| $$  / $$
| $$$$  /$$$$| $$      |  $$/ $$/
| $$ $$/$$ $$| $$$$$    \  $$$$/
| $$  $$$| $$| $$__/     >$$  $$
| $$\  $ | $$| $$       /$$/\  $$
| $$ \/  | $$| $$$$$$$$| $$  \ $$
|__/     |__/|________/|__/  |__/
`

const teamCAN string = `
  /$$$$$$   /$$$$$$  /$$   /$$
 /$$__  $$ /$$__  $$| $$$ | $$
| $$  \__/| $$  \ $$| $$$$| $$
| $$      | $$$$$$$$| $$ $$ $$
| $$      | $$__  $$| $$  $$$$
| $$    $$| $$  | $$| $$\  $$$
|  $$$$$$/| $$  | $$| $$ \  $$
 \______/ |__/  |__/|__/  \__/
`

const teamMOR string = `
 /$$      /$$  /$$$$$$  /$$$$$$$  
| $$$    /$$$ /$$__  $$| $$__  $$ 
| $$$$  /$$$$| $$  \ $$| $$  \ $$ 
| $$ $$/$$ $$| $$  | $$| $$$$$$$/
| $$  $$$| $$| $$  | $$| $$__  $$ 
| $$\  $ | $$| $$  | $$| $$  \ $$ 
| $$ \/  | $$|  $$$$$$/| $$  | $$ 
|__/     |__/ \______/ |__/  |__/ 
`

const teamJPN string = `
      /$$ /$$$$$$$  /$$   /$$
     | $$| $$__  $$| $$$ | $$
     | $$| $$  \ $$| $$$$| $$
     | $$| $$$$$$$/| $$ $$ $$
/$$  | $$| $$____/ | $$  $$$$
| $$  | $$| $$      | $$\  $$$
|  $$$$$$/| $$      | $$ \  $$
 \______/ |__/      |__/  \__/
`

const teamKOR string = `
 /$$   /$$ /$$$$$$  /$$$$$$$  
| $$  /$$//$$__  $$| $$__  $$ 
| $$ /$$/ | $$  \ $$| $$  \ $$ 
| $$$$$/  | $$  | $$| $$$$$$$/
| $$  $$  | $$  | $$| $$__  $$ 
| $$\  $$ | $$  | $$| $$  \ $$ 
| $$ \  $$|  $$$$$$/| $$  | $$ 
|__/  \__/ \______/ |__/  |__/ 
`

const teamAUS string = `
  /$$$$$$  /$$   /$$  /$$$$$$ 
 /$$__  $$| $$  | $$ /$$__  $$
| $$  \ $$| $$  | $$| $$  \__/
| $$$$$$$$| $$  | $$|  $$$$$$ 
| $$__  $$| $$  | $$ \____  $$
| $$  | $$| $$  | $$ /$$  \ $$
| $$  | $$|  $$$$$$/|  $$$$$$/
|__/  |__/ \______/  \______/ 
`

const teamITA string = `
 /$$$$$$ /$$$$$$$$ /$$$$$$ 
|_  $$_/|__  $$__//$$__  $$
  | $$     | $$  | $$  \ $$
  | $$     | $$  | $$$$$$$$
  | $$     | $$  | $$__  $$
  | $$     | $$  | $$  | $$
 /$$$$$$   | $$  | $$  | $$
|______/   |__/  |__/  |__/
`

const teamBEL string = `
 /$$$$$$$  /$$$$$$$$ /$$       
| $$__  $$| $$_____/| $$       
| $$  \ $$| $$      | $$       
| $$$$$$$ | $$$$$   | $$       
| $$__  $$| $$__/   | $$       
| $$  \ $$| $$      | $$       
| $$$$$$$/| $$$$$$$$| $$$$$$$$
|_______/ |________/|________/
`

const teamCRO string = `
  /$$$$$$  /$$$$$$$   /$$$$$$  
 /$$__  $$| $$__  $$ /$$__  $$ 
| $$  \__/| $$  \ $$| $$  \ $$ 
| $$      | $$$$$$$/| $$$$$$$$ 
| $$      | $$__  $$| $$__  $$ 
| $$    $$| $$  \ $$| $$  \ $$ 
|  $$$$$$/| $$  | $$| $$  | $$ 
 \______/ |__/  |__/|__/  |__/ 
`

const teamURU string = `
 /$$   /$$ /$$$$$$$  /$$   /$$
| $$  | $$| $$__  $$| $$  | $$
| $$  | $$| $$  \ $$| $$  | $$
| $$  | $$| $$$$$$$/| $$  | $$
| $$  | $$| $$__  $$| $$  | $$
| $$  | $$| $$  \ $$| $$  | $$
|  $$$$$$/| $$  | $$|  $$$$$$/
 \______/ |__/  |__/ \______/ 
`

const teamSEN string = `
  /$$$$$$  /$$$$$$$$ /$$   /$$
 /$$__  $$| $$_____/| $$$ | $$
| $$  \__/| $$      | $$$$| $$
|  $$$$$$ | $$$$$   | $$ $$ $$
 \____  $$| $$__/   | $$  $$$$
 /$$  \ $$| $$      | $$\  $$$
|  $$$$$$/| $$$$$$$$| $$ \  $$
 \______/ |________/|__/  \__/
`

const teamECU string = `
 /$$$$$$$$ /$$$$$$  /$$   /$$
| $$_____//$$__  $$| $$  | $$
| $$     | $$  \__/| $$  | $$
| $$$$$  | $$      | $$  | $$
| $$__/  | $$      | $$  | $$
| $$     | $$    $$| $$  | $$
| $$$$$$$$|  $$$$$$/|  $$$$$$/
|________/ \______/  \______/ 
`

// ── Digits ────────────────────────────────────────────────────────────────────

const numberOne string = `
▗ 
▜ 
▟▖
`
const numberTwo string = `
▄▖
▄▌
▙▖
`
const numberThree string = `
▄▖
▄▌
▄▌
`
const numberFour string = `
▖▖
▙▌
 ▌
`
const numberFive string = `
▄▖
▙▖
▄▌
`
const numberSix string = `
▄▖
▙▖
▙▌
`
const numberSeven string = `
▄▖
 ▌
 ▌
`
const numberEight string = `
▄▖
▙▌
▙▌
`
const numberNine string = `
▄▖
▙▌
▄▌
`
const numberZero string = `
▄▖
▛▌
█▌
`
const dash string = `
▄▖
`

// ── Lookup maps ───────────────────────────────────────────────────────────────

var teamASCII = map[string]string{
	// Full FIFA codes
	"BRA": teamBRA,
	"ARG": teamARG,
	"FRA": teamFRA,
	"ENG": teamENG,
	"ESP": teamESP,
	"GER": teamGER,
	"POR": teamPOR,
	"NED": teamNED,
	"USA": teamUSA,
	"MEX": teamMEX,
	"CAN": teamCAN,
	"MAR": teamMOR, // Morocco FIFA code
	"MOR": teamMOR,
	"JPN": teamJPN,
	"KOR": teamKOR,
	"AUS": teamAUS,
	"ITA": teamITA,
	"BEL": teamBEL,
	"CRO": teamCRO,
	"URU": teamURU,
	"SEN": teamSEN,
	"ECU": teamECU,
	// Common alternates / display names
	"BRAZIL":    teamBRA,
	"ARGENTINA": teamARG,
	"FRANCE":    teamFRA,
	"ENGLAND":   teamENG,
	"SPAIN":     teamESP,
	"GERMANY":   teamGER,
	"PORTUGAL":  teamPOR,
	"NETHERLANDS": teamNED,
	"HOLLAND":   teamNED,
	"MEXICO":    teamMEX,
	"CANADA":    teamCAN,
	"MOROCCO":   teamMOR,
	"JAPAN":     teamJPN,
	"AUSTRALIA": teamAUS,
	"ITALY":     teamITA,
	"BELGIUM":   teamBEL,
	"CROATIA":   teamCRO,
	"URUGUAY":   teamURU,
	"SENEGAL":   teamSEN,
	"ECUADOR":   teamECU,
}

var digitASCII = map[rune]string{
	'0': numberZero,
	'1': numberOne,
	'2': numberTwo,
	'3': numberThree,
	'4': numberFour,
	'5': numberFive,
	'6': numberSix,
	'7': numberSeven,
	'8': numberEight,
	'9': numberNine,
	'-': dash,
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// teamArt returns the ASCII banner for a team code/name, or "" if unknown.
func teamArt(code string) string {
	art, ok := teamASCII[strings.ToUpper(code)]
	if !ok {
		return ""
	}
	return strings.TrimLeft(art, "\n")
}

// scoreArt renders a football scoreline (e.g. "2-1") using block digit art.
func scoreArt(score string) string {
	var parts []string
	for _, ch := range score {
		art, ok := digitASCII[ch]
		if !ok {
			continue
		}
		parts = append(parts, strings.TrimPrefix(art, "\n"))
	}
	if len(parts) == 0 {
		return ""
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}