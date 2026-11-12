package sleeper

import (
	"context"
	"fmt"
)

// GetRosters lists every roster (one per franchise) in a league.
func (c *Client) GetRosters(ctx context.Context, leagueID string) ([]Roster, error) {
	var out []Roster
	if err := c.get(ctx, &out, "league", leagueID, "rosters"); err != nil {
		return nil, fmt.Errorf("fetching rosters for league %s: %w", leagueID, err)
	}
	return out, nil
}
