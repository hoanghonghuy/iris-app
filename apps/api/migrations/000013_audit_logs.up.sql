CREATE TABLE IF NOT EXISTS audit_logs (
  audit_log_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_user_id  uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  actor_role     varchar(30),
  action         varchar(120) NOT NULL,
  entity_type    varchar(80) NOT NULL,
  entity_id      uuid,
  details        jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at
  ON audit_logs(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_actor
  ON audit_logs(actor_user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_action
  ON audit_logs(action, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type
  ON audit_logs(entity_type, created_at DESC);
