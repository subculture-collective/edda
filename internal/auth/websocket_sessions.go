package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// WebSocketTicketTTL is how long a WebSocket ticket can be redeemed.
const WebSocketTicketTTL = 30 * time.Second

// WebSocketSessions issues one-use tickets that open a WebSocket for a token,
// and tracks each token's open sockets so that logout can close them.
//
// The state is in memory because the server runs as one instance. Running
// several would need a shared ticket store and a way to pass sign-outs between
// instances, such as PostgreSQL LISTEN/NOTIFY.
type WebSocketSessions struct {
	mu      sync.Mutex
	now     func() time.Time
	tickets map[string]webSocketTicket
	open    map[string]map[chan struct{}]struct{}
}

type webSocketTicket struct {
	userID  uuid.UUID
	jti     string
	expires time.Time
}

// NewWebSocketSessions creates an empty ticket store and socket registry.
func NewWebSocketSessions() *WebSocketSessions {
	return &WebSocketSessions{
		now:     time.Now,
		tickets: map[string]webSocketTicket{},
		open:    map[string]map[chan struct{}]struct{}{},
	}
}

// IssueTicket handles POST /api/v1/auth/ws-ticket. It must run behind
// JWTMiddleware. The ticket is bound to the user and the token that asked for
// it, and opens one WebSocket within WebSocketTicketTTL.
func (s *WebSocketSessions) IssueTicket(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserFromContext(r.Context())
	jti := tokenIDFromContext(r.Context())
	if !ok || jti == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	ticket, err := s.issue(userID, jti)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue ticket")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ticket": ticket, "expires_in": int(WebSocketTicketTTL.Seconds())})
}

func (s *WebSocketSessions) issue(userID uuid.UUID, jti string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	ticket := base64.RawURLEncoding.EncodeToString(b)
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	for k, t := range s.tickets {
		if !now.Before(t.expires) {
			delete(s.tickets, k)
		}
	}
	s.tickets[ticket] = webSocketTicket{userID: userID, jti: jti, expires: now.Add(WebSocketTicketTTL)}
	return ticket, nil
}

// redeem consumes a ticket. A ticket works once, and not after it expires.
func (s *WebSocketSessions) redeem(ticket string) (uuid.UUID, string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tickets[ticket]
	delete(s.tickets, ticket)
	if !ok || !s.now().Before(t.expires) {
		return uuid.Nil, "", false
	}
	return t.userID, t.jti, true
}

// track registers an open socket for jti. The returned channel closes when
// the token signs out; untrack removes the socket.
func (s *WebSocketSessions) track(jti string) (<-chan struct{}, func()) {
	ch := make(chan struct{})
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.open[jti] == nil {
		s.open[jti] = map[chan struct{}]struct{}{}
	}
	s.open[jti][ch] = struct{}{}
	return ch, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.open[jti], ch)
		if len(s.open[jti]) == 0 {
			delete(s.open, jti)
		}
	}
}

// signOut tells every open socket of jti to close.
func (s *WebSocketSessions) signOut(jti string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ch := range s.open[jti] {
		close(ch)
	}
	delete(s.open, jti)
}

type tokenContextKey struct{}

type tokenContext struct {
	jti       string
	signedOut <-chan struct{}
}

func withToken(ctx context.Context, jti string, signedOut <-chan struct{}) context.Context {
	return context.WithValue(ctx, tokenContextKey{}, tokenContext{jti: jti, signedOut: signedOut})
}

func tokenIDFromContext(ctx context.Context) string {
	t, _ := ctx.Value(tokenContextKey{}).(tokenContext)
	return t.jti
}

// SignedOut returns a channel that closes when the token that opened this
// WebSocket signs out. It is nil, and never closes, for other requests.
func SignedOut(ctx context.Context) <-chan struct{} {
	t, _ := ctx.Value(tokenContextKey{}).(tokenContext)
	return t.signedOut
}
