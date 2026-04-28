-- ============================================================================
-- IRIS Seed Data — File 2: Users & Profiles
-- Mục tiêu: 50+ users, teachers, parents, students, school_admins
--           + teacher_classes, student_parents, student_parent_codes
-- Thứ tự chạy: 2/5 (cần chạy sau seed_01_schools_classes.sql)
-- ============================================================================

DO $$
DECLARE
  -- School IDs (lấy từ DB)
  v_school_hanoi_1    uuid;
  v_school_hanoi_2    uuid;
  v_school_hcm_1      uuid;
  v_school_hcm_2      uuid;
  v_school_danang_1   uuid;
  v_school_danang_2   uuid;
  v_school_haiphong   uuid;
  v_school_cantho     uuid;

  -- Class IDs (lấy theo tên)
  v_class_ids         uuid[];

  -- User IDs
  v_super_admin_user  uuid;
  v_user_ids          uuid[] := ARRAY[]::uuid[];

  -- Profile IDs
  v_teacher_ids       uuid[] := ARRAY[]::uuid[];
  v_parent_ids        uuid[] := ARRAY[]::uuid[];
  v_admin_ids         uuid[] := ARRAY[]::uuid[];
  v_student_ids       uuid[] := ARRAY[]::uuid[];

  -- Temp
  v_uid               uuid;
  v_pid               uuid;
  v_tid               uuid;
  v_aid               uuid;
  v_sid               uuid;
  v_cid               uuid;
  i                   int;
  j                   int;
