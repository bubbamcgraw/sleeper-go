package sleeper

import (
	"context"
	"testing"
)

// TestEndpoints checks each method hits the right path and decodes Sleeper's
// response shape.
func TestEndpoints(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		wantPath string
		body     string
		call     func(t *testing.T, c *Client) error
	}{
		{
			name:     "GetUser",
			wantPath: "/v1/user/sleeperuser",
			body:     `{"user_id":"12345678","username":"sleeperuser","display_name":"SleeperUser","avatar":"cc12ec49965eb7856f84d71cf85306af"}`,
			call: func(t *testing.T, c *Client) error {
				u, err := c.GetUser(ctx, "sleeperuser")
				expect(t, "UserID", u.UserID, "12345678")
				expect(t, "DisplayName", u.DisplayName, "SleeperUser")
				return err
			},
		},
		{
			name:     "GetLeague",
			wantPath: "/v1/league/289646328504385536",
			body:     `{"league_id":"289646328504385536","name":"Sleeperbot Dynasty","season":"2018","total_rosters":12,"draft_id":"289646328508579840","roster_positions":["QB","RB"],"settings":{"playoff_teams":6}}`,
			call: func(t *testing.T, c *Client) error {
				l, err := c.GetLeague(ctx, "289646328504385536")
				expect(t, "Name", l.Name, "Sleeperbot Dynasty")
				expect(t, "TotalRosters", l.TotalRosters, 12)
				expect(t, "DraftID", l.DraftID, "289646328508579840")
				expect(t, "len(RosterPositions)", len(l.RosterPositions), 2)
				expect(t, "Settings.playoff_teams", l.Settings["playoff_teams"], any(float64(6)))
				return err
			},
		},
		{
			name:     "GetLeaguesForUser",
			wantPath: "/v1/user/12345678/leagues/nfl/2018",
			body:     `[{"league_id":"289646328504385536","name":"Sleeperbot Dynasty"}]`,
			call: func(t *testing.T, c *Client) error {
				ls, err := c.GetLeaguesForUser(ctx, "12345678", NFL, "2018")
				expect(t, "len", len(ls), 1)
				return err
			},
		},
		{
			name:     "GetLeagueUsers",
			wantPath: "/v1/league/289646328504385536/users",
			body:     `[{"user_id":"457511950237696","display_name":"2KSports","is_owner":true,"metadata":{"team_name":"Giant Dolphins"}}]`,
			call: func(t *testing.T, c *Client) error {
				us, err := c.GetLeagueUsers(ctx, "289646328504385536")
				if len(us) == 1 {
					expect(t, "IsOwner", us[0].IsOwner, true)
					expect(t, "team_name", us[0].Metadata["team_name"], any("Giant Dolphins"))
				} else {
					t.Errorf("len = %d, want 1", len(us))
				}
				return err
			},
		},
		{
			name:     "GetRosters",
			wantPath: "/v1/league/289646328504385536/rosters",
			body:     `[{"roster_id":1,"owner_id":"12345678","starters":["2307","2257"],"players":["1046","138"],"reserve":null}]`,
			call: func(t *testing.T, c *Client) error {
				rs, err := c.GetRosters(ctx, "289646328504385536")
				if len(rs) == 1 {
					expect(t, "RosterID", rs[0].RosterID, 1)
					expect(t, "len(Starters)", len(rs[0].Starters), 2)
				} else {
					t.Errorf("len = %d, want 1", len(rs))
				}
				return err
			},
		},
		{
			name:     "GetMatchups",
			wantPath: "/v1/league/289646328504385536/matchups/3",
			body:     `[{"roster_id":1,"matchup_id":2,"points":120.5,"starters_points":[20.1,10],"players_points":{"1046":20.1}}]`,
			call: func(t *testing.T, c *Client) error {
				ms, err := c.GetMatchups(ctx, "289646328504385536", 3)
				if len(ms) == 1 {
					expect(t, "Points", ms[0].Points, 120.5)
					expect(t, "PlayersPoints[1046]", ms[0].PlayersPoints["1046"], 20.1)
				} else {
					t.Errorf("len = %d, want 1", len(ms))
				}
				return err
			},
		},
		{
			name:     "GetWinnersBracket",
			wantPath: "/v1/league/289646328504385536/winners_bracket",
			body:     `[{"r":1,"m":1,"t1":3,"t2":6,"w":3,"l":6},{"r":2,"m":3,"t1_from":{"w":1},"t2_from":{"l":2},"p":3}]`,
			call: func(t *testing.T, c *Client) error {
				b, err := c.GetWinnersBracket(ctx, "289646328504385536")
				if len(b) == 2 {
					expect(t, "[0].Winner", *b[0].Winner, 3)
					expect(t, "[1].Team1", b[1].Team1, (*int)(nil))
					expect(t, "[1].Team1From.Winner", *b[1].Team1From.Winner, 1)
					expect(t, "[1].Team2From.Loser", *b[1].Team2From.Loser, 2)
					expect(t, "[1].Position", *b[1].Position, 3)
				} else {
					t.Errorf("len = %d, want 2", len(b))
				}
				return err
			},
		},
		{
			name:     "GetLosersBracket",
			wantPath: "/v1/league/289646328504385536/losers_bracket",
			body:     `[{"r":1,"m":1,"t1":1,"t2":2}]`,
			call: func(t *testing.T, c *Client) error {
				b, err := c.GetLosersBracket(ctx, "289646328504385536")
				expect(t, "len", len(b), 1)
				return err
			},
		},
		{
			name:     "GetDraft",
			wantPath: "/v1/draft/257270643320426496",
			body:     `{"draft_id":"257270643320426496","type":"snake","status":"complete","start_time":1515700800000,"draft_order":{"12345678":1}}`,
			call: func(t *testing.T, c *Client) error {
				d, err := c.GetDraft(ctx, "257270643320426496")
				expect(t, "Type", d.Type, "snake")
				expect(t, "StartTime", d.StartTime, int64(1515700800000))
				return err
			},
		},
		{
			name:     "GetDraftPicks",
			wantPath: "/v1/draft/257270643320426496/picks",
			body:     `[{"player_id":"2391","picked_by":"234343434","roster_id":1,"round":5,"draft_slot":5,"pick_no":1,"is_keeper":null},{"pick_no":2,"is_keeper":true}]`,
			call: func(t *testing.T, c *Client) error {
				ps, err := c.GetDraftPicks(ctx, "257270643320426496")
				if len(ps) == 2 {
					expect(t, "[0].Round", ps[0].Round, 5)
					expect(t, "[0].IsKeeper", ps[0].IsKeeper, (*bool)(nil))
					expect(t, "[1].IsKeeper", *ps[1].IsKeeper, true)
				} else {
					t.Errorf("len = %d, want 2", len(ps))
				}
				return err
			},
		},
		{
			name:     "GetAllPlayers",
			wantPath: "/v1/players/nfl",
			body:     `{"3086":{"player_id":"3086","espn_id":2577327,"full_name":"Tom Brady","position":"QB","fantasy_positions":["QB"],"active":true}}`,
			call: func(t *testing.T, c *Client) error {
				ps, err := c.GetAllPlayers(ctx, NFL)
				expect(t, "FullName", ps["3086"].FullName, "Tom Brady")
				expect(t, "EspnID", ps["3086"].EspnID, 2577327)
				return err
			},
		},
		{
			name:     "GetState",
			wantPath: "/v1/state/nfl",
			body:     `{"week":2,"season_type":"regular","season":"2020","display_week":3}`,
			call: func(t *testing.T, c *Client) error {
				s, err := c.GetState(ctx, NFL)
				expect(t, "Week", s.Week, 2)
				expect(t, "DisplayWeek", s.DisplayWeek, 3)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, tt.wantPath, 200, tt.body)
			if err := tt.call(t, c); err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
		})
	}
}

func expect[T comparable](t *testing.T, field string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}
