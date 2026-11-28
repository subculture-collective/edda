-- name: RevokeToken :exec
-- Lists a token as revoked and drops rows for tokens that have expired, so the
-- table holds only tokens that would otherwise still be accepted.
WITH pruned AS (
  DELETE FROM revoked_tokens WHERE revoked_tokens.expires_at <= now()
)
INSERT INTO revoked_tokens (jti, expires_at)
VALUES (@jti, @expires_at)
ON CONFLICT (jti) DO NOTHING;

-- name: IsTokenRevoked :one
SELECT EXISTS (SELECT 1 FROM revoked_tokens WHERE jti = @jti)::boolean AS revoked;
