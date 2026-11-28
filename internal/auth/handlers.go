package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// AuthQuerier defines the database methods required by auth handlers.
// This is satisfied by the generated statedb.Queries after running sqlc generate.
type AuthQuerier interface {
	CreateUserWithAuth(ctx context.Context, arg CreateUserWithAuthParams) (UserWithAuth, error)
	GetUserByEmail(ctx context.Context, email string) (UserWithAuth, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (UserWithAuth, error)
}

// CreateUserWithAuthParams holds the parameters for creating an authenticated user.
type CreateUserWithAuthParams struct {
	Name         string
	Email        string
	PasswordHash string
}

// UserWithAuth represents a user row including auth fields.
type UserWithAuth struct {
	ID           pgtype.UUID
	Name         string
	Email        pgtype.Text
	PasswordHash pgtype.Text
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

// AuthHandlers provides HTTP handlers for authentication.
type AuthHandlers struct {
	queries     AuthQuerier
	revocations RevocationQuerier
	jwtSecret   string
}

// NewAuthHandlers creates a new AuthHandlers. Logout lists signed-out tokens
// in revocations.
func NewAuthHandlers(queries AuthQuerier, revocations RevocationQuerier, jwtSecret string) *AuthHandlers {
	return &AuthHandlers{queries: queries, revocations: revocations, jwtSecret: jwtSecret}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string   `json:"token"`
	User  userJSON `json:"user"`
}

type userJSON struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "name, email, and password are required")
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.queries.CreateUserWithAuth(r.Context(), CreateUserWithAuthParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
	})
	if err != nil {
		// Check for unique constraint violation on email.
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	uid := uuidFromPgtype(user.ID)
	token, err := GenerateToken(uid, h.jwtSecret, DefaultTokenTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	IssueSessionCookie(w, r, token)

	writeJSON(w, http.StatusCreated, authResponse{
		Token: token,
		User: userJSON{
			ID:    uid.String(),
			Name:  user.Name,
			Email: user.Email.String,
		},
	})
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to query user")
		return
	}

	if !user.PasswordHash.Valid || !CheckPassword(req.Password, user.PasswordHash.String) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	uid := uuidFromPgtype(user.ID)
	token, err := GenerateToken(uid, h.jwtSecret, DefaultTokenTTL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	IssueSessionCookie(w, r, token)

	writeJSON(w, http.StatusOK, authResponse{
		Token: token,
		User: userJSON{
			ID:    uid.String(),
			Name:  user.Name,
			Email: user.Email.String,
		},
	})
}

// Logout handles POST /api/v1/auth/logout. It revokes the bearer token, or the
// session cookie's token when no bearer token is sent, so the token is refused
// from then on. It clears the session cookie unless the request shows that the
// cookie belongs to a session other than the bearer token's, so logging out one
// session leaves another one's cookie alone.
//
// Only the cookie sent with this request is visible here. If a client signs in
// while this reply is in flight, the reply can still clear the new cookie; the
// client then calls Me to set it again.
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	cookieToken := ""
	if cookie, err := r.Cookie(AuthCookieName); err == nil {
		cookieToken = strings.TrimSpace(cookie.Value)
	}
	token, err := TokenFromRequest(r)
	clearCookie := err != nil || cookieToken == token
	if err != nil {
		token = cookieToken
	}
	if token != "" {
		if err := revokeToken(r.Context(), h.revocations, token, h.jwtSecret); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to revoke session")
			return
		}
	}
	if clearCookie {
		ClearSessionCookie(w, r)
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// Me handles GET /api/v1/auth/me — returns the current authenticated user and
// sets the session cookie to the token that authenticated the request.
func (h *AuthHandlers) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var pgID pgtype.UUID
	pgID.Bytes = userID
	pgID.Valid = true

	user, err := h.queries.GetUserByID(r.Context(), pgID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	if token, err := TokenFromRequest(r); err == nil {
		IssueSessionCookie(w, r, token)
	}

	writeJSON(w, http.StatusOK, authResponse{
		User: userJSON{
			ID:    uuidFromPgtype(user.ID).String(),
			Name:  user.Name,
			Email: user.Email.String,
		},
	})
}

func uuidFromPgtype(id pgtype.UUID) uuid.UUID {
	return uuid.UUID(id.Bytes)
}
