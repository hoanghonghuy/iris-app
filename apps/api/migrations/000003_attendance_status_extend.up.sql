-- Migration: Extend attendance status to include late/excused
-- File: 000003_attendance_status_extend.up.sql

BEGIN;

-- Drop existing status constraint if present
ALTER TABLE attendance_records
  DROP CONSTRAINT IF EXISTS attendance_records_status_check;

-- Recreate constraint with extended status set
ALTER TABLE attendance_records
  ADD CONSTRAINT attendance_records_status_check
    CHECK (status IN ('present', 'absent', 'late', 'excused'));

-- Ensure default remains 'present'
ALTER TABLE attendance_records
  ALTER COLUMN status SET DEFAULT 'present';

COMMIT;
