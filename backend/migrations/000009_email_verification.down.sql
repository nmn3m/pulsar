-- Drop email_verifications table
DROP TABLE IF EXISTS email_verifications;

-- Remove email_verified column from users
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;
