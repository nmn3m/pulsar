-- Drop indexes
DROP INDEX IF EXISTS idx_team_invitations_unique_pending;
DROP INDEX IF EXISTS idx_team_invitations_status;
DROP INDEX IF EXISTS idx_team_invitations_token;
DROP INDEX IF EXISTS idx_team_invitations_team_id;
DROP INDEX IF EXISTS idx_team_invitations_email;

-- Drop table
DROP TABLE IF EXISTS team_invitations;
