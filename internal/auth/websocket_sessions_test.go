package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

type ticketFixture struct {
	sessions  *WebSocketSessions
	protected http.Handler
	issue     http.Handler
	handlers  *AuthHandlers
	seen      uuid.UUID
}

func newTicketFixture(t *testing.T) *ticketFixture {
	t.Helper()
	f := &ticketFixture{sessions: NewWebSocketSessions()}
	rev := newMemRevocations()
	mw := NewJWTMiddleware("secret", rev, f.sessions)
	f.protected = mw.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.seen, _ = UserFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))
	f.issue = mw.Authenticate(http.HandlerFunc(f.sessions.IssueTicket))
	f.handlers = NewAuthHandlers(nil, rev, f.sessions, "secret")
	return f
}

func (f *ticketFixture) ticket(t *testing.T, token string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/ws-ticket", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	f.issue.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("issue ticket = %d %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Ticket    string `json:"ticket"`
		ExpiresIn int    `json:"expires_in"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil || body.Ticket == "" || body.ExpiresIn != int(WebSocketTicketTTL.Seconds()) {
		t.Fatalf("ticket reply %+v: %v", body, err)
	}
	return body.Ticket
}

func (f *ticketFixture) upgrade(ticket string) int {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/x/ws?ticket="+ticket, nil)
	req.Header.Set("Upgrade", "websocket")
	rec := httptest.NewRecorder()
	f.protected.ServeHTTP(rec, req)
	return rec.Code
}

func TestWebSocketTicketOpensOneUpgrade(t *testing.T) {
	f := newTicketFixture(t)
	user := uuid.New()
	ticket := f.ticket(t, mustToken(t, user))

	if got := f.upgrade(ticket); got != http.StatusNoContent || f.seen != user {
		t.Fatalf("first use: status %d user %s, want 204 %s", got, f.seen, user)
	}
	if got := f.upgrade(ticket); got != http.StatusUnauthorized {
		t.Fatalf("reused ticket: status %d, want 401", got)
	}
}

func TestWebSocketTicketExpires(t *testing.T) {
	f := newTicketFixture(t)
	now := time.Now()
	f.sessions.now = func() time.Time { return now }
	ticket := f.ticket(t, mustToken(t, uuid.New()))

	now = now.Add(WebSocketTicketTTL)
	if got := f.upgrade(ticket); got != http.StatusUnauthorized {
		t.Fatalf("expired ticket: status %d, want 401", got)
	}
}

func TestWebSocketTicketOfASignedOutTokenIsRefused(t *testing.T) {
	f := newTicketFixture(t)
	token := mustToken(t, uuid.New())
	ticket := f.ticket(t, token)

	if rec := logout(t, f.handlers, token); rec.Code != http.StatusNoContent {
		t.Fatalf("logout = %d", rec.Code)
	}
	if got := f.upgrade(ticket); got != http.StatusUnauthorized {
		t.Fatalf("ticket of a signed-out token: status %d, want 401", got)
	}
}

func TestWebSocketTicketOnlyOpensAnUpgrade(t *testing.T) {
	f := newTicketFixture(t)
	ticket := f.ticket(t, mustToken(t, uuid.New()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/?ticket="+ticket, nil)
	rec := httptest.NewRecorder()
	f.protected.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ticket on a plain request: status %d, want 401", rec.Code)
	}
}

func TestWebSocketTicketNeedsABearerToken(t *testing.T) {
	f := newTicketFixture(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/ws-ticket", nil)
	rec := httptest.NewRecorder()
	f.issue.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ticket without a token: status %d, want 401", rec.Code)
	}
}

func TestSignOutClosesOnlyThatTokensSockets(t *testing.T) {
	s := NewWebSocketSessions()
	a1, untrackA1 := s.track("a")
	a2, untrackA2 := s.track("a")
	b, untrackB := s.track("b")
	defer untrackB()

	s.signOut("a")
	for i, ch := range []<-chan struct{}{a1, a2} {
		select {
		case <-ch:
		default:
			t.Fatalf("socket %d of the signed-out token is still open", i)
		}
	}
	select {
	case <-b:
		t.Fatal("another token's socket was closed")
	default:
	}
	// Untracking after sign-out must not panic or close anything twice.
	untrackA1()
	untrackA2()
	s.signOut("a")
}
