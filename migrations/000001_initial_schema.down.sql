-- File: migrations/000001_initial_schema.down.sql

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP TRIGGER IF EXISTS update_mentorship_programs_updated_at ON mentorship_programs;
DROP TRIGGER IF EXISTS update_mentorship_requests_updated_at ON mentorship_requests;
DROP TRIGGER IF EXISTS update_mentorship_sessions_updated_at ON mentorship_sessions;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in correct order due to foreign key constraints
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS session_feedback;
DROP TABLE IF EXISTS availability;
DROP TABLE IF EXISTS mentorship_sessions;
DROP TABLE IF EXISTS mentorship_requests;
DROP TABLE IF EXISTS mentorship_programs;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;
