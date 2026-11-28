package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestTokenFromRequestReadsOnlyTheBearerToken(t *testing.T) {
	token, err := GenerateToken(uuid.New(), "secret", DefaultTokenTTL)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	upgrade := func(authorization string) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns/1/ws", nil)
		r.Header.Set("Upgrade", "websocket")
		// The cookie that once authenticated WebSocket upgrades is ignored.
		r.AddCookie(&http.Cookie{Name: "gm_token", Value: token})
		if authorization != "" {
			r.Header.Set("Authorization", authorization)
		}
		return r
	}
	for _, tc := range []struct {
		name    string
		req     *http.Request
		wantTok string
	}{
		{name: "cookie only", req: upgrade("")},
		{name: "malformed", req: upgrade("Bearer")},
		{name: "not bearer", req: upgrade("Basic abc")},
		{name: "bearer", req: upgrade("Bearer " + token), wantTok: token},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := TokenFromRequest(tc.req)
			if tc.wantTok == "" {
				if err == nil {
					t.Fatalf("expected error, got token %q", got)
				}
				return
			}
			if err != nil || got != tc.wantTok {
				t.Fatalf("token = %q, %v; want %q", got, err, tc.wantTok)
			}
		})
	}
}

func TestUserFromContext(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		if userID, ok := UserFromContext(nil); ok || userID != uuid.Nil {
			t.Fatalf("expected missing user, got %s %v", userID, ok)
		}
	})

	t.Run("nil", func(t *testing.T) {
		if userID, ok := UserFromContext(WithUser(context.Background(), uuid.Nil)); ok || userID != uuid.Nil {
			t.Fatalf("expected nil user, got %s %v", userID, ok)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		ctx := contextWithValueForTest("not-a-uuid")
		if userID, ok := UserFromContext(ctx); ok || userID != uuid.Nil {
			t.Fatalf("expected wrong type to be rejected, got %s %v", userID, ok)
		}
	})

	t.Run("valid", func(t *testing.T) {
		userID := uuid.New()
		if got, ok := UserFromContext(WithUser(context.Background(), userID)); !ok || got != userID {
			t.Fatalf("expected %s, got %s %v", userID, got, ok)
		}
	})
}

func contextWithValueForTest(v any) context.Context {
	return context.WithValue(context.Background(), sessionContextKey{}, v)
}
