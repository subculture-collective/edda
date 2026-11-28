package auth

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

const AuthCookieName = "gm_token"

var errInvalidAuthorizationHeader = errors.New("invalid authorization header format")

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
}

// NewJWTMiddleware creates an auth middleware that validates JWT Bearer tokens
// and refuses tokens listed in revocations.
func NewJWTMiddleware(secret string, revocations RevocationQuerier) AuthMiddleware {
	return &JWTMiddleware{secret: secret, revocations: revocations}
}

// Authenticate reads the Authorization header, validates the JWT, and sets
// the user ID in request context. Returns 401 if the token is missing, invalid
// or signed out. A token without an ID (jti) cannot be revoked, so it is
// refused too: tokens issued before IDs were added must sign in again.
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := tokenFromRequest(r)
		if err != nil {
			if errors.Is(err, errInvalidAuthorizationHeader) {
				http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		userID, claims, err := parseToken(token, m.secret)
		if err != nil || claims.ID == "" {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}
		revoked, err := m.revocations.IsTokenRevoked(r.Context(), claims.ID)
		if err != nil {
			http.Error(w, `{"error":"could not check the session"}`, http.StatusServiceUnavailable)
			return
		}
		if revoked {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := WithUser(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tokenFromRequest(r *http.Request) (string, error) {
	return TokenFromRequest(r)
}
