-- Theo dõi thời điểm user xem hội thoại (unread + sort sidebar)
ALTER TABLE conversation_participants
  ADD COLUMN IF NOT EXISTS last_read_at timestamptz;
