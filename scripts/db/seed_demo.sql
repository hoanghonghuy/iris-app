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
  v_admin_user    uuid;
  v_teacher1_user uuid;
  v_teacher2_user uuid;
  v_parent1_user  uuid;
  v_parent2_user  uuid;

  -- Teachers
  v_teacher1 uuid; -- Cô Lan
  v_teacher2 uuid; -- Thầy Nam

  -- Parents
  v_parent1 uuid; -- Anh Minh
  v_parent2 uuid; -- Chị Hoa

  -- Students
  v_s1 uuid; -- Bé An
  v_s2 uuid; -- Bé Bông
  v_s3 uuid; -- Bé Na
  v_s4 uuid; -- Bé Dương
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
  WHERE school_id = v_school1 AND name = 'Lá Non' AND school_year = '2025-2026'
  ORDER BY created_at
  LIMIT 1;

  IF v_class1 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school1, 'Lá Non', '2025-2026')
    RETURNING class_id INTO v_class1;
  END IF;

  SELECT class_id INTO v_class2
  FROM classes
  WHERE school_id = v_school1 AND name = 'Lá Măng' AND school_year = '2025-2026'
  ORDER BY created_at
  LIMIT 1;

  IF v_class2 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school1, 'Lá Măng', '2025-2026')
    RETURNING class_id INTO v_class2;
  END IF;

  SELECT class_id INTO v_class3
  FROM classes
  WHERE school_id = v_school2 AND name = 'Lá Sen' AND school_year = '2025-2026'
  ORDER BY created_at
  LIMIT 1;

  IF v_class3 IS NULL THEN
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school2, 'Lá Sen', '2025-2026')
    RETURNING class_id INTO v_class3;
  END IF;

  -- ======================
  -- 3) USERS (password: 123456)
  -- bcrypt hash for "123456"
  -- ======================
  INSERT INTO users (email, password_hash, status)
  VALUES ('admin@iris.local', '$2b$12$6hU62IB6hIoxNRCXPh9Diu7QPaQkJhuMCfkjytTyphL/tX3ly93um', 'active')
  ON CONFLICT (email) DO UPDATE
    SET password_hash = EXCLUDED.password_hash, status = EXCLUDED.status
  RETURNING user_id INTO v_admin_user;

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
  SELECT v_admin_user, role_id FROM roles WHERE name = 'ADMIN'
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
  -- 7) STUDENTS
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
  -- 8) TEACHER-CLASS ASSIGNMENTS
  -- Cô Lan → Lá Non
  -- Thầy Nam → Lá Non + Lá Măng (multiple classes)
  -- ======================
  INSERT INTO teacher_classes (teacher_id, class_id)
  VALUES (v_teacher1, v_class1)
  ON CONFLICT DO NOTHING;

  INSERT INTO teacher_classes (teacher_id, class_id)
  VALUES (v_teacher2, v_class1)
  ON CONFLICT DO NOTHING;

  INSERT INTO teacher_classes (teacher_id, class_id)
  VALUES (v_teacher2, v_class2)
  ON CONFLICT DO NOTHING;

  -- ======================
  -- 9) STUDENT-PARENT RELATIONSHIPS
  -- Bé An: Anh Minh (father) + Chị Hoa (mother)
  -- Bé Bông: Anh Minh (father)
  -- Bé Na: Chị Hoa (mother)
  -- Bé Dương: không có parent (để test edge case)
  -- ======================
  INSERT INTO student_parents (student_id, parent_id, relationship)
  VALUES (v_s1, v_parent1, 'father')
  ON CONFLICT DO NOTHING;

  INSERT INTO student_parents (student_id, parent_id, relationship)
  VALUES (v_s1, v_parent2, 'mother')
  ON CONFLICT DO NOTHING;

  INSERT INTO student_parents (student_id, parent_id, relationship)
  VALUES (v_s2, v_parent1, 'father')
  ON CONFLICT DO NOTHING;

  INSERT INTO student_parents (student_id, parent_id, relationship)
  VALUES (v_s3, v_parent2, 'mother')
  ON CONFLICT DO NOTHING;

  -- ======================
  -- 10) ATTENDANCE RECORDS (đa dạng status, có check-in/out, date range)
  -- ======================
  IF to_regclass('public.attendance_records') IS NOT NULL THEN
    -- Bé An: present (hôm nay, có check-in/out)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '07:45:00')::timestamptz,
      (CURRENT_DATE + TIME '16:30:00')::timestamptz,
      'Đi học đúng giờ',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé An: late (hôm qua)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE - INTERVAL '1 day',
      'late',
      (CURRENT_DATE - INTERVAL '1 day' + TIME '08:30:00')::timestamptz,
      'Đi muộn do đưa đón',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé An: absent (2 ngày trước)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s1,
      CURRENT_DATE - INTERVAL '2 days',
      'absent',
      'Nghỉ ốm',
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Bông: excused (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
    VALUES (
      v_s2,
      CURRENT_DATE,
      'excused',
      'Nghỉ có phép - khám bác sĩ',
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Bông: present (hôm qua)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, recorded_by)
    VALUES (
      v_s2,
      CURRENT_DATE - INTERVAL '1 day',
      'present',
      (CURRENT_DATE - INTERVAL '1 day' + TIME '07:50:00')::timestamptz,
      (CURRENT_DATE - INTERVAL '1 day' + TIME '16:25:00')::timestamptz,
      v_teacher1_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;

    -- Bé Na: present (hôm nay)
    INSERT INTO attendance_records (student_id, date, status, check_in_at, recorded_by)
    VALUES (
      v_s3,
      CURRENT_DATE,
      'present',
      (CURRENT_DATE + TIME '08:00:00')::timestamptz,
      v_teacher2_user
    )
    ON CONFLICT (student_id, date) DO NOTHING;
  END IF;

  -- ======================
  -- 11) HEALTH LOGS (đa dạng severity, date range)
  -- ======================
  IF to_regclass('public.health_logs') IS NOT NULL THEN
    -- Bé An: normal (hôm nay)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s1,
      NOW(),
      36.7,
      'Bình thường, vui vẻ',
      'normal',
      'Sức khỏe tốt',
      v_teacher1_user
    );

    -- Bé An: watch (2 ngày trước - nhiệt độ hơi cao)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s1,
      NOW() - INTERVAL '2 days',
      37.3,
      'Hơi ấm, mệt mỏi',
      'watch',
      'Theo dõi sát, thông báo phụ huynh',
      v_teacher1_user
    );

    -- Bé Bông: normal (hôm qua)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s2,
      NOW() - INTERVAL '1 day',
      36.5,
      'Khỏe mạnh, năng động',
      'normal',
      'Ổn định',
      v_teacher2_user
    );

    -- Bé Na: urgent (hôm nay - sốt cao)
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s3,
      NOW(),
      38.5,
      'Sốt cao, mệt mỏi, không ăn',
      'urgent',
      'Đã gọi phụ huynh đón về ngay',
      v_teacher2_user
    );

    -- Bé Dương: normal
    INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
    VALUES (
      v_s4,
      NOW(),
      36.6,
      'Vui vẻ, ăn ngon',
      'normal',
      'Tốt',
      v_teacher2_user
    );
  END IF;

  -- ======================
  -- 12) POSTS (đa dạng type, scope: class + student, có attachments)
  -- ======================
  IF to_regclass('public.posts') IS NOT NULL THEN
    DECLARE
      v_post1 uuid;
      v_post2 uuid;
      v_post3 uuid;
      v_post4 uuid;
      v_post5 uuid;
    BEGIN
      -- Post 1: Class announcement (Lá Non)
      INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
      VALUES (
        v_teacher1_user,
        'class',
        v_class1,
        'announcement',
        '📢 Thông báo: Tuần sau lớp Lá Non sẽ tổ chức dã ngoại tại công viên. Phụ huynh vui lòng chuẩn bị mũ, áo chống nắng cho bé.'
      )
      RETURNING post_id INTO v_post1;

      -- Post 2: Class activity (Lá Non)
      INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
      VALUES (
        v_teacher2_user,
        'class',
        v_class1,
        'activity',
        '🎨 Hôm nay lớp Lá Non học vẽ và vận động ngoài trời. Các bé rất hào hứng và tham gia tích cực!'
      )
      RETURNING post_id INTO v_post2;

      -- Post 3: Student daily note (Bé An)
      INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
      VALUES (
        v_teacher1_user,
        'student',
        v_s1,
        'daily_note',
        '📝 Bé An hôm nay ăn uống ngon miệng, ngủ trưa đủ giấc. Bé tham gia hoạt động tích cực và vui vẻ!'
      )
      RETURNING post_id INTO v_post3;

      -- Post 4: Student health note (Bé Na - urgent health)
      INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
      VALUES (
        v_teacher2_user,
        'student',
        v_s3,
        'health_note',
        '🌡️ Bé Na có dấu hiệu sốt cao (38.5°C). Cô đã đo nhiệt độ và thông báo phụ huynh đón bé về ngay. Phụ huynh vui lòng theo dõi sức khỏe bé.'
      )
      RETURNING post_id INTO v_post4;

      -- Post 5: Class announcement (Lá Măng)
      INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
      VALUES (
        v_teacher2_user,
        'class',
        v_class2,
        'announcement',
        '📚 Nhắc nhở: Thứ 6 tuần này là ngày chia sẻ đồ chơi. Các bé nhớ mang theo 1 món đồ chơi yêu thích để chia sẻ với bạn nhé!'
      )
      RETURNING post_id INTO v_post5;

      -- ======================
      -- 13) POST ATTACHMENTS
      -- ======================
      IF to_regclass('public.post_attachments') IS NOT NULL THEN
        -- Attachments cho post 2 (activity)
        INSERT INTO post_attachments (post_id, url, mime_type)
        VALUES
          (v_post2, 'https://images.unsplash.com/photo-1587654780291-39c9404d746b?w=800', 'image/jpeg'),
          (v_post2, 'https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?w=800', 'image/jpeg');

        -- Attachment cho post 3 (daily note)
        INSERT INTO post_attachments (post_id, url, mime_type)
        VALUES
          (v_post3, 'https://images.unsplash.com/photo-1476224203421-9ac39bcb3327?w=800', 'image/jpeg');

        -- Attachment cho post 4 (health note)
        INSERT INTO post_attachments (post_id, url, mime_type)
        VALUES
          (v_post4, 'https://via.placeholder.com/800x600/4CAF50/FFFFFF?text=Temperature+Chart', 'image/png');
      END IF;
    END;
  END IF;

END $$;
