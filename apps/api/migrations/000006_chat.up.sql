-- Bảng cuộc hội thoại (nhóm hoặc trực tiếp)
CREATE TABLE IF NOT EXISTS conversations (
  conversation_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  type            varchar(10) NOT NULL DEFAULT 'direct', -- direct | group
  name            varchar(255),                          -- tên nhóm (NULL nếu direct)
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now()
);

-- Bảng thành viên tham gia cuộc hội thoại
CREATE TABLE IF NOT EXISTS conversation_participants (
  conversation_id uuid NOT NULL REFERENCES conversations(conversation_id) ON DELETE CASCADE,
  user_id         uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  joined_at       timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (conversation_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_conv_participants_user
  ON conversation_participants(user_id);

-- Bảng tin nhắn
CREATE TABLE IF NOT EXISTS messages (
  message_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id uuid NOT NULL REFERENCES conversations(conversation_id) ON DELETE CASCADE,
  sender_id       uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  content         text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation
  ON messages(conversation_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_messages_sender
  ON messages(sender_id);
