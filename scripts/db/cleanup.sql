-- Cleanup all data before fresh seed
-- Truncates all tables in correct order respecting foreign key dependencies
-- Keeps reference data (roles) untouched

-- Start transaction for safety
BEGIN;

-- Disable triggers temporarily to avoid cascading issues
SET session_replication_role = replica;

-- Truncate in order: leaf tables first, then parent tables
TRUNCATE TABLE
  attendance_records,
  health_logs,
  post_attachments,
  posts,
  student_parent_codes,
  student_parents,
  teacher_classes,
  students,
  parents,
  teachers,
  school_admins,
  user_roles,
  users,
  classes,
  schools
CASCADE;

-- Re-enable triggers
SET session_replication_role = DEFAULT;

-- Reset sequences to start from 1 (optional, useful for testing)
-- ALTER SEQUENCE users_user_id_seq RESTART WITH 1;
-- ALTER SEQUENCE schools_school_id_seq RESTART WITH 1;
-- ALTER SEQUENCE classes_class_id_seq RESTART WITH 1;
-- ALTER SEQUENCE students_student_id_seq RESTART WITH 1;
-- ALTER SEQUENCE teachers_teacher_id_seq RESTART WITH 1;
-- ALTER SEQUENCE parents_parent_id_seq RESTART WITH 1;
-- ALTER SEQUENCE school_admins_admin_id_seq RESTART WITH 1;
-- ALTER SEQUENCE roles_role_id_seq RESTART WITH 1;
-- ALTER SEQUENCE attendance_records_attendance_id_seq RESTART WITH 1;
-- ALTER SEQUENCE health_logs_health_log_id_seq RESTART WITH 1;
-- ALTER SEQUENCE posts_post_id_seq RESTART WITH 1;
-- ALTER SEQUENCE post_attachments_attachment_id_seq RESTART WITH 1;
-- ALTER SEQUENCE student_parent_codes_code_id_seq RESTART WITH 1;

-- Note: roles table is NOT truncated as it contains reference data
-- (SUPER_ADMIN, SCHOOL_ADMIN, TEACHER, PARENT)

COMMIT;

-- Verify cleanup
SELECT 'Tables cleaned successfully. Ready for fresh seed.' as status;
