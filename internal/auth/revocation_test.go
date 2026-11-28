package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// memRevocations is an in-memory RevocationQuerier.
type memRevocations struct {
	mu      sync.Mutex
	revoked map[string]time.Time
	err     error
}

func newMemRevocations() *memRevocations { return &memRevocations{revoked: map[string]time.Time{}} }

func (m *memRevocations) RevokeToken(_ context.Context, arg statedb.RevokeTokenParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return m.err
	}
	m.revoked[arg.Jti] = arg.ExpiresAt.Time
	return nil
}

func (m *memRevocations) IsTokenRevoked(_ context.Context, jti string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return false, m.err
	}
	_, ok := m.revoked[jti]
	return ok, nil
}

func revocationFixture(t *testing.T) (*memRevocations, http.Handler, *AuthHandlers) {
	t.Helper()
	rev := newMemRevocations()
	sessions := NewWebSocketSessions()
	protected := NewJWTMiddleware("secret", rev, sessions).Authenticate(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	return rev, protected, NewAuthHandlers(nil, rev, sessions, "secret")
}

func mustToken(t *testing.T, userID uuid.UUID) string {
	t.Helper()
	token, err := GenerateToken(userID, "secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func bearerStatus(h http.Handler, token string) int {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

func webSocketStatus(h http.Handler, token string) int {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/x/ws", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

func logout(t *testing.T, h *AuthHandlers, bearer string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.Logout(rec, req)
	return rec
}

func TestLogoutRevokesTheBearerTokenOnly(t *testing.T) {
	_, protected, h := revocationFixture(t)
	user := uuid.New()
	a, b := mustToken(t, user), mustToken(t, user)

	if rec := logout(t, h, a); rec.Code != http.StatusNoContent {
		t.Fatalf("logout = %d %s", rec.Code, rec.Body.String())
	}
	if got := bearerStatus(protected, a); got != http.StatusUnauthorized {
		t.Fatalf("signed-out bearer token: status %d, want 401", got)
	}
	if got := webSocketStatus(protected, a); got != http.StatusUnauthorized {
		t.Fatalf("signed-out WebSocket upgrade: status %d, want 401", got)
	}
	if got := bearerStatus(protected, b); got != http.StatusNoContent {
		t.Fatalf("other session of the same user: status %d, want 204", got)
	}
	if got := webSocketStatus(protected, b); got != http.StatusNoContent {
		t.Fatalf("other session's WebSocket: status %d, want 204", got)
	}
}

func TestLogoutIgnoresTheOldSessionCookie(t *testing.T) {
	rev, _, h := revocationFixture(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "gm_token", Value: mustToken(t, uuid.New())})
	rec := httptest.NewRecorder()
	h.Logout(rec, req)
	if rec.Code != http.StatusNoContent || len(rec.Result().Cookies()) != 0 || len(rev.revoked) != 0 {
		t.Fatalf("logout = %d, cookies %v, revoked %v", rec.Code, rec.Result().Cookies(), rev.revoked)
	}
}

func TestLogoutIgnoresATokenItCannotVerify(t *testing.T) {
	rev, _, h := revocationFixture(t)
	forged, err := GenerateToken(uuid.New(), "other-secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if rec := logout(t, h, forged); rec.Code != http.StatusNoContent {
		t.Fatalf("logout = %d", rec.Code)
	}
	if len(rev.revoked) != 0 {
		t.Fatalf("an unverified token must not be stored, got %v", rev.revoked)
	}
}

func TestLogoutReportsAFailedRevocation(t *testing.T) {
	rev, _, h := revocationFixture(t)
	rev.err = errors.New("database down")
	rec := logout(t, h, mustToken(t, uuid.New()))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("logout = %d, want 500", rec.Code)
	}
}

func TestMiddlewareRefusesATokenWithoutID(t *testing.T) {
	_, protected, _ := revocationFixture(t)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	if got := bearerStatus(protected, signed); got != http.StatusUnauthorized {
		t.Fatalf("token without jti: status %d, want 401", got)
	}
}

func TestMiddlewareFailsClosedWhenTheRevocationCheckFails(t *testing.T) {
	rev, protected, _ := revocationFixture(t)
	rev.err = errors.New("database down")
	if got := bearerStatus(protected, mustToken(t, uuid.New())); got != http.StatusServiceUnavailable {
		t.Fatalf("status %d, want 503", got)
	}
}
