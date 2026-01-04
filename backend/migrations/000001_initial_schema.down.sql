-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_organization_users_org_id;
DROP INDEX IF EXISTS idx_organization_users_user_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_organizations_slug;

-- Drop tables
DROP TABLE IF EXISTS organization_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
