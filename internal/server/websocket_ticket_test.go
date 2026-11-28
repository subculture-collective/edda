package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"git.subcult.tv/subculture-collective/edda/internal/handlers"
)

// wsTicket asks the server for a WebSocket ticket the way the web client does.
func wsTicket(t *testing.T, baseURL, token string) string {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/api/v1/auth/ws-ticket", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var body struct{ Ticket string }
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || resp.StatusCode != http.StatusOK || body.Ticket == "" {
		t.Fatalf("ticket = %d %+v: %v", resp.StatusCode, body, err)
	}
	return body.Ticket
}

func TestWebSocketTicketThroughTheRouter(t *testing.T) {
	h, _, _, user, _, campaign := ownershipFixture(t)
	srv := httptest.NewServer(h)
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/v1/campaigns/" + campaign.String() + "/ws"
	token := ownershipToken(t, user)

	ticket := wsTicket(t, srv.URL, token)
	conn, _, err := websocket.Dial(ctx, url+"?ticket="+ticket, nil)
	if err != nil {
		t.Fatalf("dial with a ticket: %v", err)
	}
	_ = conn.CloseNow()

	for name, header := range map[string]http.Header{
		"reused ticket": nil,
		"old cookie":    {"Cookie": []string{"gm_token=" + token}},
	} {
		target := url
		if header == nil {
			target += "?ticket=" + ticket
		}
		conn, resp, err := websocket.Dial(ctx, target, &websocket.DialOptions{HTTPHeader: header})
		if err == nil {
			_ = conn.CloseNow()
			t.Fatalf("%s: upgrade accepted", name)
		}
		if resp == nil || resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("%s: %v %v, want 401", name, resp, err)
		}
	}
}

func TestLogoutClosesOnlyThatTokensOpenSockets(t *testing.T) {
	h, _, _, user, _, campaign := ownershipFixture(t)
	srv := httptest.NewServer(h)
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/v1/campaigns/" + campaign.String() + "/ws"
	signedOut, other := ownershipToken(t, user), ownershipToken(t, user)

	open := func(token string) *websocket.Conn {
		t.Helper()
		conn, _, err := websocket.Dial(ctx, url+"?ticket="+wsTicket(t, srv.URL, token), nil)
		if err != nil {
			t.Fatal(err)
		}
		return conn
	}
	closed, stays := open(signedOut), open(other)
	defer func() { _ = closed.CloseNow(); _ = stays.CloseNow() }()
	staysCtx := stays.CloseRead(ctx)

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

	_, _, err = closed.Read(ctx)
	if got := websocket.CloseStatus(err); got != handlers.StatusSignedOut {
		t.Fatalf("signed-out socket: %v (status %d), want close status %d", err, got, handlers.StatusSignedOut)
	}
	if err := stays.Ping(staysCtx); err != nil {
		t.Fatalf("another session's socket: %v", err)
	}
}
