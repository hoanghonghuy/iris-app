ALTER TABLE audit_logs
  ADD COLUMN IF NOT EXISTS school_id uuid REFERENCES schools(school_id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_audit_logs_school_id_created_at
  ON audit_logs(school_id, created_at DESC);
