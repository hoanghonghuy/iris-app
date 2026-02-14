-- Đổi tên role ADMIN → SUPER_ADMIN, thêm SCHOOL_ADMIN
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_name_check;

UPDATE roles SET name = 'SUPER_ADMIN' WHERE name = 'ADMIN';

-- Thêm ràng buộc CHECK mới với tên vai trò được cập nhật
ALTER TABLE roles ADD CONSTRAINT roles_name_check
  CHECK (name IN ('SUPER_ADMIN', 'SCHOOL_ADMIN', 'TEACHER', 'PARENT'));

INSERT INTO roles (name) VALUES ('SCHOOL_ADMIN')
ON CONFLICT (name) DO NOTHING;


-- Bảng school_admins (cùng pattern với teachers/parents — liên kết user với trường)
CREATE TABLE IF NOT EXISTS school_admins (
  admin_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    -- UNIQUE: 1 user chỉ admin 1 trường (đổi trường thì mất quyền admin)
  school_id  uuid NOT NULL REFERENCES schools(school_id) ON DELETE CASCADE,
  full_name  varchar(255) NOT NULL,
  phone      varchar(30),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_school_admins_school
  ON school_admins(school_id);

CREATE INDEX IF NOT EXISTS idx_school_admins_user
  ON school_admins(user_id);
