package sleeper

import (
	"context"
	"fmt"
)

// GetAllPlayers fetches Sleeper's entire player pool for a sport, keyed by
// player_id. This is a large response (tens of thousands of players,
// several MB) that Sleeper's docs say to cache and refresh at most once a
// day — callers should not invoke this per-request.
func (c *Client) GetAllPlayers(ctx context.Context, sport Sport) (map[string]Player, error) {
	var out map[string]Player
	if err := c.get(ctx, &out, "players", string(sport)); err != nil {
		return nil, fmt.Errorf("fetching players for %s: %w", sport, err)
	}
	return out, nil
}
