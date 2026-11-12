package sleeper

import (
	"context"
	"fmt"
)

// GetUser looks up a Sleeper user by username or user ID.
func (c *Client) GetUser(ctx context.Context, usernameOrID string) (User, error) {
	var out User
	if err := c.get(ctx, &out, "user", usernameOrID); err != nil {
		return User{}, fmt.Errorf("fetching user %s: %w", usernameOrID, err)
	}
	return out, nil
}
