-- Attendance records
CREATE TABLE IF NOT EXISTS attendance_records (
  attendance_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  student_id    uuid NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
  date          date NOT NULL,
  status        varchar(20) NOT NULL DEFAULT 'present'
    CHECK (status IN ('present','absent')),
  check_in_at   timestamptz,
  check_out_at  timestamptz,
  note          text,
  recorded_by   uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  UNIQUE (student_id, date)
);

CREATE INDEX IF NOT EXISTS idx_attendance_student_date
  ON attendance_records (student_id, date);


-- Health logs
CREATE TABLE IF NOT EXISTS health_logs (
  health_log_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  student_id    uuid NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
  recorded_at   timestamptz NOT NULL DEFAULT now(),
  temperature   numeric(4,1),
  symptoms      text,
  severity      varchar(20)
    CHECK (severity IN ('normal','watch','urgent')),
  note          text,
  recorded_by   uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_health_student_recorded_at
  ON health_logs (student_id, recorded_at DESC);


-- Posts (Newsfeed)
CREATE TABLE IF NOT EXISTS posts (
  post_id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  author_user_id  uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  scope_type      varchar(20) NOT NULL DEFAULT 'class'
    CHECK (scope_type IN ('school','class','student')),
  school_id       uuid REFERENCES schools(school_id) ON DELETE CASCADE,
  class_id        uuid REFERENCES classes(class_id) ON DELETE CASCADE,
  student_id      uuid REFERENCES students(student_id) ON DELETE CASCADE,
  type            varchar(30) NOT NULL DEFAULT 'activity'
    CHECK (type IN ('announcement','activity','daily_note','health_note')),
  content         text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),

  -- Enforce scope columns based on scope_type (đỡ lộ dữ liệu sai scope)
  CHECK (
    (scope_type = 'school'  AND school_id IS NOT NULL AND class_id IS NULL AND student_id IS NULL) OR
    (scope_type = 'class'   AND class_id  IS NOT NULL AND school_id IS NULL AND student_id IS NULL) OR
    (scope_type = 'student' AND student_id IS NOT NULL AND school_id IS NULL AND class_id IS NULL)
  )
);

CREATE INDEX IF NOT EXISTS idx_posts_scope_class_created
  ON posts (scope_type, class_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_posts_scope_student_created
  ON posts (scope_type, student_id, created_at DESC);


-- Post attachments
CREATE TABLE IF NOT EXISTS post_attachments (
  attachment_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id       uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
  url           text NOT NULL,
  mime_type     varchar(100),
  created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_attachments_post_id
  ON post_attachments (post_id);
