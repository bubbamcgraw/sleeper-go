package sleeper

import (
	"context"
	"fmt"
	"strconv"
)

// GetMatchups lists every roster's performance for a single week of a
// league — rosters sharing the same Matchup.MatchupID played each other
// that week.
func (c *Client) GetMatchups(ctx context.Context, leagueID string, week int) ([]Matchup, error) {
	var out []Matchup
	if err := c.get(ctx, &out, "league", leagueID, "matchups", strconv.Itoa(week)); err != nil {
		return nil, fmt.Errorf("fetching matchups for league %s week %d: %w", leagueID, week, err)
	}
	return out, nil
}
