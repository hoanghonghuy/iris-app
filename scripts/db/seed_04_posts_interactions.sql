-- ============================================================================
-- IRIS Seed Data — File 4: Posts & Interactions
-- Mục tiêu: 50+ posts (đủ scope_type + type), 50+ comments,
--           50+ interactions (like/share), 50+ edit history, 50+ attachments
-- Thứ tự chạy: 4/5 (cần chạy sau seed_02)
-- ============================================================================

DO $$
DECLARE
  v_student_ids    uuid[];
  v_class_ids      uuid[];
  v_school_ids     uuid[];
  v_teacher_users  uuid[];
  v_parent_users   uuid[];
  v_post_ids       uuid[] := ARRAY[]::uuid[];
  v_pid            uuid;
  v_uid            uuid;
  v_cid            uuid;
  v_sid            uuid;
  v_scid           uuid;
  i                int;
  j                int;
BEGIN
  SELECT array_agg(student_id) INTO v_student_ids FROM students ORDER BY full_name;
  SELECT array_agg(class_id) INTO v_class_ids FROM classes ORDER BY school_id, name;
  SELECT array_agg(school_id) INTO v_school_ids FROM schools ORDER BY name;

  SELECT array_agg(u.user_id) INTO v_teacher_users
  FROM users u JOIN user_roles ur ON u.user_id = ur.user_id
  JOIN roles r ON ur.role_id = r.role_id WHERE r.name = 'TEACHER' ORDER BY u.email;

  SELECT array_agg(u.user_id) INTO v_parent_users
  FROM users u JOIN user_roles ur ON u.user_id = ur.user_id
  JOIN roles r ON ur.role_id = r.role_id WHERE r.name = 'PARENT' ORDER BY u.email;

  -- ======================
  -- 1) POSTS (56 posts — đủ scope_type × type)
  -- ======================

  -- scope_type = 'school' (8 posts — mỗi trường 1 announcement)
  FOR i IN 1..8 LOOP
    v_uid := v_teacher_users[i];
    v_scid := v_school_ids[i];

    INSERT INTO posts (author_user_id, scope_type, school_id, type, content)
    VALUES (v_uid, 'school', v_scid, 'announcement',
            CASE (i % 4)
              WHEN 0 THEN 'Thông báo toàn trường: Lịch nghỉ lễ và các hoạt động ngoại khóa tháng tới. Phụ huynh vui lòng theo dõi trên ứng dụng.'
              WHEN 1 THEN 'Thông báo: Nhà trường tổ chức khám sức khỏe định kỳ cho tất cả các bé vào tuần sau.'
              WHEN 2 THEN 'Thông báo: Lịch họp phụ huynh định kỳ tháng này, đề nghị quý phụ huynh sắp xếp thời gian tham dự.'
              ELSE 'Thông báo: Thay đổi giờ đón trẻ buổi chiều, phụ huynh vui lòng cập nhật thông tin.'
            END)
    ON CONFLICT DO NOTHING;

    SELECT post_id INTO v_pid FROM posts
    WHERE author_user_id = v_uid AND scope_type = 'school' AND school_id = v_scid
    ORDER BY created_at DESC LIMIT 1;

    v_post_ids := array_append(v_post_ids, v_pid);
  END LOOP;

  -- scope_type = 'class' (24 posts — mỗi loại type 6 posts)
  FOR i IN 1..24 LOOP
    v_uid := v_teacher_users[((i-1) % 24) + 1];
    v_cid := v_class_ids[((i-1) % 56) + 1];

    INSERT INTO posts (author_user_id, scope_type, class_id, type, content)
    VALUES (v_uid, 'class', v_cid,
            CASE (i % 4)
              WHEN 0 THEN 'announcement'
              WHEN 1 THEN 'activity'
              WHEN 2 THEN 'health_note'
              ELSE 'announcement'
            END,
            CASE (i % 4)
              WHEN 0 THEN 'Thông báo lớp: Thứ Sáu tuần này các con tham gia hoạt động trải nghiệm tại sân trường. Phụ huynh chuẩn bị mũ và bình nước.'
              WHEN 1 THEN 'Hoạt động hôm nay: Các con tập vẽ tranh chủ đề gia đình, rất hào hứng và sáng tạo. Ảnh đã được cập nhật trong album lớp.'
              WHEN 2 THEN 'Lưu ý sức khỏe: Hôm nay lớp có 1 bé ho nhẹ, giáo viên đã theo dõi và nhắc uống nước ấm thường xuyên.'
              ELSE 'Thông báo: Tuần sau lớp bắt đầu chủ đề mới "Thế giới động vật", phụ huynh có thể cho bé mang sách/tranh về động vật.'
            END)
    ON CONFLICT DO NOTHING;

    SELECT post_id INTO v_pid FROM posts
    WHERE author_user_id = v_uid AND scope_type = 'class' AND class_id = v_cid
    ORDER BY created_at DESC LIMIT 1;

    v_post_ids := array_append(v_post_ids, v_pid);
  END LOOP;

  -- scope_type = 'student' (24 posts — daily_note + health_note)
  FOR i IN 1..24 LOOP
    v_uid := v_teacher_users[((i-1) % 24) + 1];
    v_sid := v_student_ids[((i-1) % 56) + 1];

    INSERT INTO posts (author_user_id, scope_type, student_id, type, content)
    VALUES (v_uid, 'student', v_sid,
            CASE WHEN i % 2 = 0 THEN 'daily_note' ELSE 'health_note' END,
            CASE WHEN i % 2 = 0 THEN
              'Nhật ký hôm nay: Bé ăn hết suất, ngủ trưa đúng giờ, tham gia tốt hoạt động nhóm. Gia đình duy trì giờ ngủ tối trước 21h.'
            ELSE
              'Ghi chú sức khỏe: Bé có biểu hiện ho nhẹ sau giờ ngủ trưa, giáo viên đã theo dõi và cho uống nước ấm. Phụ huynh theo dõi thêm tại nhà.'
            END)
    ON CONFLICT DO NOTHING;

    SELECT post_id INTO v_pid FROM posts
    WHERE author_user_id = v_uid AND scope_type = 'student' AND student_id = v_sid
    ORDER BY created_at DESC LIMIT 1;

    v_post_ids := array_append(v_post_ids, v_pid);
  END LOOP;

  -- ======================
  -- 2) POST ATTACHMENTS (56 attachments — mỗi post 1 ảnh)
  -- ======================
  IF to_regclass('public.post_attachments') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      INSERT INTO post_attachments (post_id, url, mime_type)
      VALUES (v_post_ids[i],
              'https://picsum.photos/seed/iris' || i || '/800/600',
              'image/jpeg')
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  -- ======================
  -- 3) POST COMMENTS (56 comments — từ parents + teachers)
  -- ======================
  IF to_regclass('public.post_comments') IS NOT NULL THEN
    FOR i IN 1..56 LOOP
      v_pid := v_post_ids[i];
      v_uid := CASE WHEN i % 2 = 0
                THEN v_parent_users[((i-1) % 24) + 1]
                ELSE v_teacher_users[((i-1) % 24) + 1]
              END;

      INSERT INTO post_comments (post_id, author_user_id, content)
      VALUES (v_pid, v_uid,
              CASE (i % 5)
                WHEN 0 THEN 'Cảm ơn cô giáo đã cập nhật thông tin ạ!'
                WHEN 1 THEN 'Bé ở nhà cũng rất thích hoạt động này.'
                WHEN 2 THEN 'Phụ huynh đã nhận được thông báo, cảm ơn nhà trường.'
                WHEN 3 THEN 'Cô giáo ơi cho em hỏi thêm về lịch tuần sau với ạ.'
                ELSE 'Rất hữu ích, cảm ơn cô!'
              END)
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  -- ======================
  -- 4) POST INTERACTIONS (56 likes + 20 shares = 76 records)
  -- ======================
  IF to_regclass('public.post_interactions') IS NOT NULL THEN
    -- Likes: mỗi post có 1-2 likes
    FOR i IN 1..56 LOOP
      v_pid := v_post_ids[i];
      v_uid := v_parent_users[((i-1) % 24) + 1];

      INSERT INTO post_interactions (post_id, user_id, action_type)
      VALUES (v_pid, v_uid, 'like')
      ON CONFLICT DO NOTHING;

      -- Thêm 1 like nữa từ teacher cho 1/2 số post
      IF i % 2 = 0 THEN
        v_uid := v_teacher_users[((i-1) % 24) + 1];
        INSERT INTO post_interactions (post_id, user_id, action_type)
        VALUES (v_pid, v_uid, 'like')
        ON CONFLICT DO NOTHING;
      END IF;
    END LOOP;

    -- Shares: 20 shares
    FOR i IN 1..20 LOOP
      v_pid := v_post_ids[i];
      v_uid := v_parent_users[((i-1) % 24) + 1];

      INSERT INTO post_interactions (post_id, user_id, action_type)
      VALUES (v_pid, v_uid, 'share')
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  -- ======================
  -- 5) POST EDIT HISTORY (50 records)
  -- ======================
  IF to_regclass('public.post_edit_history') IS NOT NULL THEN
    FOR i IN 1..50 LOOP
      v_pid := v_post_ids[i];
      v_uid := v_teacher_users[((i-1) % 24) + 1];

      INSERT INTO post_edit_history (post_id, old_content, new_content, edited_by)
      VALUES (v_pid,
              'Nội dung cũ trước khi chỉnh sửa lần ' || i,
              'Nội dung mới sau khi cập nhật lần ' || i || ' — đã bổ sung thêm chi tiết.',
              v_uid)
      ON CONFLICT DO NOTHING;
    END LOOP;
  END IF;

  RAISE NOTICE '✅ seed_04: 56 posts + 56 attachments + 56 comments + 76 interactions + 50 edit history inserted';
END $$;