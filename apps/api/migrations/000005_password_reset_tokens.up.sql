-- Password reset tokens (separate table for security & audit)
CREATE TABLE IF NOT EXISTS password_reset_tokens (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  token_hash varchar(64) NOT NULL,       -- SHA-256 hex digest (64 chars)
  expires_at timestamptz NOT NULL,
  used_at    timestamptz,                -- NULL = not yet used
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reset_tokens_hash
  ON password_reset_tokens(token_hash) WHERE used_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_reset_tokens_user
  ON password_reset_tokens(user_id);
