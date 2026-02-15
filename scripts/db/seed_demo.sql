DO $$
DECLARE
  -- Schools
  v_school1 uuid;
  v_school2 uuid;

  -- Classes
  v_class1  uuid; -- Lá Non (school1)
  v_class2  uuid; -- Lá Măng (school1)
  v_class3  uuid; -- Lá Sen (school2)

  -- Users
  v_super_admin_user uuid;
  v_school_admin_user uuid;
  v_teacher1_user     uuid;
  v_teacher2_user     uuid;
  v_parent1_user      uuid;
  v_parent2_user      uuid;

  -- Teachers
  v_teacher1 uuid; -- Cô Lan
  v_teacher2 uuid; -- Thầy Nam

  -- Parents
  v_parent1 uuid; -- Anh Minh
  v_parent2 uuid; -- Chị Hoa

  -- School Admins
  v_school_admin uuid; -- Cô Hương (SCHOOL_ADMIN của school1)

  -- Students
  v_s1 uuid; -- Bé An
  v_s2 uuid; -- Bé Bông
  v_s3 uuid; -- Bé Na
  v_s4 uuid; -- Bé Dương

  -- Parent Codes
  v_code1 uuid; -- Parent code for Bé An
  v_code2 uuid; -- Parent code for Bé Bông
