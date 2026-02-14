DROP INDEX IF EXISTS idx_parent_codes_code;
DROP INDEX IF EXISTS idx_parent_codes_student;
DROP TABLE IF EXISTS student_parent_codes;

DROP INDEX IF EXISTS idx_users_activation_token;
ALTER TABLE users
DROP COLUMN IF EXISTS token_expires_at,
DROP COLUMN IF EXISTS activation_token;
