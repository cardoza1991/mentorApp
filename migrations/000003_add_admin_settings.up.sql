-- Add admin settings table
CREATE TABLE IF NOT EXISTS admin_settings (
    id SERIAL PRIMARY KEY,
    settings_key VARCHAR(50) UNIQUE NOT NULL,
    settings_value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add initial settings
INSERT INTO admin_settings (settings_key, settings_value) 
VALUES ('first_admin_email', ''),
       ('registration_enabled', 'false'),
       ('require_email_verification', 'true'),
       ('require_profile_approval', 'true');

-- Add is_admin column to users if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT false;

-- Add is_approved column to profiles if it doesn't exist
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS is_approved BOOLEAN DEFAULT false NOT NULL;
