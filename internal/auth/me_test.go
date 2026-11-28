package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type meQuerier struct{ AuthQuerier }

func (meQuerier) GetUserByID(_ context.Context, id pgtype.UUID) (UserWithAuth, error) {
	return UserWithAuth{ID: id, Name: "Player"}, nil
}

func TestMeSetsNoCookie(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID, "secret", DefaultTokenTTL)
	if err != nil {
		t.Fatal(err)
	}
	h := NewAuthHandlers(meQuerier{}, newMemRevocations(), NewWebSocketSessions(), "secret")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	NewJWTMiddleware("secret", newMemRevocations(), NewWebSocketSessions()).Authenticate(http.HandlerFunc(h.Me)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", rec.Code, rec.Body.String())
	}
	if cookies := rec.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("expected no cookie, got %v", cookies)
	}
}
