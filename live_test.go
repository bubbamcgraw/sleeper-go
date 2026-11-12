package sleeper

import (
	"context"
	"errors"
	"testing"
)

// Live tests hit the real Sleeper API using the example IDs from Sleeper's
// docs. Slow, network-dependent, and will break if Sleeper changes the
// response shape — skipped with -short.
const (
	liveLeagueID = "289646328504385536"
	liveDraftID  = "257270643320426496"
	// owner of liveLeagueID; the docs' example user_id 12345678 owns nothing
	liveUserID = "457511950237696"
)

func liveClient(t *testing.T) (*Client, context.Context) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping live network test in -short mode")
	}
	return NewClient(nil), context.Background()
}

func TestLive_User(t *testing.T) {
	c, ctx := liveClient(t)

	u, err := c.GetUser(ctx, liveUserID)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "UserID", u.UserID, liveUserID)

	if _, err := c.GetUser(ctx, "doesnotexist_zzq_sleeper_go"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown user err = %v, want ErrNotFound", err)
	}
}

func TestLive_League(t *testing.T) {
	c, ctx := liveClient(t)

	l, err := c.GetLeague(ctx, liveLeagueID)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "LeagueID", l.LeagueID, liveLeagueID)

	if _, err := c.GetLeague(ctx, "0"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown league err = %v, want ErrNotFound", err)
	}

	users, err := c.GetLeagueUsers(ctx, liveLeagueID)
	if err != nil || len(users) == 0 {
		t.Errorf("GetLeagueUsers: %d users, err %v", len(users), err)
	}

	rosters, err := c.GetRosters(ctx, liveLeagueID)
	if err != nil || len(rosters) == 0 {
		t.Errorf("GetRosters: %d rosters, err %v", len(rosters), err)
	}

	if _, err := c.GetMatchups(ctx, liveLeagueID, 1); err != nil {
		t.Errorf("GetMatchups: %v", err)
	}
	if _, err := c.GetWinnersBracket(ctx, liveLeagueID); err != nil {
		t.Errorf("GetWinnersBracket: %v", err)
	}
	if _, err := c.GetLosersBracket(ctx, liveLeagueID); err != nil {
		t.Errorf("GetLosersBracket: %v", err)
	}
}

func TestLive_LeaguesForUser(t *testing.T) {
	c, ctx := liveClient(t)

	leagues, err := c.GetLeaguesForUser(ctx, liveUserID, NFL, "2018")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range leagues {
		if l.LeagueID == liveLeagueID {
			return
		}
	}
	t.Errorf("expected league %s in results", liveLeagueID)
}

func TestLive_Draft(t *testing.T) {
	c, ctx := liveClient(t)

	d, err := c.GetDraft(ctx, liveDraftID)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, "DraftID", d.DraftID, liveDraftID)

	picks, err := c.GetDraftPicks(ctx, liveDraftID)
	if err != nil || len(picks) == 0 {
		t.Errorf("GetDraftPicks: %d picks, err %v", len(picks), err)
	}
}

func TestLive_State(t *testing.T) {
	c, ctx := liveClient(t)

	s, err := c.GetState(ctx, NFL)
	if err != nil {
		t.Fatal(err)
	}
	if s.Season == "" {
		t.Error("expected Season to be populated")
	}
}
