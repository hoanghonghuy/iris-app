-- Clean all seed data
-- Removes application data in dependency order and keeps reference data intact.

BEGIN;

DELETE FROM audit_logs;
DELETE FROM messages;
DELETE FROM conversations_participants;
DELETE FROM conversations;
DELETE FROM appointment_slots;
DELETE FROM appointments;
DELETE FROM post_interactions;
DELETE FROM comments;
DELETE FROM posts;
DELETE FROM google_links;
DELETE FROM password_reset_tokens;
DELETE FROM attendance_change_logs;
DELETE FROM attendance;
DELETE FROM health_logs;
DELETE FROM student_parents;
DELETE FROM user_roles;
DELETE FROM teachers;
DELETE FROM parents;
DELETE FROM school_admins;
DELETE FROM students;
DELETE FROM classes;
DELETE FROM schools;
DELETE FROM users;

COMMIT;

SELECT COUNT(*) AS remaining_records FROM users;
