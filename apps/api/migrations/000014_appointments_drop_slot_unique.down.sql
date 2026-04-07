DROP INDEX IF EXISTS idx_appointments_slot_active;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM appointments
    GROUP BY slot_id
    HAVING COUNT(*) > 1
  ) THEN
    RAISE EXCEPTION 'cannot restore appointments_slot_id_key: duplicate slot_id values exist';
  END IF;
END $$;

ALTER TABLE appointments
  ADD CONSTRAINT appointments_slot_id_key UNIQUE (slot_id);