BEGIN
  -- ======================
  -- 0) Lấy school IDs
  -- ======================
  SELECT school_id INTO v_school_hanoi_1  FROM schools WHERE name = 'Trường Mầm non Bình Minh Cầu Giấy' LIMIT 1;
  SELECT school_id INTO v_school_hanoi_2  FROM schools WHERE name = 'Trường Mầm non Tuổi Thần Tiên Đống Đa' LIMIT 1;
  SELECT school_id INTO v_school_hcm_1    FROM schools WHERE name = 'Trường Mầm non Hoa Mai Quận 7' LIMIT 1;
  SELECT school_id INTO v_school_hcm_2    FROM schools WHERE name = 'Trường Mầm non Vườn Xanh Thủ Đức' LIMIT 1;
  SELECT school_id INTO v_school_danang_1 FROM schools WHERE name = 'Trường Mầm non Ánh Dương Hải Châu' LIMIT 1;
  SELECT school_id INTO v_school_danang_2 FROM schools WHERE name = 'Trường Mầm non Ngọc Lan Sơn Trà' LIMIT 1;
  SELECT school_id INTO v_school_haiphong FROM schools WHERE name = 'Trường Mầm non Hướng Dương Hải Phòng' LIMIT 1;
  SELECT school_id INTO v_school_cantho   FROM schools WHERE name = 'Trường Mầm non Cánh Én Cần Thơ' LIMIT 1;

  -- ======================
  -- 1) SUPER ADMIN (1 user)
  -- ======================
  INSERT INTO users (email, password_hash, status)
  VALUES ('admin@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_super_admin_user;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_super_admin_user, role_id FROM roles WHERE name = 'SUPER_ADMIN'
  ON CONFLICT DO NOTHING;

  -- ======================
  -- 2) SCHOOL ADMINS (8 users — mỗi trường 1 admin)
  -- ======================
  INSERT INTO users (email, password_hash, status)
  VALUES
    ('admin-hn1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-hn2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-hp@iris.local',  '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-dn1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-dn2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-hcm1@iris.local','$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-hcm2@iris.local','$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active'),
    ('admin-ct@iris.local',  '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO NOTHING;

  -- Gán role SCHOOL_ADMIN + tạo profile
  FOR i IN 1..8 LOOP
    CASE i
      WHEN 1 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-hn1@iris.local';
      WHEN 2 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-hn2@iris.local';
      WHEN 3 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-hp@iris.local';
      WHEN 4 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-dn1@iris.local';
      WHEN 5 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-dn2@iris.local';
      WHEN 6 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-hcm1@iris.local';
      WHEN 7 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-hcm2@iris.local';
      WHEN 8 THEN SELECT user_id INTO v_uid FROM users WHERE email = 'admin-ct@iris.local';
    END CASE;

    INSERT INTO user_roles (user_id, role_id)
    SELECT v_uid, role_id FROM roles WHERE name = 'SCHOOL_ADMIN'
    ON CONFLICT DO NOTHING;

    CASE i
      WHEN 1 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_hanoi_1, 'Bùi Thanh Hương', '0972331184')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 2 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_hanoi_2, 'Nguyễn Thị Mai Lan', '0987654321')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 3 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_haiphong, 'Trần Văn Hùng', '0912345678')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 4 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_danang_1, 'Lê Thị Hồng Nhung', '0905123456')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 5 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_danang_2, 'Phạm Quốc Bảo', '0938111222')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 6 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_hcm_1, 'Võ Thị Kim Chi', '0967555444')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 7 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_hcm_2, 'Đặng Minh Tuấn', '0979888777')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
      WHEN 8 THEN
        INSERT INTO school_admins (user_id, school_id, full_name, phone)
        VALUES (v_uid, v_school_cantho, 'Huỳnh Ngọc Trâm', '0918777666')
        ON CONFLICT (user_id) DO NOTHING
        RETURNING admin_id INTO v_aid;
    END CASE;

    v_admin_ids := array_append(v_admin_ids, v_aid);
  END LOOP;

  -- ======================
  -- 3) TEACHERS (24 users — mỗi trường 3 giáo viên)
  -- ======================
  -- Hanoi 1
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hn1-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hn1-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hn1-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- Hanoi 2
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hn2-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hn2-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hn2-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- Hai Phong
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hp-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hp-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hp-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- Da Nang 1
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-dn1-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-dn1-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-dn1-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- Da Nang 2
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-dn2-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-dn2-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-dn2-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- HCM 1
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hcm1-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hcm1-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hcm1-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- HCM 2
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hcm2-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hcm2-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-hcm2-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  -- Can Tho
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-ct-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-ct-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('teacher-ct-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;

  -- Teacher profiles + role
  DECLARE
    teacher_emails text[] := ARRAY[
      'teacher-hn1-1@iris.local','teacher-hn1-2@iris.local','teacher-hn1-3@iris.local',
      'teacher-hn2-1@iris.local','teacher-hn2-2@iris.local','teacher-hn2-3@iris.local',
      'teacher-hp-1@iris.local','teacher-hp-2@iris.local','teacher-hp-3@iris.local',
      'teacher-dn1-1@iris.local','teacher-dn1-2@iris.local','teacher-dn1-3@iris.local',
      'teacher-dn2-1@iris.local','teacher-dn2-2@iris.local','teacher-dn2-3@iris.local',
      'teacher-hcm1-1@iris.local','teacher-hcm1-2@iris.local','teacher-hcm1-3@iris.local',
      'teacher-hcm2-1@iris.local','teacher-hcm2-2@iris.local','teacher-hcm2-3@iris.local',
      'teacher-ct-1@iris.local','teacher-ct-2@iris.local','teacher-ct-3@iris.local'
    ];
    teacher_names text[] := ARRAY[
      'Nguyễn Thị Lan Anh','Trần Hoàng Nam','Phạm Thu Thủy',
      'Lê Văn Đức','Đỗ Thị Hồng','Ngô Minh Phương',
      'Vũ Thị Hải Yến','Bùi Quang Huy','Đinh Thị Ngọc',
      'Đặng Minh Phương','Lý Thị Thanh','Hoàng Văn Tùng',
      'Trương Thị Mỹ Linh','Nguyễn Văn Thành','Phan Thị Kim Ngân',
      'Võ Thị Cẩm Tú','Lâm Minh Nhật','Trịnh Thị Bích',
      'Đoàn Văn Khánh','Nguyễn Thị Diễm','Mai Thanh Sơn',
      'Huỳnh Thị Thảo','Lê Quốc Việt','Dương Thị Mỹ Duyên'
    ];
    teacher_phones text[] := ARRAY[
      '0983124501','0917642208','0978555123',
      '0967123456','0908123987','0934555789',
      '0912345601','0978234502','0909456703',
      '0908765123','0932123456','0977456789',
      '0918234567','0965345678','0909876543',
      '0987654321','0912123987','0938456712',
      '0979345612','0901234567','0966789012',
      '0915456789','0932678901','0988012345'
    ];
    teacher_schools uuid[] := ARRAY[
      v_school_hanoi_1,v_school_hanoi_1,v_school_hanoi_1,
      v_school_hanoi_2,v_school_hanoi_2,v_school_hanoi_2,
      v_school_haiphong,v_school_haiphong,v_school_haiphong,
      v_school_danang_1,v_school_danang_1,v_school_danang_1,
      v_school_danang_2,v_school_danang_2,v_school_danang_2,
      v_school_hcm_1,v_school_hcm_1,v_school_hcm_1,
      v_school_hcm_2,v_school_hcm_2,v_school_hcm_2,
      v_school_cantho,v_school_cantho,v_school_cantho
    ];
  BEGIN
    FOR i IN 1..24 LOOP
      SELECT user_id INTO v_uid FROM users WHERE email = teacher_emails[i];

      INSERT INTO user_roles (user_id, role_id)
      SELECT v_uid, role_id FROM roles WHERE name = 'TEACHER'
      ON CONFLICT DO NOTHING;

      INSERT INTO teachers (user_id, school_id, full_name, phone)
      VALUES (v_uid, teacher_schools[i], teacher_names[i], teacher_phones[i])
      ON CONFLICT (user_id) DO NOTHING
      RETURNING teacher_id INTO v_tid;

      v_teacher_ids := array_append(v_teacher_ids, v_tid);
    END LOOP;
  END;

  -- ======================
  -- 4) PARENTS (24 users — mỗi trường 3 phụ huynh)
  -- ======================
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hn1-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hn1-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hn1-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hn2-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hn2-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hn2-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hp-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hp-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hp-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-dn1-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-dn1-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-dn1-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-dn2-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-dn2-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-dn2-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hcm1-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hcm1-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hcm1-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hcm2-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hcm2-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-hcm2-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-ct-1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-ct-2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;
  INSERT INTO users (email, password_hash, status) VALUES ('parent-ct-3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active') ON CONFLICT (email) DO NOTHING;

  -- Parent profiles + role
  DECLARE
    parent_emails text[] := ARRAY[
      'parent-hn1-1@iris.local','parent-hn1-2@iris.local','parent-hn1-3@iris.local',
      'parent-hn2-1@iris.local','parent-hn2-2@iris.local','parent-hn2-3@iris.local',
      'parent-hp-1@iris.local','parent-hp-2@iris.local','parent-hp-3@iris.local',
      'parent-dn1-1@iris.local','parent-dn1-2@iris.local','parent-dn1-3@iris.local',
      'parent-dn2-1@iris.local','parent-dn2-2@iris.local','parent-dn2-3@iris.local',
      'parent-hcm1-1@iris.local','parent-hcm1-2@iris.local','parent-hcm1-3@iris.local',
      'parent-hcm2-1@iris.local','parent-hcm2-2@iris.local','parent-hcm2-3@iris.local',
      'parent-ct-1@iris.local','parent-ct-2@iris.local','parent-ct-3@iris.local'
    ];
    parent_names text[] := ARRAY[
      'Lê Quang Minh','Phạm Thu Hoa','Nguyễn Văn Hải',
      'Trần Thị Thanh','Đỗ Đức Anh','Ngô Thị Bích Ngọc',
      'Vũ Văn Cường','Bùi Thị Hạnh','Đinh Công Thành',
      'Lý Văn Phúc','Hoàng Thị Xuân','Nguyễn Đức Tài',
      'Trương Minh Quân','Phan Thị Ngọc Ánh','Lê Văn Lợi',
      'Võ Văn Quân','Lâm Thị Kim Oanh','Trịnh Hoàng Long',
      'Đoàn Thị Mỹ Hạnh','Nguyễn Thanh Bình','Mai Văn Phát',
      'Huỳnh Văn Tèo','Lê Thị Cẩm Vân','Dương Quốc Huy'
    ];
    parent_phones text[] := ARRAY[
      '0935587612','0964412097','0978123456',
      '0912345098','0908123765','0934555432',
      '0987123401','0978234102','0909456303',
      '0918234104','0932123455','0977456306',
      '0965345107','0909876208','0987654309',
      '0912123410','0938456311','0979345612',
      '0901234313','0966789014','0915456315',
      '0932678416','0988012317','0909123418'
    ];
    parent_schools uuid[] := ARRAY[
      v_school_hanoi_1,v_school_hanoi_1,v_school_hanoi_1,
      v_school_hanoi_2,v_school_hanoi_2,v_school_hanoi_2,
      v_school_haiphong,v_school_haiphong,v_school_haiphong,
      v_school_danang_1,v_school_danang_1,v_school_danang_1,
      v_school_danang_2,v_school_danang_2,v_school_danang_2,
      v_school_hcm_1,v_school_hcm_1,v_school_hcm_1,
      v_school_hcm_2,v_school_hcm_2,v_school_hcm_2,
      v_school_cantho,v_school_cantho,v_school_cantho
    ];
  BEGIN
    FOR i IN 1..24 LOOP
      SELECT user_id INTO v_uid FROM users WHERE email = parent_emails[i];

      INSERT INTO user_roles (user_id, role_id)
      SELECT v_uid, role_id FROM roles WHERE name = 'PARENT'
      ON CONFLICT DO NOTHING;

      INSERT INTO parents (user_id, school_id, full_name, phone)
      VALUES (v_uid, parent_schools[i], parent_names[i], parent_phones[i])
      ON CONFLICT (user_id) DO NOTHING
      RETURNING parent_id INTO v_pid;

      v_parent_ids := array_append(v_parent_ids, v_pid);
    END LOOP;
  END;

  -- ======================
  -- 5) STUDENTS (56 students — mỗi lớp 1 học sinh, phân bổ đều)
  -- ======================
  -- Lấy tất cả class_id theo school
  DECLARE
    student_names text[] := ARRAY[
      'Nguyễn Gia An','Trần Khánh Băng','Lê Ngọc Nhi','Phạm Đức Dương','Đỗ Minh Khang','Ngô Phương My','Vũ Bảo Long',
      'Bùi Thanh Tùng','Đinh Gia Trâm','Hoàng Minh Đức','Lý Bảo Châu','Nguyễn Hải Đăng','Phan Thanh Hằng','Trương Quốc Bảo',
      'Đoàn Mỹ Linh','Mai Thanh Tú','Dương Hồng Nhung','Huỳnh Bảo Khôi','Lâm Gia Huy','Trịnh Ngọc Hân','Võ Minh Triết',
      'Đặng Hoàng Phúc','Lê Thị Cẩm Tiên','Ngô Quang Vinh','Phạm Ngọc Trân','Trần Thanh Phong','Vũ Thị Kim Ngân',
      'Bùi Đức Thịnh','Đỗ Hồng Sơn','Hoàng Thị Bích Liên','Lý Minh Bảo','Nguyễn Thanh Trúc','Phan Văn Kiệt','Trương Hạ Vy',
      'Đoàn Quốc Đạt','Mai Thị Hồng Đào','Dương Minh Nhật','Huỳnh Thị Ngọc Mai','Lâm Văn Tài','Trịnh Công Minh','Võ Thị Thùy Dương',
      'Đặng Gia Hân','Lê Quang Huy','Ngô Thị Bích Phượng','Phạm Minh Khôi','Trần Ngọc Bích','Vũ Đức Anh',
      'Bùi Thị Thanh Trúc','Đỗ Minh Quân','Hoàng Gia Bảo','Lý Thị Kim Chi','Nguyễn Đức Thắng','Phan Thị Mỹ Dung','Trương Văn Nam',
      'Đoàn Thị Hồng Nhung','Mai Quốc Cường'
    ];
    student_dobs date[] := ARRAY[
      '2020-05-12','2020-09-03','2021-01-20','2020-11-15','2020-07-08','2021-02-14','2020-03-19',
      '2020-06-22','2020-12-02','2021-03-10','2020-08-17','2020-04-25','2021-05-30','2020-10-05',
      '2020-02-14','2021-06-18','2020-09-28','2020-03-19','2021-01-07','2020-07-22','2020-11-30',
      '2020-04-15','2021-02-28','2020-08-09','2020-12-24','2020-06-11','2021-04-03',
      '2020-10-18','2020-03-07','2021-07-12','2021-04-11','2020-09-21','2020-05-06','2020-10-27',
      '2020-01-30','2021-08-14','2020-06-25','2020-11-08','2021-03-22','2020-07-15','2020-12-09',
      '2020-04-28','2021-01-13','2020-08-31','2020-05-19','2021-06-02','2020-02-07',
      '2020-09-14','2021-04-26','2020-07-29','2020-12-16','2020-03-25','2021-05-08','2020-10-11',
      '2020-06-04','2021-02-20'
    ];
    student_genders text[] := ARRAY[
      'male','female','female','male','male','female','male',
      'male','female','male','female','male','female','male',
      'female','male','female','male','male','female','male',
      'male','female','male','female','male','female',
      'male','male','female','male','female','male','female',
      'male','female','male','female','male','female','male',
      'female','male','female','male','female','male',
      'female','male','female','male','female','male','female',
      'female','male'
    ];
    class_records RECORD;
    class_idx int := 0;
  BEGIN
    FOR class_records IN
      SELECT c.class_id, c.school_id
      FROM classes c
      ORDER BY c.school_id, c.name
    LOOP
      class_idx := class_idx + 1;
      IF class_idx > 56 THEN EXIT; END IF;

      INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
      VALUES (class_records.school_id, class_records.class_id,
              student_names[class_idx], student_dobs[class_idx], student_genders[class_idx])
      ON CONFLICT DO NOTHING
      RETURNING student_id INTO v_sid;

      v_student_ids := array_append(v_student_ids, v_sid);
    END LOOP;
  END;

  -- ======================
  -- 6) TEACHER-CLASS ASSIGNMENTS (mỗi giáo viên dạy 2-3 lớp)
  -- ======================
  IF to_regclass('public.teacher_classes') IS NOT NULL THEN
    -- Mỗi teacher dạy 2 lớp trong trường của họ
    FOR i IN 1..24 LOOP
      -- Lấy 2 class_id trong cùng school của teacher
      FOR class_records IN
        SELECT c.class_id
        FROM classes c
        JOIN teachers t ON t.teacher_id = v_teacher_ids[i] AND c.school_id = t.school_id
        ORDER BY c.name
        LIMIT 2
      LOOP
        INSERT INTO teacher_classes (teacher_id, class_id)
        VALUES (v_teacher_ids[i], class_records.class_id)
        ON CONFLICT DO NOTHING;
      END LOOP;
    END LOOP;
  END IF;

  -- ======================
  -- 7) STUDENT-PARENT ASSIGNMENTS (mỗi học sinh 1-2 phụ huynh)
  -- ======================
  IF to_regclass('public.student_parents') IS NOT NULL THEN
    -- Phân bổ: mỗi parent trong school được gán 2-3 students
    FOR i IN 1..24 LOOP
      -- Lấy school_id của parent
      DECLARE
        p_school uuid;
        rel text;
      BEGIN
        SELECT school_id INTO p_school FROM parents WHERE parent_id = v_parent_ids[i];

        -- Gán parent này với 2 students trong cùng school
        j := 0;
        FOR class_records IN
          SELECT s.student_id
          FROM students s
          WHERE s.school_id = p_school
          ORDER BY s.full_name
        LOOP
          j := j + 1;
          IF j > 2 THEN EXIT; END IF;

          rel := CASE WHEN j = 1 THEN 'father' ELSE 'mother' END;

          INSERT INTO student_parents (student_id, parent_id, relationship)
          VALUES (class_records.student_id, v_parent_ids[i], rel)
          ON CONFLICT (student_id, parent_id) DO NOTHING;
        END LOOP;
      END;
    END LOOP;
  END IF;

  -- ======================
  -- 8) STUDENT PARENT CODES (mỗi học sinh 1 code)
  -- ======================
  IF to_regclass('public.student_parent_codes') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
      VALUES (
        v_student_ids[i],
        'PC' || LPAD(i::text, 4, '0'),
        4,
        NOW() + INTERVAL '90 days'
      )
      ON CONFLICT (code) DO NOTHING;
    END LOOP;
  END IF;

  RAISE NOTICE '✅ seed_02: 1 super admin + 8 school admins + 24 teachers + 24 parents + 56 students + assignments + codes inserted';
END $$;