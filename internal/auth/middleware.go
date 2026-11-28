package auth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var errInvalidAuthorizationHeader = errors.New("invalid authorization header format")
var errMissingAuthorizationHeader = errors.New("missing authorization header")

// AuthMiddleware defines middleware that authenticates a request and enriches context.
type AuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

// ---------------------------------------------------------------------------
// NoOpMiddleware — injects a hardcoded default user (used by TUI)
// ---------------------------------------------------------------------------

// NoOpMiddleware injects a default user ID into every request context.
type NoOpMiddleware struct {
	defaultUserID uuid.UUID
}

// NewNoOpMiddleware creates an auth middleware that always authenticates as defaultUserID.
func NewNoOpMiddleware(defaultUserID uuid.UUID) AuthMiddleware {
	return &NoOpMiddleware{defaultUserID: defaultUserID}
}

// Authenticate sets the default user ID in request context and passes through.
func (m *NoOpMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithUser(r.Context(), m.defaultUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ---------------------------------------------------------------------------
// JWTMiddleware — validates Bearer tokens from the Authorization header
// ---------------------------------------------------------------------------

// JWTMiddleware validates JWT tokens and injects the user ID into context.
type JWTMiddleware struct {
	secret      string
	revocations RevocationQuerier
	sessions    *WebSocketSessions
}

// NewJWTMiddleware creates an auth middleware that validates JWT Bearer tokens,
// refuses tokens listed in revocations and redeems WebSocket tickets from
// sessions.
func NewJWTMiddleware(secret string, revocations RevocationQuerier, sessions *WebSocketSessions) AuthMiddleware {
	return &JWTMiddleware{secret: secret, revocations: revocations, sessions: sessions}
}

// Authenticate validates the bearer token, or for a WebSocket upgrade without
// one the ticket in the ticket query parameter, and sets the user ID in request
// context. Returns 401 if the credential is missing, invalid or signed out. A
// token without an ID (jti) cannot be revoked, so it is refused too: tokens
// issued before IDs were added must sign in again.
//
// A WebSocket upgrade is tracked until the handler returns, and SignedOut on
// its context closes when its token signs out.
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID uuid.UUID
		var jti string
		if ticket := webSocketTicketFromRequest(r); ticket != "" {
			var ok bool
			if userID, jti, ok = m.sessions.redeem(ticket); !ok {
				http.Error(w, `{"error":"invalid or expired ticket"}`, http.StatusUnauthorized)
				return
			}
		} else {
			token, err := TokenFromRequest(r)
			if err != nil {
				if errors.Is(err, errInvalidAuthorizationHeader) {
					http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			var claims *jwt.RegisteredClaims
			userID, claims, err = parseToken(token, m.secret)
			if err != nil || claims.ID == "" {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			jti = claims.ID
		}

		var signedOut <-chan struct{}
		if isWebSocketUpgradeRequest(r) {
			// Track before checking: a logout that lands in between is then
			// either seen by the check or closes this socket.
			var untrack func()
			signedOut, untrack = m.sessions.track(jti)
			defer untrack()
		}
		revoked, err := m.revocations.IsTokenRevoked(r.Context(), jti)
		if err != nil {
			http.Error(w, `{"error":"could not check the session"}`, http.StatusServiceUnavailable)
			return
		}
		if revoked {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := withToken(WithUser(r.Context(), userID), jti, signedOut)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
