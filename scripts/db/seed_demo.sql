DO $$
DECLARE
  -- Schools
  v_school1 uuid;
  v_school2 uuid;
  v_school3 uuid;

  -- Classes
  v_class1  uuid; -- Lá Non (school1)
  v_class2  uuid; -- Lá Măng (school1)
  v_class3  uuid; -- Lá Sen (school2)
  v_class4  uuid; -- Chồi Xanh (school3)
  v_class5  uuid; -- Mầm Sáng (school3)

  -- Users
  v_super_admin_user uuid;
  v_school_admin_user uuid;
  v_teacher1_user     uuid;
  v_teacher2_user     uuid;
  v_teacher3_user     uuid;
  v_parent1_user      uuid;
  v_parent2_user      uuid;
  v_parent3_user      uuid;
  v_school_admin2_user uuid;

  -- Teachers
  v_teacher1 uuid; -- Cô Lan
  v_teacher2 uuid; -- Thầy Nam
  v_teacher3 uuid; -- Cô Phương

  -- Parents
  v_parent1 uuid; -- Anh Minh
  v_parent2 uuid; -- Chị Hoa
  v_parent3 uuid; -- Anh Quân

  -- School Admins
  v_school_admin uuid; -- Cô Hương (SCHOOL_ADMIN của school1)
  v_school_admin2 uuid; -- Thầy Duy (SCHOOL_ADMIN của school3)

  -- Students
  v_s1 uuid; -- Bé An
  v_s2 uuid; -- Bé Bông
  v_s3 uuid; -- Bé Na
  v_s4 uuid; -- Bé Dương
  v_s5 uuid; -- Bé My
  v_s6 uuid; -- Bé Phúc
  v_s7 uuid; -- Bé Khôi
  v_s8 uuid; -- Bé Trâm
  v_s9 uuid; -- Bé Bảo
  v_s10 uuid; -- Bé Vy

