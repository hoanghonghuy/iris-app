ALTER TABLE users
ADD COLUMN IF NOT EXISTS google_sub varchar(255),
ADD COLUMN IF NOT EXISTS google_linked_at timestamptz;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_sub_unique
  ON users(google_sub)
  WHERE google_sub IS NOT NULL;
