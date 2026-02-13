-- Activation tokens for teachers
ALTER TABLE users
ADD COLUMN activation_token varchar(255),
ADD COLUMN token_expires_at timestamptz;

CREATE INDEX IF NOT EXISTS idx_users_activation_token
  ON users(activation_token) WHERE activation_token IS NOT NULL;


-- Parent codes
CREATE TABLE IF NOT EXISTS student_parent_codes (
  code_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  student_id   uuid NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
  code         varchar(8) UNIQUE NOT NULL
    -- Validation regex:
        -- ^: bắt đầu chuỗi
        -- [A-Z0-9]: chỉ chấp nhận ký tự hoa hoặc số
        -- {6,8}: độ dài từ 6 đến 8 ký tự
        -- $: kết thúc chuỗi
    CHECK (code ~ '^[A-Z0-9]{6,8}$'),
  usage_count  int NOT NULL DEFAULT 0
    CHECK (usage_count >= 0),
  max_usage    int NOT NULL DEFAULT 4
    CHECK (max_usage > 0 AND max_usage <= 10),
  created_at   timestamptz NOT NULL DEFAULT now(),
  expires_at   timestamptz NOT NULL DEFAULT (now() + interval '1 year'),
  CHECK (usage_count <= max_usage)
);

CREATE INDEX IF NOT EXISTS idx_parent_codes_student
  ON student_parent_codes(student_id);

CREATE INDEX IF NOT EXISTS idx_parent_codes_code
  ON student_parent_codes(code);
