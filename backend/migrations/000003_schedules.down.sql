-- Drop triggers
DROP TRIGGER IF EXISTS update_schedule_overrides_updated_at ON schedule_overrides;
DROP TRIGGER IF EXISTS update_schedule_rotations_updated_at ON schedule_rotations;
DROP TRIGGER IF EXISTS update_schedules_updated_at ON schedules;

-- Drop indexes
DROP INDEX IF EXISTS idx_schedule_overrides_time_range;
DROP INDEX IF EXISTS idx_schedule_overrides_schedule_id;
DROP INDEX IF EXISTS idx_schedule_rotation_participants_user_id;
DROP INDEX IF EXISTS idx_schedule_rotation_participants_rotation_id;
DROP INDEX IF EXISTS idx_schedule_rotations_schedule_id;
DROP INDEX IF EXISTS idx_schedules_team_id;
DROP INDEX IF EXISTS idx_schedules_organization_id;

-- Drop tables
DROP TABLE IF EXISTS schedule_overrides;
DROP TABLE IF EXISTS schedule_rotation_participants;
DROP TABLE IF EXISTS schedule_rotations;
DROP TABLE IF EXISTS schedules;
