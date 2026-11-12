package sleeper

import (
	"context"
	"fmt"
)

// GetLeague fetches a single league by ID.
func (c *Client) GetLeague(ctx context.Context, leagueID string) (League, error) {
	var out League
	if err := c.get(ctx, &out, "league", leagueID); err != nil {
		return League{}, fmt.Errorf("fetching league %s: %w", leagueID, err)
	}
	return out, nil
}

// GetLeaguesForUser lists all leagues a user belongs to for a given sport and season.
func (c *Client) GetLeaguesForUser(ctx context.Context, userID string, sport Sport, season string) ([]League, error) {
	var out []League
	if err := c.get(ctx, &out, "user", userID, "leagues", string(sport), season); err != nil {
		return nil, fmt.Errorf("fetching leagues for user %s: %w", userID, err)
	}
	return out, nil
}

// GetLeagueUsers lists every member of a league.
func (c *Client) GetLeagueUsers(ctx context.Context, leagueID string) ([]LeagueUser, error) {
	var out []LeagueUser
	if err := c.get(ctx, &out, "league", leagueID, "users"); err != nil {
		return nil, fmt.Errorf("fetching users for league %s: %w", leagueID, err)
	}
	return out, nil
}
