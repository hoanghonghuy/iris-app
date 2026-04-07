ALTER TABLE appointments
  DROP CONSTRAINT IF EXISTS appointments_slot_id_key;

CREATE INDEX IF NOT EXISTS idx_appointments_slot_active
  ON appointments(slot_id)
  WHERE status IN ('pending', 'confirmed', 'completed', 'no_show');
