CREATE TABLE IF NOT EXISTS appointment_slots (
  slot_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  teacher_id   uuid NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
  class_id     uuid NOT NULL REFERENCES classes(class_id) ON DELETE CASCADE,
  start_time   timestamptz NOT NULL,
  end_time     timestamptz NOT NULL,
  note         text,
  is_active    boolean NOT NULL DEFAULT true,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_appointment_slots_teacher_time
  ON appointment_slots(teacher_id, start_time DESC);

CREATE INDEX IF NOT EXISTS idx_appointment_slots_class_time
  ON appointment_slots(class_id, start_time DESC);

CREATE TABLE IF NOT EXISTS appointments (
  appointment_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slot_id        uuid NOT NULL REFERENCES appointment_slots(slot_id) ON DELETE CASCADE,
  parent_id      uuid NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
  student_id     uuid NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
  status         varchar(20) NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed', 'no_show')),
  note           text,
  cancel_reason  text,
  confirmed_at   timestamptz,
  completed_at   timestamptz,
  cancelled_at   timestamptz,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  UNIQUE (slot_id)
);

CREATE INDEX IF NOT EXISTS idx_appointments_parent_status
  ON appointments(parent_id, status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_appointments_student_status
  ON appointments(student_id, status, created_at DESC);
