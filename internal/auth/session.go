package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type sessionContextKey struct{}

// WithUser stores an authenticated user ID in context.
func WithUser(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionContextKey{}, userID)
}

// UserFromContext extracts an authenticated user ID from context.
// A uuid.Nil value is treated as unauthenticated and returns ok=false.
func UserFromContext(ctx context.Context) (uuid.UUID, bool) {
	if ctx == nil {
		return uuid.Nil, false
	}
	switch userID := ctx.Value(sessionContextKey{}).(type) {
	case uuid.UUID:
		if userID == uuid.Nil {
			return uuid.Nil, false
		}
		return userID, true
	default:
		return uuid.Nil, false
	}
}

// TokenFromRequest extracts the bearer token from the Authorization header.
func TokenFromRequest(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", errMissingAuthorizationHeader
	}
	token, found := strings.CutPrefix(authHeader, "Bearer ")
	if !found || token == "" {
		return "", errInvalidAuthorizationHeader
	}
	return token, nil
}

// webSocketTicketFromRequest returns the ticket of a WebSocket upgrade that
// has no Authorization header. Browsers cannot set headers on a WebSocket, so
// the ticket travels in the URL; it works once and only for a few seconds.
func webSocketTicketFromRequest(r *http.Request) string {
	if !isWebSocketUpgradeRequest(r) || r.Header.Get("Authorization") != "" {
		return ""
	}
	return r.URL.Query().Get("ticket")
}

func isWebSocketUpgradeRequest(r *http.Request) bool {
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("Upgrade")), "websocket")
}
