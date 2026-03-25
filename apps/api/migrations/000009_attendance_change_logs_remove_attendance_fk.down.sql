ALTER TABLE attendance_change_logs
ADD CONSTRAINT attendance_change_logs_attendance_id_fkey
FOREIGN KEY (attendance_id)
REFERENCES attendance_records(attendance_id)
ON DELETE CASCADE;
