-- Drop is_approved column from profiles if it exists
ALTER TABLE profiles DROP COLUMN IF EXISTS is_approved;

-- Drop is_admin column from users if it exists
ALTER TABLE users DROP COLUMN IF EXISTS is_admin;

-- Drop admin_settings table if it exists
DROP TABLE IF EXISTS admin_settings;
