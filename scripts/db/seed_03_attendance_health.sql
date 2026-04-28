-- ============================================================================
-- IRIS Seed Data — File 3: Attendance & Health
-- Mục tiêu: 50+ attendance records, 50+ health logs, 50+ attendance change logs
--           Mỗi status/severity/change_type đều có đại diện
-- Thứ tự chạy: 3/5 (cần chạy sau seed_02)
-- ============================================================================

DO $$
DECLARE
  v_student_ids  uuid[];
  v_teacher_users uuid[];
  v_sid          uuid;
  v_tuid         uuid;
  v_aid          uuid;
  v_date         date;
  i              int;
  j              int;
BEGIN
  -- Lấy tất cả student IDs và teacher user IDs
  SELECT array_agg(student_id) INTO v_student_ids FROM students ORDER BY full_name;
  SELECT array_agg(u.user_id) INTO v_teacher_users
  FROM users u
  JOIN user_roles ur ON u.user_id = ur.user_id
  JOIN roles r ON ur.role_id = r.role_id
  WHERE r.name = 'TEACHER'
  ORDER BY u.email;

  -- ======================
  -- 1) ATTENDANCE RECORDS (56 students × 3 days = 168 records, đủ 4 status)
  -- ======================
  IF to_regclass('public.attendance_records') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_sid := v_student_ids[i];
      -- Lấy teacher user từ school của student
      SELECT u.user_id INTO v_tuid
      FROM users u
      JOIN teachers t ON t.user_id = u.user_id
      JOIN students s ON s.school_id = t.school_id AND s.student_id = v_sid
      LIMIT 1;

      -- Day 1: Hôm nay — đa dạng status
      v_date := CURRENT_DATE;
      CASE (i % 4)
        WHEN 0 THEN
          -- present + check-in/out
          INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
          VALUES (v_sid, v_date, 'present',
                  (v_date + TIME '08:' || LPAD(((i*7)%60)::text,2,'0') || ':00')::timestamptz,
                  (v_date + TIME '16:' || LPAD(((i*11)%60)::text,2,'0') || ':00')::timestamptz,
                  CASE (i % 3)
                    WHEN 0 THEN 'Đi học đúng giờ, tham gia tốt hoạt động nhóm'
                    WHEN 1 THEN 'Ăn trưa đủ suất, ngủ trưa ngoan'
                    ELSE 'Hoàn thành tốt bài tập trong lớp'
                  END,
                  v_tuid)
          ON CONFLICT (student_id, date) DO NOTHING;
        WHEN 1 THEN
          -- late
          INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
          VALUES (v_sid, v_date, 'late',
                  (v_date + TIME '09:' || LPAD(((i*13)%60)::text,2,'0') || ':00')::timestamptz,
                  CASE (i % 3)
                    WHEN 0 THEN 'Đi học muộn do tắc đường'
                    WHEN 1 THEN 'Đến lớp muộn vì trời mưa'
                    ELSE 'Phụ huynh báo đến muộn do việc gia đình'
                  END,
                  v_tuid)
          ON CONFLICT (student_id, date) DO NOTHING;
        WHEN 2 THEN
          -- absent
          INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
          VALUES (v_sid, v_date, 'absent',
                  CASE (i % 3)
                    WHEN 0 THEN 'Nghỉ học không phép'
                    WHEN 1 THEN 'Nghỉ do ốm, phụ huynh chưa kịp báo'
                    ELSE 'Vắng mặt chưa rõ lý do, giáo viên sẽ liên hệ'
                  END,
                  v_tuid)
          ON CONFLICT (student_id, date) DO NOTHING;
        WHEN 3 THEN
          -- excused
          INSERT INTO attendance_records (student_id, date, status, note, recorded_by)
          VALUES (v_sid, v_date, 'excused',
                  CASE (i % 3)
                    WHEN 0 THEN 'Nghỉ ốm có xác nhận của phụ huynh'
                    WHEN 1 THEN 'Nghỉ khám sức khỏe định kỳ'
                    ELSE 'Nghỉ việc gia đình, đã báo trước'
                  END,
                  v_tuid)
          ON CONFLICT (student_id, date) DO NOTHING;
      END CASE;

      -- Day 2: Hôm qua — present hoặc late
      v_date := CURRENT_DATE - 1;
      INSERT INTO attendance_records (student_id, date, status, check_in_at, check_out_at, note, recorded_by)
      VALUES (v_sid, v_date,
              CASE WHEN i % 3 = 0 THEN 'late' ELSE 'present' END,
              (v_date + TIME '08:' || LPAD(((i*5+10)%60)::text,2,'0') || ':00')::timestamptz,
              (v_date + TIME '16:' || LPAD(((i*7+15)%60)::text,2,'0') || ':00')::timestamptz,
              '',
              v_tuid)
      ON CONFLICT (student_id, date) DO NOTHING;

      -- Day 3: 2 ngày trước — present
      v_date := CURRENT_DATE - 2;
      INSERT INTO attendance_records (student_id, date, status, check_in_at, note, recorded_by)
      VALUES (v_sid, v_date, 'present',
              (v_date + TIME '08:' || LPAD(((i*3+20)%60)::text,2,'0') || ':00')::timestamptz,
              '',
              v_tuid)
      ON CONFLICT (student_id, date) DO NOTHING;
    END LOOP;
  END IF;

  -- ======================
  -- 2) HEALTH LOGS (56 students × 2 logs = 112 records, đủ 3 severity)
  -- ======================
  IF to_regclass('public.health_logs') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_sid := v_student_ids[i];
      SELECT u.user_id INTO v_tuid
      FROM users u
      JOIN teachers t ON t.user_id = u.user_id
      JOIN students s ON s.school_id = t.school_id AND s.student_id = v_sid
      LIMIT 1;

      -- Log 1: Hôm nay
      CASE (i % 3)
        WHEN 0 THEN
          -- normal
          INSERT INTO health_logs (student_id, recorded_at, temperature, severity, note, recorded_by)
          VALUES (v_sid, NOW() - (i * INTERVAL '5 minutes'), 36.0 + (i % 10) * 0.1, 'normal',
                  CASE (i % 3)
                    WHEN 0 THEN 'Sức khỏe bình thường, ăn ngủ tốt'
                    WHEN 1 THEN 'Không có dấu hiệu bất thường'
                    ELSE 'Thể trạng ổn định'
                  END,
                  v_tuid)
          ON CONFLICT DO NOTHING;
        WHEN 1 THEN
          -- watch
          INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
          VALUES (v_sid, NOW() - (i * INTERVAL '5 minutes'), 37.0 + (i % 10) * 0.1,
                  CASE (i % 3)
                    WHEN 0 THEN 'Ho nhẹ, sổ mũi'
                    WHEN 1 THEN 'Hắt hơi, hơi mệt'
                    ELSE 'Đau bụng nhẹ sau ăn'
                  END,
                  'watch',
                  'Đang theo dõi thêm, chưa cần can thiệp y tế',
                  v_tuid)
          ON CONFLICT DO NOTHING;
        WHEN 2 THEN
          -- urgent
          INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
          VALUES (v_sid, NOW() - (i * INTERVAL '5 minutes'), 38.5 + (i % 10) * 0.1,
                  CASE (i % 3)
                    WHEN 0 THEN 'Sốt cao, mệt, đau đầu'
                    WHEN 1 THEN 'Nôn nhiều, đau bụng dữ dội'
                    ELSE 'Phát ban toàn thân, ngứa'
                  END,
                  'urgent',
                  'Đã gọi phụ huynh đón về, đề nghị đưa đi khám',
                  v_tuid)
          ON CONFLICT DO NOTHING;
      END CASE;

      -- Log 2: 1-7 ngày trước (normal hoặc watch)
      INSERT INTO health_logs (student_id, recorded_at, temperature, symptoms, severity, note, recorded_by)
      VALUES (v_sid, NOW() - ((i % 7 + 1) * INTERVAL '1 day'),
              36.0 + (i % 15) * 0.1,
              CASE WHEN i % 5 = 0 THEN 'Ho khan' ELSE '' END,
              CASE WHEN i % 5 = 0 THEN 'watch' ELSE 'normal' END,
              '',
              v_tuid)
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  -- ======================
  -- 3) ATTENDANCE CHANGE LOGS (50+ records, đủ 3 change_type)
  -- ======================
  IF to_regclass('public.attendance_change_logs') IS NOT NULL THEN
    -- Lấy 50 attendance records đầu tiên để tạo change logs
    FOR i IN 1..50 LOOP
      -- Lấy 1 attendance record
      SELECT a.attendance_id, a.student_id, a.date, a.status, a.note
      INTO v_aid, v_sid, v_date, NULL, NULL
      FROM attendance_records a
      ORDER BY a.date DESC, a.student_id
      OFFSET (i - 1) LIMIT 1;

      IF v_aid IS NOT NULL THEN
        SELECT u.user_id INTO v_tuid
        FROM users u
        JOIN teachers t ON t.user_id = u.user_id
        JOIN students s ON s.school_id = t.school_id AND s.student_id = v_sid
        LIMIT 1;

        CASE (i % 3)
          WHEN 0 THEN
            -- create
            INSERT INTO attendance_change_logs (attendance_id, student_id, date, change_type, new_status, new_note, changed_by)
            VALUES (v_aid, v_sid, v_date, 'create', 'present', 'Điểm danh lần đầu', v_tuid)
            ON CONFLICT DO NOTHING;
          WHEN 1 THEN
            -- update
            INSERT INTO attendance_change_logs (attendance_id, student_id, date, change_type, old_status, new_status, old_note, new_note, changed_by)
            VALUES (v_aid, v_sid, v_date, 'update', 'absent', 'excused', 'Nghỉ không phép', 'Đã xác nhận nghỉ ốm với phụ huynh', v_tuid)
            ON CONFLICT DO NOTHING;
          WHEN 2 THEN
            -- delete (soft — new_status IS NULL)
            INSERT INTO attendance_change_logs (attendance_id, student_id, date, change_type, old_status, new_status, old_note, changed_by)
            VALUES (v_aid, v_sid, v_date, 'delete', 'present', NULL, 'Xóa do điểm danh nhầm', v_tuid)
            ON CONFLICT DO NOTHING;
        END CASE;
      END IF;
    END LOOP;
  END IF;

  RAISE NOTICE '✅ seed_03: attendance + health + change logs inserted';
END $$;