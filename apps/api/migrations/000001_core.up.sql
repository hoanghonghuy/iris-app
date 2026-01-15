CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS schools (
  school_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name        varchar(255) NOT NULL,
  address     text,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
  user_id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email         varchar(255) NOT NULL UNIQUE,
  password_hash text NOT NULL,
  status        varchar(20) NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'locked')),
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS roles (
  role_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name    varchar(30) NOT NULL UNIQUE
    CHECK (name IN ('ADMIN', 'TEACHER', 'PARENT'))
);

INSERT INTO roles (name) VALUES ('ADMIN'), ('TEACHER'), ('PARENT')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS user_roles (
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  role_id uuid NOT NULL REFERENCES roles(role_id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS classes (
  class_id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id   uuid NOT NULL REFERENCES schools(school_id) ON DELETE CASCADE,
  name        varchar(100) NOT NULL,
  school_year varchar(20) NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS students (
  student_id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id        uuid NOT NULL REFERENCES schools(school_id) ON DELETE CASCADE,
  current_class_id uuid REFERENCES classes(class_id) ON DELETE SET NULL,
  full_name        varchar(255) NOT NULL,
  dob              date,
  gender           varchar(10) CHECK (gender IN ('male','female','other')),
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS teachers (
  teacher_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
  school_id  uuid NOT NULL REFERENCES schools(school_id) ON DELETE CASCADE,
  full_name  varchar(255) NOT NULL,
  phone      varchar(30),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS parents (
  parent_id  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
  school_id  uuid NOT NULL REFERENCES schools(school_id) ON DELETE CASCADE,
  full_name  varchar(255) NOT NULL,
  phone      varchar(30),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS teacher_classes (
  teacher_id uuid NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
  class_id   uuid NOT NULL REFERENCES classes(class_id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (teacher_id, class_id)
);

CREATE TABLE IF NOT EXISTS student_parents (
  student_id   uuid NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
  parent_id    uuid NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
  relationship varchar(50),
  created_at   timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (student_id, parent_id)
);
