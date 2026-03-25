-- Audit trail cho thay đổi điểm danh
CREATE TABLE IF NOT EXISTS attendance_change_logs (
  change_id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  attendance_id   uuid NOT NULL REFERENCES attendance_records(attendance_id) ON DELETE CASCADE,
  student_id      uuid NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
  date            date NOT NULL,
  change_type     varchar(10) NOT NULL CHECK (change_type IN ('create', 'update')),
  old_status      varchar(20) CHECK (old_status IN ('present','absent','late','excused')),
  new_status      varchar(20) NOT NULL CHECK (new_status IN ('present','absent','late','excused')),
  old_note        text,
  new_note        text,
  changed_by      uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  changed_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_attendance_change_logs_student_date
  ON attendance_change_logs (student_id, date, changed_at DESC);

CREATE INDEX IF NOT EXISTS idx_attendance_change_logs_attendance
  ON attendance_change_logs (attendance_id, changed_at DESC);
