package sleeper

import (
	"context"
	"fmt"
)

// GetDraft fetches a single draft by ID.
func (c *Client) GetDraft(ctx context.Context, draftID string) (Draft, error) {
	var out Draft
	if err := c.get(ctx, &out, "draft", draftID); err != nil {
		return Draft{}, fmt.Errorf("fetching draft %s: %w", draftID, err)
	}
	return out, nil
}

// GetDraftPicks lists every pick made in a draft, in pick order.
func (c *Client) GetDraftPicks(ctx context.Context, draftID string) ([]DraftPick, error) {
	var out []DraftPick
	if err := c.get(ctx, &out, "draft", draftID, "picks"); err != nil {
		return nil, fmt.Errorf("fetching picks for draft %s: %w", draftID, err)
	}
	return out, nil
}
