-- +goose Up
-- revoked_tokens lists the jti of each signed-out token until the token expires.
-- The auth middleware refuses a token whose jti is listed here.
CREATE TABLE revoked_tokens (
  jti        text        PRIMARY KEY,
  expires_at timestamptz NOT NULL
);
CREATE INDEX revoked_tokens_expires_at_idx ON revoked_tokens (expires_at);

-- +goose Down
DROP TABLE revoked_tokens;
