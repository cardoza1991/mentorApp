-- Down migration:
-- File: migrations/000002_update_constraints.down.sql

ALTER TABLE availability
    DROP CONSTRAINT IF EXISTS availability_time_check;

ALTER TABLE availability
    ALTER COLUMN start_time TYPE TIME USING start_time::time,
    ALTER COLUMN end_time TYPE TIME USING end_time::time;

ALTER TABLE availability
    DROP CONSTRAINT IF EXISTS availability_mentor_timeslot_unique;
