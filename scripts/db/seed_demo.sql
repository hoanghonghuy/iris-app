DO $$
DECLARE
  v_school uuid;
  v_class  uuid;

  v_admin_user   uuid;
  v_teacher_user uuid;
  v_parent_user  uuid;

  v_teacher uuid;
  v_parent  uuid;

  v_s1 uuid;
  v_s2 uuid;
  v_s3 uuid;
BEGIN
  -- 1) School
  SELECT school_id INTO v_school
  FROM schools
  WHERE name = 'IRIS Demo School'
  ORDER BY created_at
  LIMIT 1;

  IF v_school IS NULL THEN
    INSERT INTO schools (name, address)
    VALUES ('IRIS Demo School', 'Hà Nội')
    RETURNING school_id INTO v_school;
  END IF;

  -- 2) Class
  SELECT class_id INTO v_class
  FROM classes
  WHERE school_id = v_school AND name = 'Lá Non' AND school_year = '2025-2026'
  ORDER BY created_at
  LIMIT 1;

  IF v_class IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school, 'Lá Non', '2025-2026')
    RETURNING class_id INTO v_class;
  END IF;

  -- 3) Users (password: 123456)
  -- bcrypt hash for "123456"
  INSERT INTO users (email, password_hash, status)
  VALUES ('admin@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_admin_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('teacher@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_teacher_user;

  INSERT INTO users (email, password_hash, status)
  VALUES ('parent@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_parent_user;

  -- 4) User roles
  INSERT INTO user_roles (user_id, role_id)
  SELECT v_admin_user, role_id FROM roles WHERE name = 'ADMIN'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_teacher_user, role_id FROM roles WHERE name = 'TEACHER'
  ON CONFLICT DO NOTHING;

  INSERT INTO user_roles (user_id, role_id)
  SELECT v_parent_user, role_id FROM roles WHERE name = 'PARENT'
  ON CONFLICT DO NOTHING;

  -- 5) Teacher/Parent profiles
  SELECT teacher_id INTO v_teacher FROM teachers WHERE user_id = v_teacher_user;
  IF v_teacher IS NULL THEN
    INSERT INTO teachers (user_id, school_id, full_name, phone)
    VALUES (v_teacher_user, v_school, 'Cô Lan', '0900000001')
    RETURNING teacher_id INTO v_teacher;
  END IF;

  SELECT parent_id INTO v_parent FROM parents WHERE user_id = v_parent_user;
  IF v_parent IS NULL THEN
    INSERT INTO parents (user_id, school_id, full_name, phone)
    VALUES (v_parent_user, v_school, 'Anh Minh', '0900000002')
    RETURNING parent_id INTO v_parent;
  END IF;

  -- 6) Students
  SELECT student_id INTO v_s1 FROM students
  WHERE school_id = v_school AND full_name = 'Bé An' AND dob = DATE '2021-05-12'
  LIMIT 1;
  IF v_s1 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school, v_class, 'Bé An', DATE '2021-05-12', 'male')
    RETURNING student_id INTO v_s1;
  END IF;

  SELECT student_id INTO v_s2 FROM students
  WHERE school_id = v_school AND full_name = 'Bé Bông' AND dob = DATE '2021-09-03'
  LIMIT 1;
  IF v_s2 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school, v_class, 'Bé Bông', DATE '2021-09-03', 'female')
    RETURNING student_id INTO v_s2;
  END IF;

  SELECT student_id INTO v_s3 FROM students
  WHERE school_id = v_school AND full_name = 'Bé Na' AND dob = DATE '2022-01-20'
  LIMIT 1;
  IF v_s3 IS NULL THEN
    INSERT INTO students (school_id, current_class_id, full_name, dob, gender)
    VALUES (v_school, v_class, 'Bé Na', DATE '2022-01-20', 'other')
    RETURNING student_id INTO v_s3;
  END IF;

  -- 7) Relations
  INSERT INTO teacher_classes (teacher_id, class_id)
  VALUES (v_teacher, v_class)
  ON CONFLICT DO NOTHING;

  INSERT INTO student_parents (student_id, parent_id, relationship)
  VALUES (v_s1, v_parent, 'father')
  ON CONFLICT DO NOTHING;

  INSERT INTO student_parents (student_id, parent_id, relationship)
  VALUES (v_s2, v_parent, 'father')
  ON CONFLICT DO NOTHING;

  -- 8) Optional demo data for attendance/health/posts (only if tables exist)
  IF to_regclass('public.attendance_records') IS NOT NULL THEN
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (v_s1, CURRENT_DATE, 'present', 'Đi học đúng giờ', v_teacher_user)
    ON CONFLICT (student_id, date) DO NOTHING;
  END IF;

  IF to_regclass('public.health_logs') IS NOT NULL THEN
    INSERT INTO health_logs (student_id, temperature, symptoms, severity, note, recorded_by)
    VALUES (v_s1, 36.7, 'Bình thường', 'normal', 'Ổn', v_teacher_user);
  END IF;

  IF to_regclass('public.posts') IS NOT NULL THEN
    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (v_teacher_user, 'class', v_class, 'activity',
            'Hôm nay lớp Lá Non học vẽ và vận động ngoài trời.');
  END IF;
END $$;
