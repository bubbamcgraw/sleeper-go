package sleeper

import (
	"context"
	"fmt"
)

// GetWinnersBracket fetches a league's championship playoff bracket.
func (c *Client) GetWinnersBracket(ctx context.Context, leagueID string) ([]BracketMatchup, error) {
	var out []BracketMatchup
	if err := c.get(ctx, &out, "league", leagueID, "winners_bracket"); err != nil {
		return nil, fmt.Errorf("fetching winners bracket for league %s: %w", leagueID, err)
	}
	return out, nil
}

// GetLosersBracket fetches a league's consolation ("toilet bowl") bracket.
func (c *Client) GetLosersBracket(ctx context.Context, leagueID string) ([]BracketMatchup, error) {
	var out []BracketMatchup
	if err := c.get(ctx, &out, "league", leagueID, "losers_bracket"); err != nil {
		return nil, fmt.Errorf("fetching losers bracket for league %s: %w", leagueID, err)
	}
	return out, nil
}
