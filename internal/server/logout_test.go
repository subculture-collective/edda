package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
)

func TestLogoutClearsOnlyTheSessionBeingLoggedOut(t *testing.T) {
	h, _, _, user, _, _ := ownershipFixture(t)
	// Two sign-ins by one user, usually within the same second.
	a, b := ownershipToken(t, user), ownershipToken(t, user)
	other := ownershipToken(t, uuid.New())

	for _, tc := range []struct {
		name, authorization, cookie string
		wantCleared                 bool
	}{
		{"own cookie", "Bearer " + a, a, true},
		{"cookie of a newer sign-in", "Bearer " + a, b, false},
		{"cookie of another user", "Bearer " + a, other, false},
		{"no cookie sent", "Bearer " + a, "", false},
		{"no bearer token", "", a, true},
		{"unusable authorization header", "Basic abc", a, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
			if tc.authorization != "" {
				req.Header.Set("Authorization", tc.authorization)
			}
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: auth.AuthCookieName, Value: tc.cookie})
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
			}
			cookies := rec.Result().Cookies()
			if !tc.wantCleared {
				if len(cookies) != 0 {
					t.Fatalf("logout must leave the cookie alone, got %v", cookies)
				}
				return
			}
			if len(cookies) != 1 || cookies[0].Name != auth.AuthCookieName || cookies[0].MaxAge != -1 {
				t.Fatalf("logout must clear the cookie, got %v", cookies)
			}
		})
	}
}
