-- Drop triggers
DROP TRIGGER IF EXISTS update_alerts_updated_at ON alerts;
DROP TRIGGER IF EXISTS update_teams_updated_at ON teams;

-- Drop indexes for alerts
DROP INDEX IF EXISTS idx_alerts_source;
DROP INDEX IF EXISTS idx_alerts_status;
DROP INDEX IF EXISTS idx_alerts_created;
DROP INDEX IF EXISTS idx_alerts_assigned_team;
DROP INDEX IF EXISTS idx_alerts_assigned_user;
DROP INDEX IF EXISTS idx_alerts_org_priority;
DROP INDEX IF EXISTS idx_alerts_org_status;
DROP INDEX IF EXISTS idx_alerts_org_id;

-- Drop indexes for teams
DROP INDEX IF EXISTS idx_team_members_user_id;
DROP INDEX IF EXISTS idx_team_members_team_id;
DROP INDEX IF EXISTS idx_teams_org_id;

-- Drop tables
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
