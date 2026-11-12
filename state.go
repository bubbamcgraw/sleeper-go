package sleeper

import (
	"context"
	"fmt"
)

// GetState returns Sleeper's current view of the season/week for a sport,
// which callers typically use to figure out which week's matchups are current.
func (c *Client) GetState(ctx context.Context, sport Sport) (State, error) {
	var out State
	if err := c.get(ctx, &out, "state", string(sport)); err != nil {
		return State{}, fmt.Errorf("fetching state for %s: %w", sport, err)
	}
	return out, nil
}