BEGIN
  -- ======================
  -- 1) SCHOOLS
  -- ======================
  SELECT school_id INTO v_school1
  FROM schools
  WHERE name = 'Trường Mầm non Bình Minh Cầu Giấy'
  ORDER BY created_at
  LIMIT 1;

  IF v_school1 IS NULL THEN
    INSERT INTO schools (name, address)
    VALUES ('Trường Mầm non Bình Minh Cầu Giấy', 'Số 18 ngõ 165 Xuân Thủy, Dịch Vọng Hậu, Cầu Giấy, Hà Nội')
    RETURNING school_id INTO v_school1;
  END IF;

  SELECT school_id INTO v_school2
  FROM schools
  WHERE name = 'Trường Mầm non Hoa Mai Quận 7'
  ORDER BY created_at
  LIMIT 1;

  IF v_school2 IS NULL THEN
    INSERT INTO schools (name, address)
    VALUES ('Trường Mầm non Hoa Mai Quận 7', 'Số 12 đường số 3, KDC Him Lam, Tân Hưng, Quận 7, TP. Hồ Chí Minh')
    RETURNING school_id INTO v_school2;
  END IF;

  SELECT school_id INTO v_school3
  FROM schools
  WHERE name = 'Trường Mầm non Ánh Dương Hải Châu'
  ORDER BY created_at
  LIMIT 1;

  IF v_school3 IS NULL THEN
    INSERT INTO schools (name, address)
    VALUES ('Trường Mầm non Ánh Dương Hải Châu', 'Số 25 Lê Đình Dương, Phước Ninh, Hải Châu, Đà Nẵng')
    RETURNING school_id INTO v_school3;
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

  SELECT class_id INTO v_class4
  FROM classes
  WHERE school_id = v_school3 AND name = 'Chồi Xanh'
  ORDER BY created_at
  LIMIT 1;

  IF v_class4 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school3, 'Chồi Xanh', '2025-2026')
    RETURNING class_id INTO v_class4;
  END IF;

  SELECT class_id INTO v_class5
  FROM classes
  WHERE school_id = v_school3 AND name = 'Mầm Sáng'
  ORDER BY created_at
  LIMIT 1;

  IF v_class5 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school3, 'Mầm Sáng', '2025-2026')
    RETURNING class_id INTO v_class5;
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
  VALUES ('teacher3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_teacher3_user;

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

  INSERT INTO users (email, password_hash, status)
  VALUES ('parent3@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_parent3_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('school-admin2@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_school_admin2_user;

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
  SELECT v_teacher3_user, role_id FROM roles WHERE name = 'TEACHER'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_parent1_user, role_id FROM roles WHERE name = 'PARENT'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_parent2_user, role_id FROM roles WHERE name = 'PARENT'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_parent3_user, role_id FROM roles WHERE name = 'PARENT'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_school_admin2_user, role_id FROM roles WHERE name = 'SCHOOL_ADMIN'
  ON CONFLICT DO NOTHING;

  -- ======================
  -- 5) TEACHER PROFILES
  -- ======================
  SELECT teacher_id INTO v_teacher1 FROM teachers WHERE user_id = v_teacher1_user;
  IF v_teacher1 IS NULL THEN
    INSERT INTO teachers (user_id, school_id, full_name, phone)
    VALUES (v_teacher1_user, v_school1, 'Nguyễn Thị Lan Anh', '0983124501')
    RETURNING teacher_id INTO v_teacher1;
  END IF;

  SELECT teacher_id INTO v_teacher2 FROM teachers WHERE user_id = v_teacher2_user;
  IF v_teacher2 IS NULL THEN
    INSERT INTO teachers (user_id, school_id, full_name, phone)
    VALUES (v_teacher2_user, v_school1, 'Trần Hoàng Nam', '0917642208')
    RETURNING teacher_id INTO v_teacher2;
  END IF;

  SELECT teacher_id INTO v_teacher3 FROM teachers WHERE user_id = v_teacher3_user;
  IF v_teacher3 IS NULL THEN
    INSERT INTO teachers (user_id, school_id, full_name, phone)
    VALUES (v_teacher3_user, v_school3, 'Đặng Minh Phương', '0908765123')
    RETURNING teacher_id INTO v_teacher3;
  END IF;

  -- ======================
  -- 6) PARENT PROFILES
  -- ======================
  SELECT parent_id INTO v_parent1 FROM parents WHERE user_id = v_parent1_user;
  IF v_parent1 IS NULL THEN
    INSERT INTO parents (user_id, school_id, full_name, phone)
    VALUES (v_parent1_user, v_school1, 'Lê Quang Minh', '0935587612')
    RETURNING parent_id INTO v_parent1;
  END IF;

  SELECT parent_id INTO v_parent2 FROM parents WHERE user_id = v_parent2_user;
  IF v_parent2 IS NULL THEN
    INSERT INTO parents (user_id, school_id, full_name, phone)
    VALUES (v_parent2_user, v_school1, 'Phạm Thu Hoa', '0964412097')
    RETURNING parent_id INTO v_parent2;
  END IF;

  SELECT parent_id INTO v_parent3 FROM parents WHERE user_id = v_parent3_user;
  IF v_parent3 IS NULL THEN
    INSERT INTO parents (user_id, school_id, full_name, phone)
    VALUES (v_parent3_user, v_school3, 'Võ Văn Quân', '0978123456')
    RETURNING parent_id INTO v_parent3;
  END IF;

  -- ======================
  -- 7) SCHOOL ADMIN PROFILES
  -- ======================
  SELECT admin_id INTO v_school_admin FROM school_admins WHERE user_id = v_school_admin_user;
  IF v_school_admin IS NULL THEN
    INSERT INTO school_admins (user_id, school_id, full_name, phone)
    VALUES (v_school_admin_user, v_school1, 'Bùi Thanh Hương', '0972331184')
    RETURNING admin_id INTO v_school_admin;
  END IF;

  SELECT admin_id INTO v_school_admin2 FROM school_admins WHERE user_id = v_school_admin2_user;
  IF v_school_admin2 IS NULL THEN
    INSERT INTO school_admins (user_id, school_id, full_name, phone)
    VALUES (v_school_admin2_user, v_school3, 'Ngô Quang Duy', '0933123123')
    RETURNING admin_id INTO v_school_admin2;
  END IF;

  -- ======================
  -- 8) STUDENTS
  -- ======================
  SELECT student_id INTO v_s1 FROM students
  WHERE school_id = v_school1 AND full_name = 'Nguyễn Gia An' AND dob = DATE '2020-05-12'
  LIMIT 1;
  IF v_s1 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class1, 'Nguyễn Gia An', DATE '2020-05-12', 'male')
    RETURNING student_id INTO v_s1;
  END IF;

  SELECT student_id INTO v_s2 FROM students
  WHERE school_id = v_school1 AND full_name = 'Trần Khánh Băng' AND dob = DATE '2020-09-03'
  LIMIT 1;
  IF v_s2 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class1, 'Trần Khánh Băng', DATE '2020-09-03', 'female')
    RETURNING student_id INTO v_s2;
  END IF;

  SELECT student_id INTO v_s3 FROM students
  WHERE school_id = v_school1 AND full_name = 'Lê Ngọc Nhi' AND dob = DATE '2021-01-20'
  LIMIT 1;
  IF v_s3 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class2, 'Lê Ngọc Nhi', DATE '2021-01-20', 'female')
    RETURNING student_id INTO v_s3;
  END IF;

  SELECT student_id INTO v_s4 FROM students
  WHERE school_id = v_school1 AND full_name = 'Phạm Đức Dương' AND dob = DATE '2020-11-15'
  LIMIT 1;
  IF v_s4 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class2, 'Phạm Đức Dương', DATE '2020-11-15', 'male')
    RETURNING student_id INTO v_s4;
  END IF;

  SELECT student_id INTO v_s5 FROM students
  WHERE school_id = v_school1 AND full_name = 'Đỗ Minh Khang' AND dob = DATE '2020-07-08'
  LIMIT 1;
  IF v_s5 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class2, 'Đỗ Minh Khang', DATE '2020-07-08', 'male')
    RETURNING student_id INTO v_s5;
  END IF;

  SELECT student_id INTO v_s6 FROM students
  WHERE school_id = v_school1 AND full_name = 'Ngô Phương My' AND dob = DATE '2021-02-14'
  LIMIT 1;
  IF v_s6 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school1, v_class1, 'Ngô Phương My', DATE '2021-02-14', 'female')
    RETURNING student_id INTO v_s6;
  END IF;

  SELECT student_id INTO v_s7 FROM students
  WHERE school_id = v_school2 AND full_name = 'Huỳnh Bảo Khôi' AND dob = DATE '2020-03-19'
  LIMIT 1;
  IF v_s7 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school2, v_class3, 'Huỳnh Bảo Khôi', DATE '2020-03-19', 'male')
    RETURNING student_id INTO v_s7;
  END IF;

  SELECT student_id INTO v_s8 FROM students
  WHERE school_id = v_school2 AND full_name = 'Đinh Gia Trâm' AND dob = DATE '2020-12-02'
  LIMIT 1;
  IF v_s8 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school2, v_class3, 'Đinh Gia Trâm', DATE '2020-12-02', 'female')
    RETURNING student_id INTO v_s8;
  END IF;

  SELECT student_id INTO v_s9 FROM students
  WHERE school_id = v_school3 AND full_name = 'Lý Minh Bảo' AND dob = DATE '2021-04-11'
  LIMIT 1;
  IF v_s9 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school3, v_class4, 'Lý Minh Bảo', DATE '2021-04-11', 'male')
    RETURNING student_id INTO v_s9;
  END IF;

  SELECT student_id INTO v_s10 FROM students
  WHERE school_id = v_school3 AND full_name = 'Trương Hạ Vy' AND dob = DATE '2020-10-27'
  LIMIT 1;
  IF v_s10 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school3, v_class5, 'Trương Hạ Vy', DATE '2020-10-27', 'female')
    RETURNING student_id INTO v_s10;
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

    -- Cô Phương → Chồi Xanh + Mầm Sáng
    INSERT INTO teacher_classes (teacher_id, class_id)
    VALUES (v_teacher3, v_class4)
    ON CONFLICT DO NOTHING;

    INSERT INTO teacher_classes (teacher_id, class_id)
    VALUES (v_teacher3, v_class5)
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

    -- Bé Khang: father (Anh Minh)
    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s5, v_parent1, 'father')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    -- Bé My: mother (Chị Hoa)
    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s6, v_parent2, 'mother')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    -- Bé Bảo + Bé Vy: father (Anh Quân)
    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s9, v_parent3, 'father')
    ON CONFLICT (student_id, parent_id) DO NOTHING;

    INSERT INTO student_parents (student_id, parent_id, relationship)
    VALUES (v_s10, v_parent3, 'father')
    ON CONFLICT (student_id, parent_id) DO NOTHING;
  END IF;

  -- ======================
  -- 11) PARENT CODES
  -- ======================
  IF to_regclass('public.student_parent_codes') IS NOT NULL THEN
    -- Parent code for Nguyễn Gia An (max 4 parents, expires in 30 days)
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s1, 'NGA2401', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    -- Parent code for Trần Khánh Băng
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s2, 'TKB2402', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    -- Parent code for Lê Ngọc Nhi
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s3, 'LNN2403', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    -- Parent code for Phạm Đức Dương
    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s4, 'PDD2404', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s5, 'DMK2405', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s6, 'NPM2406', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s7, 'HBK2407', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s8, 'DGT2408', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s9, 'LMB2409', 4, NOW() + INTERVAL '30 days')
    ON CONFLICT (code) DO NOTHING;

    INSERT INTO student_parent_codes (student_id, code, max_usage, expires_at)
    VALUES (v_s10, 'THV2410', 4, NOW() + INTERVAL '30 days')
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
      'Đến lớp đúng giờ, tinh thần vui vẻ',
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
      'Đến lớp muộn do mưa lớn buổi sáng',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé An: absent (2 ngày trước)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE - INTERVAL '2 days',
      'absent',
      'Nghỉ học không phép, nhà trường đã liên hệ gia đình',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Bông: excused (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s2,
      CURRENT_DATE,
      'excused',
      'Nghỉ ốm có xác nhận phụ huynh qua ứng dụng',
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
      'Nghỉ do việc gia đình, đã báo trước với giáo viên chủ nhiệm',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Khang: present (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
    VALUES (
      v_s5,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:12:00')::timestamptz,
      (CURRENT_DATE + TIME '16:08:00')::timestamptz,
      'Hoàn thành tốt hoạt động nhóm, ăn trưa đủ suất',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé My: late (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
    VALUES (
      v_s6,
      CURRENT_DATE,
      'late',
      (CURRENT_DATE + TIME '09:05:00')::timestamptz,
      'Đến lớp muộn do kẹt xe giờ cao điểm',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Khôi: present (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
    VALUES (
      v_s7,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:20:00')::timestamptz,
      (CURRENT_DATE + TIME '16:20:00')::timestamptz,
      'Tự giác dọn đồ chơi sau giờ học',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Trâm: excused (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s8,
      CURRENT_DATE,
      'excused',
      'Nghỉ khám răng định kỳ, phụ huynh đã báo trước',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Bảo: present (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
    VALUES (
      v_s9,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:18:00')::timestamptz,
      'Ổn định nề nếp đầu giờ tốt',
      v_teacher3_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Vy: present (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
    VALUES (
      v_s10,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:26:00')::timestamptz,
      'Tham gia tích cực tiết âm nhạc',
      v_teacher3_user
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
      'Ăn ngủ bình thường, hoạt động thể chất tốt',
      v_teacher1_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé An: watch (hôm qua - sốt nhẹ)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s1,
      NOW() - INTERVAL '1 day',
      37.8,
      'Ho nhẹ, sổ mũi',
      'watch',
      'Đã theo dõi 2 giờ, chưa cần đưa đi khám',
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
      'Sốt cao, mệt và đau đầu',
      'urgent',
      'Đã gọi phụ huynh đón về và hướng dẫn theo dõi tại nhà',
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

    -- Bé Khang: watch (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s5,
      NOW(),
      37.4,
      'Ho nhẹ buổi chiều',
      'watch',
      'Đã nhắc uống nước ấm và theo dõi thêm',
      v_teacher2_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé My: normal (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, severity, note, recorded_by)
    VALUES (
      v_s6,
      NOW(),
      36.4,
      'normal',
      'Không ghi nhận dấu hiệu bất thường',
      v_teacher1_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Khôi: normal (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, severity, note, recorded_by)
    VALUES (
      v_s7,
      NOW(),
      36.6,
      'normal',
      'Ăn sáng đầy đủ trước khi đến lớp',
      v_teacher2_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Trâm: watch (2 ngày trước)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s8,
      NOW() - INTERVAL '2 days',
      37.6,
      'Hắt hơi, nghẹt mũi',
      'watch',
      'Theo dõi tại lớp, tình trạng cải thiện sau giờ nghỉ trưa',
      v_teacher2_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Bảo: normal (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, severity, note, recorded_by)
    VALUES (
      v_s9,
      NOW(),
      36.5,
      'normal',
      'Thể trạng ổn định',
      v_teacher3_user
    )
    ON CONFLICT DO NOTHING;

    -- Bé Vy: normal (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, severity, note, recorded_by)
    VALUES (
      v_s10,
      NOW(),
      36.7,
      'normal',
      'Tinh thần vui vẻ, hợp tác tốt',
      v_teacher3_user
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
      'Thông báo: Thứ Sáu tuần này lớp Lá Non tham gia hoạt động trải nghiệm tại sân trường. Phụ huynh chuẩn bị mũ, bình nước cá nhân và giày thể thao cho các con.'
    )
    ON CONFLICT DO NOTHING;

    -- Class post: activity (Lá Non)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher1_user,
      'class',
      v_class1,
      'activity',
      'Hoạt động vẽ tranh sáng nay các con rất hào hứng, biết chia sẻ bút màu và hỗ trợ bạn cùng nhóm. Nhà trường đã cập nhật ảnh trong mục album lớp.'
    )
    ON CONFLICT DO NOTHING;

    -- Student post: daily_note (Bé An)
    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (
      v_teacher1_user,
      'student',
      v_s1,
      'daily_note',
      'Gia An hôm nay tự giác ăn hết suất trưa, ngủ trưa đúng giờ và hợp tác tốt trong giờ học kỹ năng. Đề nghị gia đình duy trì nếp sinh hoạt tối trước 21h.'
    )
    ON CONFLICT DO NOTHING;

    -- Student post: health_note (Bé Na - urgent từ health log)
    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (
      v_teacher2_user,
      'student',
      v_s3,
      'health_note',
      'Ngọc Nhi có biểu hiện sốt cao 39.2°C vào 10:20, giáo viên đã liên hệ gia đình và hỗ trợ bé nghỉ tại phòng y tế trước khi phụ huynh đến đón.'
    )
    ON CONFLICT DO NOTHING;

    -- Class post: health_note (Lá Măng)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher2_user,
      'class',
      v_class2,
      'health_note',
      'Lưu ý sức khỏe: Hôm nay lớp có 1 trường hợp ho nhẹ. Giáo viên đã nhắc cả lớp rửa tay đúng cách và vệ sinh cá nhân. Phụ huynh theo dõi thêm triệu chứng tại nhà.'
    )
    ON CONFLICT DO NOTHING;

    -- School post: announcement (school3)
    INSERT INTO posts (author_user_id, scope_type, school_id, type, content)
    VALUES (
      v_school_admin2_user,
      'school',
      v_school3,
      'announcement',
      'Thông báo toàn trường: Tuần tới tổ chức chuyên đề “An toàn giao thông cho bé”. Đề nghị phụ huynh phối hợp đưa đón đúng khu vực cổng trường theo phân luồng.'
    )
    ON CONFLICT DO NOTHING;

    -- Class post: activity (Chồi Xanh)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher3_user,
      'class',
      v_class4,
      'activity',
      'Lớp Chồi Xanh hôm nay thực hành kỹ năng tự phục vụ: xếp dép đúng vị trí, cất balo gọn gàng và rửa tay theo 6 bước trước giờ ăn.'
    )
    ON CONFLICT DO NOTHING;

    -- Class post: announcement (Mầm Sáng)
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (
      v_teacher3_user,
      'class',
      v_class5,
      'announcement',
      'Nhắc lịch: Thứ Tư các con mặc đồng phục thể dục để tham gia hoạt động vận động ngoài trời lúc 8:30.'
    )
    ON CONFLICT DO NOTHING;

    -- Student post: daily_note (Bé Khang)
    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (
      v_teacher2_user,
      'student',
      v_s5,
      'daily_note',
      'Minh Khang hôm nay chủ động chào cô, tập trung tốt trong tiết kể chuyện và biết hỗ trợ bạn thu dọn đồ dùng học tập.'
    )
    ON CONFLICT DO NOTHING;

    -- Student post: daily_note (Bé Vy)
    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (
      v_teacher3_user,
      'student',
      v_s10,
      'daily_note',
      'Hạ Vy tích cực tham gia hoạt động âm nhạc, phát âm rõ khi hát tập thể và phối hợp tốt với nhóm trong giờ thủ công.'
    )
    ON CONFLICT DO NOTHING;
  END IF;

END $$;
