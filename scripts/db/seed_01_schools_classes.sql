-- ============================================================================
-- IRIS Seed Data — File 1: Schools & Classes
-- Mục tiêu: 5+ schools, 50+ classes phân bổ đều
-- Thứ tự chạy: 1/5
-- ============================================================================

DO $$
DECLARE
  -- School IDs
  v_school_hanoi_1    uuid;
  v_school_hanoi_2    uuid;
  v_school_hcm_1      uuid;
  v_school_hcm_2      uuid;
  v_school_danang_1   uuid;
  v_school_danang_2   uuid;
  v_school_haiphong   uuid;
  v_school_cantho     uuid;

  -- Class holder
  v_class uuid;
BEGIN
  -- ======================
  -- 1) SCHOOLS (8 trường — 3 miền Bắc/Trung/Nam)
  -- ======================

  -- Miền Bắc
  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Bình Minh Cầu Giấy', 'Số 18 ngõ 165 Xuân Thủy, Dịch Vọng Hậu, Cầu Giấy, Hà Nội')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_hanoi_1;

  IF v_school_hanoi_1 IS NULL THEN
    SELECT school_id INTO v_school_hanoi_1 FROM schools WHERE name = 'Trường Mầm non Bình Minh Cầu Giấy' LIMIT 1;
  END IF;

  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Tuổi Thần Tiên Đống Đa', 'Số 42 ngõ 88 Láng Hạ, Đống Đa, Hà Nội')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_hanoi_2;

  IF v_school_hanoi_2 IS NULL THEN
    SELECT school_id INTO v_school_hanoi_2 FROM schools WHERE name = 'Trường Mầm non Tuổi Thần Tiên Đống Đa' LIMIT 1;
  END IF;

  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Hướng Dương Hải Phòng', 'Số 5 Lạch Tray, Ngô Quyền, Hải Phòng')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_haiphong;

  IF v_school_haiphong IS NULL THEN
    SELECT school_id INTO v_school_haiphong FROM schools WHERE name = 'Trường Mầm non Hướng Dương Hải Phòng' LIMIT 1;
  END IF;

  -- Miền Trung
  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Ánh Dương Hải Châu', 'Số 25 Lê Đình Dương, Phước Ninh, Hải Châu, Đà Nẵng')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_danang_1;

  IF v_school_danang_1 IS NULL THEN
    SELECT school_id INTO v_school_danang_1 FROM schools WHERE name = 'Trường Mầm non Ánh Dương Hải Châu' LIMIT 1;
  END IF;

  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Ngọc Lan Sơn Trà', 'Số 120 Ngô Quyền, Sơn Trà, Đà Nẵng')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_danang_2;

  IF v_school_danang_2 IS NULL THEN
    SELECT school_id INTO v_school_danang_2 FROM schools WHERE name = 'Trường Mầm non Ngọc Lan Sơn Trà' LIMIT 1;
  END IF;

  -- Miền Nam
  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Hoa Mai Quận 7', 'Số 12 đường số 3, KDC Him Lam, Tân Hưng, Quận 7, TP. Hồ Chí Minh')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_hcm_1;

  IF v_school_hcm_1 IS NULL THEN
    SELECT school_id INTO v_school_hcm_1 FROM schools WHERE name = 'Trường Mầm non Hoa Mai Quận 7' LIMIT 1;
  END IF;

  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Vườn Xanh Thủ Đức', 'Số 88 Võ Văn Ngân, Linh Chiểu, Thủ Đức, TP. Hồ Chí Minh')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_hcm_2;

  IF v_school_hcm_2 IS NULL THEN
    SELECT school_id INTO v_school_hcm_2 FROM schools WHERE name = 'Trường Mầm non Vườn Xanh Thủ Đức' LIMIT 1;
  END IF;

  INSERT INTO schools (name, address)
  VALUES ('Trường Mầm non Cánh Én Cần Thơ', 'Số 30 đường 30/4, Ninh Kiều, Cần Thơ')
  ON CONFLICT DO NOTHING
  RETURNING school_id INTO v_school_cantho;

  IF v_school_cantho IS NULL THEN
    SELECT school_id INTO v_school_cantho FROM schools WHERE name = 'Trường Mầm non Cánh Én Cần Thơ' LIMIT 1;
  END IF;

  -- ======================
  -- 2) CLASSES (56 lớp — mỗi trường 7 lớp, phân bổ 3 khối tuổi)
  -- ======================
  -- Mỗi trường: 2 lớp Nhà trẻ (18-24 tháng), 3 lớp Mầm/Chồi (3-4 tuổi), 2 lớp Lá (5 tuổi)

  -- Hanoi 1: Bình Minh Cầu Giấy
  FOREACH v_class IN ARRAY ARRAY[
    -- Nhà trẻ
    'Bé Ngoan 1', 'Bé Ngoan 2',
    -- Mầm
    'Mầm Vàng', 'Mầm Xanh', 'Mầm Đỏ',
    -- Lá
    'Lá Non', 'Lá Măng'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_hanoi_1, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- Hanoi 2: Tuổi Thần Tiên Đống Đa
  FOREACH v_class IN ARRAY ARRAY[
    'Gấu Con 1', 'Gấu Con 2',
    'Chồi Xanh', 'Chồi Vàng', 'Chồi Đỏ',
    'Lá Phong', 'Lá Sen'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_hanoi_2, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- Hai Phong
  FOREACH v_class IN ARRAY ARRAY[
    'Cún Con 1', 'Cún Con 2',
    'Mầm Hồng', 'Mầm Tím', 'Mầm Cam',
    'Lá Cọ', 'Lá Tre'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_haiphong, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- Da Nang 1: Ánh Dương Hải Châu
  FOREACH v_class IN ARRAY ARRAY[
    'Sóc Con 1', 'Sóc Con 2',
    'Chồi Biển', 'Chồi Cát', 'Chồi Nắng',
    'Lá Dừa', 'Lá Sóng'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_danang_1, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- Da Nang 2: Ngọc Lan Sơn Trà
  FOREACH v_class IN ARRAY ARRAY[
    'Thỏ Ngọc 1', 'Thỏ Ngọc 2',
    'Mầm Sơn', 'Mầm Trà', 'Mầm Hương',
    'Lá Ngọc', 'Lá Lan'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_danang_2, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- HCM 1: Hoa Mai Quận 7
  FOREACH v_class IN ARRAY ARRAY[
    'Bướm Vàng 1', 'Bướm Vàng 2',
    'Chồi Mai', 'Chồi Đào', 'Chồi Cúc',
    'Lá Mai', 'Lá Trúc'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_hcm_1, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- HCM 2: Vườn Xanh Thủ Đức
  FOREACH v_class IN ARRAY ARRAY[
    'Chim Non 1', 'Chim Non 2',
    'Mầm Lộc', 'Mầm Chồi', 'Mầm Hoa',
    'Lá Xanh', 'Lá Vàng'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_hcm_2, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  -- Can Tho
  FOREACH v_class IN ARRAY ARRAY[
    'Cá Nhỏ 1', 'Cá Nhỏ 2',
    'Chồi Sông', 'Chồi Nước', 'Chồi Mưa',
    'Lá Én', 'Lá Gió'
  ] LOOP
    INSERT INTO classes (school_id, name, school_year)
    VALUES (v_school_cantho, v_class, '2025-2026')
    ON CONFLICT DO NOTHING;
  END LOOP;

  RAISE NOTICE '✅ seed_01: 8 schools + 56 classes inserted';
END $$;