BEGIN
  -- ======================
  -- 1) SCHOOLS
  -- ======================
  SELECT school_id INTO v_school1
  FROM schools
  WHERE name = 'IRIS Demo School'
  ORDER BY created_at
  LIMIT 1;

  IF v_school1 IS NULL THEN
    INSERT INTO schools (name, address)
    VALUES ('IRIS Demo School', 'Hà Nội')
    RETURNING school_id INTO v_school1;
  END IF;

  SELECT school_id INTO v_school2
  FROM schools
  WHERE name = 'Mầm Non Hoa Mai'
  ORDER BY created_at
  LIMIT 1;

  IF v_school2 IS NULL THEN
    INSERT INTO schools (name, address)
    VALUES ('Mầm Non Hoa Mai', 'Hồ Chí Minh')
    RETURNING school_id INTO v_school2;
  END IF;

  -- ======================
  -- 2) CLASSES
  -- ======================
  SELECT class_id INTO v_class1
  FROM classes
  WHERE school_id = v_school1 AND name = 'Lá Non'
  ORDER BY created_at
  LIMIT 1;

  IF v_class1 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school1, 'Lá Non', '2025-2026')
    RETURNING class_id INTO v_class1;
  END IF;

  SELECT class_id INTO v_class2
  FROM classes
  WHERE school_id = v_school1 AND name = 'Lá Măng'
  ORDER BY created_at
  LIMIT 1;

  IF v_class2 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school1, 'Lá Măng', '2025-2026')
    RETURNING class_id INTO v_class2;
  END IF;

  SELECT class_id INTO v_class3
  FROM classes
  WHERE school_id = v_school2 AND name = 'Lá Sen'
  ORDER BY created_at
  LIMIT 1;

  IF v_class3 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school2, 'Lá Sen', '2025-2026')
    RETURNING class_id INTO v_class3;
  END IF;

  -- ======================
  -- 3) USERS
  -- ======================
  -- bcrypt hash for "123456"
  INSERT INTO users (email, password_hash, status)
  VALUES ('admin@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_super_admin_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('school-admin@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_school_admin_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('teacher1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_teacher1_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('teacher2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_teacher2_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('parent1@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_parent1_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('parent2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_parent2_user;

  -- ======================
  -- 4) USER ROLES
  -- ======================
  INSERT INTO user_roles (user_id, role_id)
  SELECT v_super_admin_user, role_id FROM roles WHERE name = 'SUPER_ADMIN'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_school_admin_user, role_id FROM roles WHERE name = 'SCHOOL_ADMIN'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_teacher1_user, role_id FROM roles WHERE name = 'TEACHER'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_teacher2_user, role_id FROM roles WHERE name = 'TEACHER'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_parent1_user, role_id FROM roles WHERE name = 'PARENT'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_parent2_user, role_id FROM roles WHERE name = 'PARENT'
  ON CONFLICT DO NOTHING;

  -- ======================
  -- 5) TEACHER PROFILES
  -- ======================
  SELECT teacher_id INTO v_teacher1 FROM teachers WHERE user_id = v_teacher1_user;
  IF v_teacher1 IS NULL THEN
    INSERT INTO teachers (user_id, school_id, full_name, phone)
    VALUES (v_teacher1_user, v_school1, 'Cô Lan', '0900000001')
    RETURNING teacher_id INTO v_teacher1;
  END IF;

  SELECT teacher_id INTO v_teacher2 FROM teachers WHERE user_id = v_teacher2_user;
  IF v_teacher2 IS NULL THEN
    INSERT INTO teachers (user_id, school_id, full_name, phone)
    VALUES (v_teacher2_user, v_school1, 'Thầy Nam', '0900000002')
    RETURNING teacher_id INTO v_teacher2;
  END IF;

  -- ======================
  -- 6) PARENT PROFILES
  -- ======================
  SELECT parent_id INTO v_parent1 FROM parents WHERE user_id = v_parent1_user;
  IF v_parent1 IS NULL THEN
    INSERT INTO parents (user_id, school_id, full_name, phone)
    VALUES (v_parent1_user, v_school1, 'Anh Minh', '0900000003')
    RETURNING parent_id INTO v_parent1;
  END IF;

  SELECT parent_id INTO v_parent2 FROM parents WHERE user_id = v_parent2_user;
  IF v_parent2 IS NULL THEN
    INSERT INTO parents (user_id, school_id, full_name, phone)
    VALUES (v_parent2_user, v_school1, 'Chị Hoa', '0900000004')
    RETURNING parent_id INTO v_parent2;
  END IF;

  -- ======================
  -- 7) SCHOOL ADMIN PROFILES
  -- ======================
  SELECT admin_id INTO v_school_admin FROM school_admins WHERE user_id = v_school_admin_user;
  IF v_school_admin IS NULL THEN
    INSERT INTO school_admins (user_id, school_id, full_name, phone)
    VALUES (v_school_admin_user, v_school1, 'Cô Hương', '0900000005')
    RETURNING admin_id INTO v_school_admin;
  END IF;

  -- ======================
  -- 8) STUDENTS
  -- ======================
  SELECT student_id INTO v_s1 FROM students
  WHERE school_id = v_school1 AND full_name = 'Bé An' AND dob = DATE '2021-05-12'
  LIMIT 1;
  IF v_s1 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class1, 'Bé An', DATE '2021-05-12', 'male')
    RETURNING student_id INTO v_s1;
  END IF;

  SELECT student_id INTO v_s2 FROM students
  WHERE school_id = v_school1 AND full_name = 'Bé Bông' AND dob = DATE '2021-09-03'
  LIMIT 1;
  IF v_s2 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class1, 'Bé Bông', DATE '2021-09-03', 'female')
    RETURNING student_id INTO v_s2;
  END IF;

  SELECT student_id INTO v_s3 FROM students
  WHERE school_id = v_school1 AND full_name = 'Bé Na' AND dob = DATE '2022-01-20'
  LIMIT 1;
  IF v_s3 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class2, 'Bé Na', DATE '2022-01-20', 'female')
    RETURNING student_id INTO v_s3;
  END IF;

  SELECT student_id INTO v_s4 FROM students
  WHERE school_id = v_school1 AND full_name = 'Bé Dương' AND dob = DATE '2021-11-15'
  LIMIT 1;
  IF v_s4 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class2, 'Bé Dương', DATE '2021-11-15', 'male')
    RETURNING student_id INTO v_s4;
  END IF;

  -- ======================
  -- 9) TEACHER-CLASS ASSIGNMENTS
  -- ======================
  IF to_regclass('public.teacher_classes') IS NOT NULL THEN
    -- Cô Lan → Lá Non
    INSERT INTO teacher_classes (teacher_id, class_id)
    VALUES (v_teacher1, v_class1)
    ON CONFLICT DO NOTHING;

    -- Thầy Nam → Lá Non + Lá Măng (multiple classes)
    INSERT INTO teacher_classes (teacher_id, class_id)
    VALUES (v_teacher2, v_class1)
    ON CONFLICT DO NOTHING;

    INSERT INTO teacher_classes (teacher_id, class_id)
    VALUES (v_teacher2, v_class2)
    ON CONFLICT DO NOTHING;
  END IF;

  -- ======================
  -- 10) STUDENT-PARENT ASSIGNMENTS
  -- ======================
  IF to_regclass('public.student_parents') IS NOT NULL THEN
    -- Bé An: father (Anh Minh) + mother (Chị Hoa)
    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s1, v_parent1, 'father')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s1, v_parent2, 'mother')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    -- Bé Bông: father only (Anh Minh)
    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s2, v_parent1, 'father')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    -- Bé Na: mother only (Chị Hoa)
    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s3, v_parent2, 'mother')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    -- Bé Dương: no parents
  END IF;

  -- ======================
  -- 11) PARENT CODES
  -- ======================
  IF to_regclass('public.student_parent_codes') IS NOT NULL THEN
    -- Parent code for Bé An (max 4 parents, expires in 30 days)
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s1, 'BEAN25', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    -- Parent code for Bé Bông
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s2, 'BEBONG', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    -- Parent code for Bé Na
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s3, 'BENA25', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    -- Parent code for Bé Dương
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s4, 'BEDUONG', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;
  END IF;

  -- ======================
  -- 12) ATTENDANCE RECORDS (đa dạng status, có check-in/out, date range)
  -- ======================
  IF to_regclass('public.attendance_records') IS NOT NULL THEN
    -- Bé An: present (hôm nay, có check-in/out)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:30:00')::timestamptz,
      (CURRENT_DATE + TIME '16:30:00')::timestamptz,
      'Đến đúng giờ',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé An: late (hôm qua)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE - INTERVAL '1 day',
      'late',
      (CURRENT_DATE - INTERVAL '1 day' + TIME '09:15:00')::timestamptz,
      'Đi muộn 45 phút',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé An: absent (2 ngày trước)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE - INTERVAL '2 days',
      'absent',
      'Nghỉ không phép',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Bông: excused (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s2,
      CURRENT_DATE,
      'excused',
      'Nghỉ ốm có giấy',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Bông: present (hôm qua)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
    VALUES (
      v_s2,
      CURRENT_DATE - INTERVAL '1 day',
      'present',
      (CURRENT_DATE - INTERVAL '1 day' + TIME '08:00:00')::timestamptz,
      (CURRENT_DATE - INTERVAL '1 day' + TIME '16:00:00')::timestamptz,
      '',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Na: present (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
    VALUES (
      v_s3,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:45:00')::timestamptz,
      '',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Dương: absent (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s4,
      CURRENT_DATE,
      'absent',
      'Nghỉ gia đình',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;
  END IF;

  -- ======================
  -- 13) HEALTH LOGS (đa dạng severity, date range)
  -- ======================
  IF to_regclass('public.health_logs') IS NOT NULL THEN
    -- Bé An: normal (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s1,
      NOW(),
      36.5,
      '',
      'normal',
      'Sức khỏe tốt',
      v_teacher1_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé An: watch (hôm qua - sốt nhẹ)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s1,
      NOW() - INTERVAL '1 day',
      37.8,
      'Ho nhẹ',
      'watch',
      'Cần theo dõi',
      v_teacher1_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Bông: normal (3 ngày trước)
    INSERT INTO health_logs (student_id, recorded_at, temperature, severity, recorded_by)
    VALUES (
      v_s2,
      NOW() - INTERVAL '3 days',
      36.3,
      'normal',
      v_teacher1_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Na: urgent (hôm nay - sốt cao)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s3,
      NOW(),
      39.2,
      'Sốt cao, đau đầu',
      'urgent',
      'Cần liên hệ gia đình',
      v_teacher2_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Dương: normal (tuần trước)
    INSERT INTO health_logs (student_id, recorded_at, temperature, severity, recorded_by)
    VALUES (
      v_s4,
      NOW() - INTERVAL '7 days',
      36.7,
      'normal',
      v_teacher2_user
    )
    ON CONFLICT DO NOTHING;
  END IF;

  -- ======================
  -- 14) POSTS (đa dạng type, scope: class/student)
  -- ======================
  IF to_regclass('public.posts') IS NOT NULL THEN
    -- Class post: announcement (Lá Non)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher1_user,
      'class',
      v_class1,
      'announcement',
      'Thông báo: Ngày mai lớp Lá Non có hoạt động ngoại khóa. Phụ huynh chuẩn bị quần áo thể thao cho các bé.'
    )
    ON CONFLICT DO NOTHING;

    -- Class post: activity (Lá Non)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher1_user,
      'class',
      v_class1,
      'activity',
      'Hoạt động vẽ tranh sáng nay các bạn rất hào hứng. Mời phụ huynh xem thêm ảnh đã gửi vào group!'
    )
    ON CONFLICT DO NOTHING;

    -- Student post: daily_note (Bé An)
    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (
      v_teacher1_user,
      'student',
      v_s1,
      'daily_note',
      'Bé An hôm nay tự giác ăn hết cơm trưa và ngủ trưa rất ngoan. Tiếp tục cố gắng nhé!'
    )
    ON CONFLICT DO NOTHING;

    -- Student post: health_note (Bé Na - urgent từ health log)
    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (
      v_teacher2_user,
      'student',
      v_s3,
      'health_note',
      'Bé Na bị sốt cao (39.2°C), đã báo phụ huynh. Bé đang ở y trường để theo dõi.'
    )
    ON CONFLICT DO NOTHING;

    -- Class post: health_note (Lá Măng)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher2_user,
      'class',
      v_class2,
      'health_note',
      'Lưu ý: Hôm nay trong lớp có 1 bạn bị ho nhẹ. Các bạn khác chú ý vệ sinh, đeo khẩu trang nếu cần.'
    )
    ON CONFLICT DO NOTHING;
  END IF;

END $$;
