-- ============================================================================
-- IRIS Seed Data — MASTER RUNNER
-- Chạy tất cả các file seed theo thứ tự ưu tiên.
-- 
-- Cách chạy:
--   psql -U <user> -d <db> -f scripts/db/seed_01_schools_classes.sql
--   psql -U <user> -d <db> -f scripts/db/seed_02_users_profiles.sql
--   psql -U <user> -d <db> -f scripts/db/seed_03_attendance_health.sql
--   psql -U <user> -d <db> -f scripts/db/seed_04_posts_interactions.sql
--   psql -U <user> -d <db> -f scripts/db/seed_05_appointments_chat_audit.sql
--
-- Hoặc chạy file này (đã include tất cả):
--   psql -U <user> -d <db> -f scripts/db/seed_master.sql
-- ============================================================================

\echo '=== IRIS Seed Data — Master Runner ==='
\echo ''

\echo '[1/5] Schools & Classes...'
\i scripts/db/seed_01_schools_classes.sql

\echo '[2/5] Users & Profiles...'
\i scripts/db/seed_02_users_profiles.sql

\echo '[3/5] Attendance & Health...'
\i scripts/db/seed_03_attendance_health.sql

\echo '[4/5] Posts & Interactions...'
\i scripts/db/seed_04_posts_interactions.sql

\echo '[5/5] Appointments, Chat & Audit Logs...'
\i scripts/db/seed_05_appointments_chat_audit.sql

\echo ''
\echo '=== ALL DONE ==='

-- ============================================================================
-- SUMMARY OF EXPECTED ROW COUNTS (minimum):
--   schools:              8
--   classes:             56
--   users:               57 (1 super admin + 8 school admins + 24 teachers + 24 parents)
--   user_roles:          57
--   school_admins:        8
--   teachers:            24
--   parents:             24
--   students:            56
--   teacher_classes:     48 (24 teachers × 2 classes)
--   student_parents:     48 (24 parents × 2 students)
--   student_parent_codes: 56
--   attendance_records: 168 (56 students × 3 days)
--   health_logs:        112 (56 students × 2 logs)
--   attendance_change_logs: 50
--   posts:               56 (8 school + 24 class + 24 student)
--   post_attachments:    56
--   post_comments:       56
--   post_interactions:   76 (56 likes + 20 shares)
--   post_edit_history:   50
--   appointment_slots:   56
--   appointments:        56
--   conversations:       50 (30 direct + 20 group)
--   conversation_participants: ~130
--   messages:            56
--   audit_logs:          56
-- ============================================================================