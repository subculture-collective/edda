package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// RevocationQuerier stores the IDs (jti) of signed-out tokens until they
// expire. statedb.Querier satisfies it.
type RevocationQuerier interface {
	RevokeToken(ctx context.Context, arg statedb.RevokeTokenParams) error
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)
}

// revokeToken lists a valid token as signed out until it expires. It does
// nothing for a token that is invalid, already expired or has no ID.
func revokeToken(ctx context.Context, revocations RevocationQuerier, token, secret string) error {
	_, claims, err := parseToken(token, secret)
	if err != nil || claims.ID == "" || claims.ExpiresAt == nil {
		return nil
	}
	return revocations.RevokeToken(ctx, statedb.RevokeTokenParams{
		Jti:       claims.ID,
		ExpiresAt: pgtype.Timestamptz{Time: claims.ExpiresAt.Time, Valid: true},
	})
}
