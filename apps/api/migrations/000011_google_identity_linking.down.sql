DROP INDEX IF EXISTS idx_users_google_sub_unique;

ALTER TABLE users
DROP COLUMN IF EXISTS google_linked_at,
DROP COLUMN IF EXISTS google_sub;
