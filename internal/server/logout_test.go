package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
)

func TestLogoutClearsOnlyTheSessionBeingLoggedOut(t *testing.T) {
	h, _, _, user, _, _ := ownershipFixture(t)
	// Two sign-ins by one user, usually within the same second.
	a, b := ownershipToken(t, user), ownershipToken(t, user)
	other := ownershipToken(t, uuid.New())

	for _, tc := range []struct {
		name, authorization, cookie string
		wantCleared                 bool
	}{
		{"own cookie", "Bearer " + a, a, true},
		{"cookie of a newer sign-in", "Bearer " + a, b, false},
		{"cookie of another user", "Bearer " + a, other, false},
		{"no cookie sent", "Bearer " + a, "", false},
		{"no bearer token", "", a, true},
		{"unusable authorization header", "Basic abc", a, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
			if tc.authorization != "" {
				req.Header.Set("Authorization", tc.authorization)
			}
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: auth.AuthCookieName, Value: tc.cookie})
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
			}
			cookies := rec.Result().Cookies()
			if !tc.wantCleared {
				if len(cookies) != 0 {
					t.Fatalf("logout must leave the cookie alone, got %v", cookies)
				}
				return
			}
			if len(cookies) != 1 || cookies[0].Name != auth.AuthCookieName || cookies[0].MaxAge != -1 {
				t.Fatalf("logout must clear the cookie, got %v", cookies)
			}
		})
	}
}

func TestLogoutRevokesTheTokenOnEveryRoute(t *testing.T) {
	h, _, _, user, _, campaign := ownershipFixture(t)
	signedOut, other := ownershipToken(t, user), ownershipToken(t, user)
	srv := httptest.NewServer(h)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+signedOut)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("logout = %d", resp.StatusCode)
	}

	get := func(token string) int {
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/campaigns/"+campaign.String(), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		return resp.StatusCode
	}
	if got := get(signedOut); got != http.StatusUnauthorized {
		t.Fatalf("signed-out token on a campaign route: %d, want 401", got)
	}
	if got := get(other); got != http.StatusOK {
		t.Fatalf("another session of the same user: %d, want 200", got)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/v1/campaigns/" + campaign.String() + "/ws"
	dial := func(token string) (*websocket.Conn, *http.Response, error) {
		headers := http.Header{}
		headers.Set("Cookie", auth.AuthCookieName+"="+token)
		return websocket.Dial(ctx, url, &websocket.DialOptions{HTTPHeader: headers})
	}
	conn, resp, err := dial(signedOut)
	if err == nil {
		_ = conn.CloseNow()
		t.Fatal("WebSocket upgrade accepted a signed-out token")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("signed-out WebSocket upgrade: %v %v", resp, err)
	}
	conn, _, err = dial(other)
	if err != nil {
		t.Fatalf("another session's WebSocket upgrade: %v", err)
	}
	_ = conn.CloseNow()
}
