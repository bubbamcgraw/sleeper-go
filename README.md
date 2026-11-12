# sleeper-go

Go client for the read-only [Sleeper API](https://docs.sleeper.com). No dependencies.

```sh
go get github.com/bubbamcgraw/sleeper-go
```

## Usage

```go
import "github.com/bubbamcgraw/sleeper-go"

client := sleeper.NewClient(nil) // or pass your own *http.Client
ctx := context.Background()

state, err := client.GetState(ctx, sleeper.NFL)
user, err := client.GetUser(ctx, "sleeperuser")
leagues, err := client.GetLeaguesForUser(ctx, user.UserID, sleeper.NFL, state.Season)
league, err := client.GetLeague(ctx, "289646328504385536")
```

| Method | Endpoint |
| --- | --- |
| `GetUser(ctx, usernameOrID)` | `/user/<id>` |
| `GetLeaguesForUser(ctx, userID, sport, season)` | `/user/<id>/leagues/<sport>/<season>` |
| `GetLeague(ctx, leagueID)` | `/league/<id>` |
| `GetLeagueUsers(ctx, leagueID)` | `/league/<id>/users` |
| `GetRosters(ctx, leagueID)` | `/league/<id>/rosters` |
| `GetMatchups(ctx, leagueID, week)` | `/league/<id>/matchups/<week>` |
| `GetWinnersBracket(ctx, leagueID)` | `/league/<id>/winners_bracket` |
| `GetLosersBracket(ctx, leagueID)` | `/league/<id>/losers_bracket` |
| `GetDraft(ctx, draftID)` | `/draft/<id>` |
| `GetDraftPicks(ctx, draftID)` | `/draft/<id>/picks` |
| `GetAllPlayers(ctx, sport)` | `/players/<sport>` — several MB, cache it daily |
| `GetState(ctx, sport)` | `/state/<sport>` |

### Not found

Sleeper returns either a 404 or a 200 with a `null` body for missing resources.
Both come back as `sleeper.ErrNotFound`:

```go
if errors.Is(err, sleeper.ErrNotFound) { ... }
```

## Tests

```sh
go test -short ./...  # mocked only
go test ./...         # also hits the live API
```

## License

MIT
