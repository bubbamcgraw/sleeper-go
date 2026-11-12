package sleeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://api.sleeper.app/v1"

// ErrNotFound is returned when Sleeper has no such resource. Sleeper signals
// this either with a 404 or with a 200 and a JSON null body, depending on the
// endpoint; both map here. Check with errors.Is.
var ErrNotFound = errors.New("sleeper: not found")

// Client wraps an *http.Client and centralizes the fetch/decode logic shared
// by every Sleeper endpoint.
type Client struct {
	HTTP    *http.Client
	BaseURL string
}

// NewClient builds a Client against DefaultBaseURL. A nil httpClient falls
// back to one with a 30s timeout.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{HTTP: httpClient, BaseURL: DefaultBaseURL}
}

// get issues a GET to BaseURL joined with the escaped path segments and
// decodes the JSON response body into out.
func (c *Client) get(ctx context.Context, out any, segments ...string) error {
	escaped := make([]string, len(segments))
	for i, s := range segments {
		escaped[i] = url.PathEscape(s)
	}
	u := c.BaseURL + "/" + strings.Join(escaped, "/")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("building request %s: %w", u, err)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, u)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading %s: %w", u, err)
	}
	if bytes.Equal(bytes.TrimSpace(body), []byte("null")) {
		return ErrNotFound
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decoding %s: %w", u, err)
	}

	return nil
}
