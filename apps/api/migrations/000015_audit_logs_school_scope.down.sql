DROP INDEX IF EXISTS idx_audit_logs_school_id_created_at;

ALTER TABLE audit_logs
  DROP COLUMN IF EXISTS school_id;
