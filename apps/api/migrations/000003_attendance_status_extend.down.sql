-- Migration: Revert attendance status extension (late/excused) back to present/absent
-- File: 000003_attendance_status_extend.down.sql

BEGIN;

-- Drop the extended constraint
ALTER TABLE attendance_records
  DROP CONSTRAINT IF EXISTS attendance_records_status_check;

-- Restore original constraint (present/absent)
ALTER TABLE attendance_records
  ADD CONSTRAINT attendance_records_status_check
    CHECK (status IN ('present', 'absent'));

-- Ensure default remains 'present'
ALTER TABLE attendance_records
  ALTER COLUMN status SET DEFAULT 'present';

COMMIT;
