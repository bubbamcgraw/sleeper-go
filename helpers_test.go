package sleeper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient returns a Client pointed at a server that asserts the request
// path (escaped form) and replies with status and body.
func newTestClient(t *testing.T, wantPath string, status int, body string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.URL.EscapedPath(); got != wantPath {
			t.Errorf("path = %q, want %q", got, wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(srv.Client())
	c.BaseURL = srv.URL + "/v1"
	return c
}
