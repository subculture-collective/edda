package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// DefaultTokenTTL is the default lifetime for issued JWT tokens (7 days).
const DefaultTokenTTL = 7 * 24 * time.Hour

// GenerateToken creates a signed HS256 JWT with the given user ID as the subject claim.
func GenerateToken(userID uuid.UUID, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		// Timestamps have one-second precision; the ID keeps each sign-in's token distinct.
		ID:        uuid.NewString(),
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken parses and validates the JWT, returning the user ID from the subject claim.
func ValidateToken(tokenString, secret string) (uuid.UUID, error) {
	userID, _, err := parseToken(tokenString, secret)
	return userID, err
}

// parseToken validates the JWT and returns its user ID and claims.
func parseToken(tokenString, secret string) (uuid.UUID, *jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return uuid.Nil, nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, claims, nil
}
