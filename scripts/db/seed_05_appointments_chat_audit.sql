-- ============================================================================
-- IRIS Seed Data — File 5: Appointments, Chat & Audit Logs
-- Mục tiêu: 50+ appointment slots, 50+ appointments (đủ status),
--           50+ conversations, 50+ messages, 50+ audit logs
-- Thứ tự chạy: 5/5 (cần chạy sau seed_02)
-- ============================================================================

DO $$
DECLARE
  v_teacher_ids    uuid[];
  v_parent_ids     uuid[];
  v_student_ids    uuid[];
  v_class_ids      uuid[];
  v_user_ids       uuid[];
  v_school_ids     uuid[];
  v_slot_ids       uuid[] := ARRAY[]::uuid[];
  v_conv_ids       uuid[] := ARRAY[]::uuid[];
  v_slot_id        uuid;
  v_tid            uuid;
  v_pid            uuid;
  v_sid            uuid;
  v_cid            uuid;
  v_uid            uuid;
  v_scid           uuid;
  v_apt_id         uuid;
  i                int;
  j                int;
BEGIN
  SELECT array_agg(teacher_id) INTO v_teacher_ids FROM teachers ORDER BY full_name;
  SELECT array_agg(parent_id) INTO v_parent_ids FROM parents ORDER BY full_name;
  SELECT array_agg(student_id) INTO v_student_ids FROM students ORDER BY full_name;
  SELECT array_agg(class_id) INTO v_class_ids FROM classes ORDER BY school_id, name;
  SELECT array_agg(school_id) INTO v_school_ids FROM schools ORDER BY name;
  SELECT array_agg(user_id) INTO v_user_ids FROM users ORDER BY email;

  -- ======================
  -- 1) APPOINTMENT SLOTS (56 slots — mỗi teacher 2-3 slots)
  -- ======================
  IF to_regclass('public.appointment_slots') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_tid := v_teacher_ids[((i-1) % 24) + 1];
      v_cid := v_class_ids[((i-1) % 56) + 1];

      INSERT INTO appointment_slots (teacher_id, class_id, start_time, end_time, note, is_active)
      VALUES (
        v_tid, v_cid,
        date_trunc('day', NOW()) + ((i % 5 + 1) * INTERVAL '1 day') + (TIME '08:00' + (i % 8) * INTERVAL '30 minutes'),
        date_trunc('day', NOW()) + ((i % 5 + 1) * INTERVAL '1 day') + (TIME '08:00' + (i % 8) * INTERVAL '30 minutes') + INTERVAL '20 minutes',
        CASE (i % 4)
          WHEN 0 THEN 'Tư vấn về kỹ năng tự phục vụ'
          WHEN 1 THEN 'Trao đổi về chế độ dinh dưỡng'
          WHEN 2 THEN 'Đánh giá tiến bộ học tập'
          ELSE 'Họp phụ huynh định kỳ'
        END,
        true
      )
      ON CONFLICT DO NOTHING;

      SELECT slot_id INTO v_slot_id FROM appointment_slots
      WHERE teacher_id = v_tid AND class_id = v_cid
        AND start_time = date_trunc('day', NOW()) + ((i % 5 + 1) * INTERVAL '1 day') + (TIME '08:00' + (i % 8) * INTERVAL '30 minutes')
      ORDER BY created_at DESC LIMIT 1;

      v_slot_ids := array_append(v_slot_ids, v_slot_id);
    END LOOP;
  END IF;

  -- ======================
  -- 2) APPOINTMENTS (56 appointments — đủ 5 status)
  -- ======================
  IF to_regclass('public.appointments') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_slot_id := v_slot_ids[i];
      v_pid := v_parent_ids[((i-1) % 24) + 1];
      v_sid := v_student_ids[((i-1) % 56) + 1];

      CASE (i % 5)
        WHEN 0 THEN
          -- pending
          INSERT INTO appointments (slot_id, parent_id, student_id, status, note)
          VALUES (v_slot_id, v_pid, v_sid, 'pending',
                  'Phụ huynh muốn trao đổi về tình hình học tập của bé')
          ON CONFLICT DO NOTHING;
        WHEN 1 THEN
          -- confirmed
          INSERT INTO appointments (slot_id, parent_id, student_id, status, note, confirmed_at)
          VALUES (v_slot_id, v_pid, v_sid, 'confirmed',
                  'Đã xác nhận lịch hẹn', NOW() - INTERVAL '1 hour')
          ON CONFLICT DO NOTHING;
        WHEN 2 THEN
          -- cancelled
          INSERT INTO appointments (slot_id, parent_id, student_id, status, note, cancel_reason, cancelled_at)
          VALUES (v_slot_id, v_pid, v_sid, 'cancelled',
                  'Muốn trao đổi về sức khỏe', 'Phụ huynh bận đột xuất', NOW() - INTERVAL '30 minutes')
          ON CONFLICT DO NOTHING;
        WHEN 3 THEN
          -- completed
          INSERT INTO appointments (slot_id, parent_id, student_id, status, note, confirmed_at, completed_at)
          VALUES (v_slot_id, v_pid, v_sid, 'completed',
                  'Đã trao đổi xong về kỹ năng xã hội', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day')
          ON CONFLICT DO NOTHING;
        WHEN 4 THEN
          -- no_show
          INSERT INTO appointments (slot_id, parent_id, student_id, status, note, confirmed_at)
          VALUES (v_slot_id, v_pid, v_sid, 'no_show',
                  'Phụ huynh không đến theo lịch đã xác nhận', NOW() - INTERVAL '3 days')
          ON CONFLICT DO NOTHING;
      END CASE;
    END LOOP;
  END IF;

  -- ======================
  -- 3) CONVERSATIONS (50 conversations — direct + group)
  -- ======================
  IF to_regclass('public.conversations') IS NOT NULL THEN
    -- 30 direct conversations
    FOR i IN 1..30 LOOP
      INSERT INTO conversations (type, name)
      VALUES ('direct', NULL)
      ON CONFLICT DO NOTHING
      RETURNING conversation_id INTO v_cid;
      v_conv_ids := array_append(v_conv_ids, v_cid);
    END LOOP;

    -- 20 group conversations
    FOR i IN 1..20 LOOP
      INSERT INTO conversations (type, name)
      VALUES ('group',
              CASE (i % 5)
                WHEN 0 THEN 'Nhóm trao đổi lớp ' || i
                WHEN 1 THEN 'Hội phụ huynh ' || i
                WHEN 2 THEN 'Nhóm giáo viên khối ' || i
                WHEN 3 THEN 'Ban phụ huynh trường'
                ELSE 'Nhóm ngoại khóa ' || i
              END)
      ON CONFLICT DO NOTHING
      RETURNING conversation_id INTO v_cid;
      v_conv_ids := array_append(v_conv_ids, v_cid);
    END LOOP;
  END IF;

  -- ======================
  -- 4) CONVERSATION PARTICIPANTS (mỗi conversation 2-4 người)
  -- ======================
  IF to_regclass('public.conversation_participants') IS NOT NULL THEN
    FOR i IN 1..50 LOOP
      v_cid := v_conv_ids[i];

      -- Direct: 2 participants
      -- Group: 3-4 participants
      j := CASE WHEN i <= 30 THEN 2 ELSE 3 + (i % 2) END;

      FOR k IN 1..j LOOP
        v_uid := v_user_ids[((i + k - 1) % array_length(v_user_ids, 1)) + 1];
        INSERT INTO conversation_participants (conversation_id, user_id)
        VALUES (v_cid, v_uid)
        ON CONFLICT DO NOTHING;
      END LOOP;
    END LOOP;
  END IF;

  -- ======================
  -- 5) MESSAGES (50+ messages — phân bổ vào các conversations)
  -- ======================
  IF to_regclass('public.messages') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_cid := v_conv_ids[((i-1) % 50) + 1];
      v_uid := v_user_ids[((i-1) % array_length(v_user_ids, 1)) + 1];

      INSERT INTO messages (conversation_id, sender_id, content)
      VALUES (v_cid, v_uid,
              CASE (i % 6)
                WHEN 0 THEN 'Chào cô, hôm nay bé nhà em đi học có ngoan không ạ?'
                WHEN 1 THEN 'Dạ bé hôm nay ăn hết suất và ngủ trưa rất ngoan ạ.'
                WHEN 2 THEN 'Cô ơi cho em xin thực đơn tuần này với ạ.'
                WHEN 3 THEN 'Phụ huynh lưu ý ngày mai bé cần mang theo bình nước cá nhân ạ.'
                WHEN 4 THEN 'Cảm ơn cô đã thông báo, em sẽ chuẩn bị cho bé.'
                ELSE 'Chúc cô và các con một ngày vui vẻ!'
              END)
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  -- ======================
  -- 6) AUDIT LOGS (56 logs — đa dạng action, entity_type, actor_role)
  -- ======================
  IF to_regclass('public.audit_logs') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_uid := v_user_ids[((i-1) % array_length(v_user_ids, 1)) + 1];
      v_scid := v_school_ids[((i-1) % 8) + 1];

      INSERT INTO audit_logs (actor_user_id, actor_role, school_id, action, entity_type, entity_id, details, created_at)
      VALUES (
        v_uid,
        CASE (i % 4)
          WHEN 0 THEN 'SUPER_ADMIN'
          WHEN 1 THEN 'SCHOOL_ADMIN'
          WHEN 2 THEN 'TEACHER'
          ELSE 'PARENT'
        END,
        CASE WHEN i % 4 = 0 THEN NULL ELSE v_scid END,
        CASE (i % 7)
          WHEN 0 THEN 'users.login'
          WHEN 1 THEN 'attendance.create'
          WHEN 2 THEN 'posts.create'
          WHEN 3 THEN 'appointments.book'
          WHEN 4 THEN 'appointments.confirm'
          WHEN 5 THEN 'audit_logs.list'
          ELSE 'students.view'
        END,
        CASE (i % 7)
          WHEN 0 THEN 'users'
          WHEN 1 THEN 'attendance_records'
          WHEN 2 THEN 'posts'
          WHEN 3 THEN 'appointments'
          WHEN 4 THEN 'appointments'
          WHEN 5 THEN 'audit_logs'
          ELSE 'students'
        END,
        CASE WHEN i % 7 IN (1,2,3,4) THEN gen_random_uuid() ELSE NULL END,
        jsonb_build_object(
          'seed_key', 'seed_audit_' || LPAD(i::text, 3, '0'),
          'source', 'seed_05',
          'ip', '192.168.1.' || (i % 254 + 1),
          'user_agent', 'Mozilla/5.0 IrisApp'
        ),
        NOW() - (i * INTERVAL '10 minutes')
      )
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  RAISE NOTICE '✅ seed_05: 56 slots + 56 appointments + 50 conversations + 56 messages + 56 audit logs inserted';
END $$;