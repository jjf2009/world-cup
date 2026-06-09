# Data Schemas

The Phase 1 JSON repositories load these files once at startup and normalize them
into typed Go domain models.

## football.teams.json

```json
{
  "id": "1",
  "name_en": "Mexico",
  "flag": "https://flagcdn.com/w80/mx.png",
  "fifa_code": "MEX",
  "iso2": "MX",
  "groups": "A"
}
```

## football.stadiums.json

```json
{
  "id": "1",
  "name_en": "Estadio Azteca",
  "fifa_name": "Mexico City Stadium",
  "city_en": "Mexico City",
  "country_en": "Mexico",
  "capacity": 83000,
  "region": "Central"
}
```

## football.matches.json

```json
{
  "id": "1",
  "home_team_id": "1",
  "away_team_id": "2",
  "home_score": "0",
  "away_score": "0",
  "home_scorers": "null",
  "away_scorers": "null",
  "group": "A",
  "matchday": "1",
  "local_date": "06/11/2026 13:00",
  "stadium_id": "1",
  "finished": "FALSE",
  "time_elapsed": "notstarted",
  "type": "group"
}
```

Knockout fixtures may use team ID `"0"` with `home_team_label` or
`away_team_label`, for example `"Winner Match 101"`.

## football.matchtables.json

```json
{
  "group": "A",
  "teams": [
    {
      "team_id": "1",
      "mp": "0",
      "w": "0",
      "l": "0",
      "d": "0",
      "pts": "0",
      "gf": "0",
      "ga": "0",
      "gd": "0"
    }
  ]
}
```

## winners.json

```json
{
  "year": 2022,
  "winner": "Argentina",
  "runner_up": "France",
  "venue": "Lusail Stadium"
}
```

`venue` is optional.

## Normalization Rules

- Numeric fields represented as strings are parsed into integers.
- `finished` accepts `"TRUE"` as true; everything else is false.
- Scorer values of `"null"` or an empty string become an empty scorer list.
- Match dates are parsed with `MM/DD/YYYY HH:mm`.
- The Bubble Tea UI receives resolved view models, not raw IDs.
