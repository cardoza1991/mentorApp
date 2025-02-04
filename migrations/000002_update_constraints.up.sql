-- File: migrations/000002_update_constraints.up.sql

-- Add unique constraint to availability to prevent overlapping slots
ALTER TABLE availability 
    ADD CONSTRAINT availability_mentor_timeslot_unique 
    UNIQUE (mentor_id, day_of_week);

-- Convert time columns to TIMESTAMP for better timezone handling
ALTER TABLE availability 
    ALTER COLUMN start_time TYPE TIMESTAMP USING date_trunc('day', CURRENT_DATE) + start_time,
    ALTER COLUMN end_time TYPE TIMESTAMP USING date_trunc('day', CURRENT_DATE) + end_time;

-- Add constraint to ensure end_time is after start_time
ALTER TABLE availability
    ADD CONSTRAINT availability_time_check
    CHECK (end_time > start_time);

