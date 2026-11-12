package sleeper

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantErr    string
		wantNotFnd bool
	}{
		{name: "ok", status: 200, body: `{"user_id":"1"}`},
		{name: "404", status: 404, body: `null`, wantNotFnd: true},
		{name: "200 null", status: 200, body: "null\n", wantNotFnd: true},
		{name: "500", status: 500, body: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", status: 200, body: `{`, wantErr: "decoding"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, "/v1/user/1", tt.status, tt.body)

			var out User
			err := c.get(context.Background(), &out, "user", "1")

			switch {
			case tt.wantNotFnd:
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("err = %v, want ErrNotFound", err)
				}
			case tt.wantErr != "":
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				if out.UserID != "1" {
					t.Errorf("UserID = %q, want %q", out.UserID, "1")
				}
			}
		})
	}
}

func TestGet_EscapesSegments(t *testing.T) {
	c := newTestClient(t, "/v1/user/has%20space%2Fslash", 200, `{"user_id":"1"}`)

	if _, err := c.GetUser(context.Background(), "has space/slash"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestGet_NotFoundWrappedByEndpoint(t *testing.T) {
	c := newTestClient(t, "/v1/league/0", 404, `null`)

	_, err := c.GetLeague(context.Background(), "0")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if !strings.Contains(err.Error(), "league 0") {
		t.Errorf("err = %q, want league id in message", err)
	}
}

func TestGet_ContextCanceled(t *testing.T) {
	c := newTestClient(t, "/v1/user/1", 200, `{}`)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetUser(ctx, "1")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient(nil)
	if c.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", c.BaseURL, DefaultBaseURL)
	}
	if c.HTTP == nil || c.HTTP.Timeout == 0 {
		t.Error("expected default HTTP client with a timeout")
	}

	custom := &http.Client{}
	if NewClient(custom).HTTP != custom {
		t.Error("expected custom HTTP client to be used")
	}
}

func TestGet_TransportError(t *testing.T) {
	c := NewClient(nil)
	c.BaseURL = "http://127.0.0.1:0"

	_, err := c.GetUser(context.Background(), "1")
	if err == nil || !strings.Contains(err.Error(), "fetching") {
		t.Fatalf("err = %v, want fetch error", err)
	}
}
