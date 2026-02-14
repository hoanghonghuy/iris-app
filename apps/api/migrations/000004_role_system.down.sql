-- Xóa bảng school_admins
DROP INDEX IF EXISTS idx_school_admins_user;
DROP INDEX IF EXISTS idx_school_admins_school;
DROP TABLE IF EXISTS school_admins;

-- Xóa role SCHOOL_ADMIN (xóa user_roles trước để tránh FK conflict)
DELETE FROM user_roles WHERE role_id IN (
  SELECT role_id FROM roles WHERE name = 'SCHOOL_ADMIN'
);
DELETE FROM roles WHERE name = 'SCHOOL_ADMIN';

-- Đổi tên SUPER_ADMIN → ADMIN (khôi phục lại tên cũ)
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_name_check;

UPDATE roles SET name = 'ADMIN' WHERE name = 'SUPER_ADMIN';

ALTER TABLE roles ADD CONSTRAINT roles_name_check
  CHECK (name IN ('ADMIN', 'TEACHER', 'PARENT'));
