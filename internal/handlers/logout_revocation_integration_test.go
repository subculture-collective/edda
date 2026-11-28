//go:build integration

package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/coder/websocket"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/server"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

func TestIntegrationLogoutRevokesToken(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	campaign, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Logout", Status: "active", CreatedBy: db.ToPgUUID(testUserID), Themes: []string{}})
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{}
	cfg.Server.JWTSecret = "integration-logout"
	srv := httptest.NewServer(server.NewRouterWithProvider(log.New(io.Discard), &stubIntegrationEngine{}, q, nil, testPool, testUserID, cfg, nil))
	defer srv.Close()
	token := func() string {
		t.Helper()
		tok, err := auth.GenerateToken(testUserID, cfg.Server.JWTSecret, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		return tok
	}
	do := func(method, path, tok string) int {
		req, _ := http.NewRequest(method, srv.URL+path, nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return 0
		}
		_ = resp.Body.Close()
		return resp.StatusCode
	}
	signedOut, other := token(), token()
	campaignPath := "/api/v1/campaigns/" + db.FromPgUUID(campaign.ID).String()

	// Concurrent logouts of one token all succeed and leave one row.
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := do(http.MethodPost, "/api/v1/auth/logout", signedOut); got != http.StatusNoContent {
				t.Errorf("concurrent logout = %d", got)
			}
		}()
	}
	wg.Wait()

	if got := do(http.MethodGet, campaignPath, signedOut); got != http.StatusUnauthorized {
		t.Fatalf("signed-out token: %d, want 401", got)
	}
	if got := do(http.MethodGet, campaignPath, other); got != http.StatusOK {
		t.Fatalf("another session of the same user: %d, want 200", got)
	}

	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + campaignPath + "/ws"
	dial := func(tok string) (*websocket.Conn, *http.Response, error) {
		h := http.Header{}
		h.Set("Cookie", auth.AuthCookieName+"="+tok)
		return websocket.Dial(dialCtx, url, &websocket.DialOptions{HTTPHeader: h})
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
