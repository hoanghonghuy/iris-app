-- Lưu lịch sử chỉnh sửa nội dung bài đăng
CREATE TABLE IF NOT EXISTS post_edit_history (
  edit_id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id       uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
  old_content   text NOT NULL,
  new_content   text NOT NULL,
  edited_by     uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  edited_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_post_edit_history_post_edited_at
  ON post_edit_history (post_id, edited_at DESC);
