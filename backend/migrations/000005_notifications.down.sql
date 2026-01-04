-- Drop indexes
DROP INDEX IF EXISTS idx_notification_logs_created_at;
DROP INDEX IF EXISTS idx_notification_logs_status;
DROP INDEX IF EXISTS idx_notification_logs_alert_id;
DROP INDEX IF EXISTS idx_notification_logs_user_id;
DROP INDEX IF EXISTS idx_notification_logs_channel_id;
DROP INDEX IF EXISTS idx_notification_logs_organization_id;
DROP INDEX IF EXISTS idx_user_notification_preferences_channel_id;
DROP INDEX IF EXISTS idx_user_notification_preferences_user_id;
DROP INDEX IF EXISTS idx_notification_channels_type;
DROP INDEX IF EXISTS idx_notification_channels_organization_id;

-- Drop triggers
DROP TRIGGER IF EXISTS update_user_notification_preferences_updated_at ON user_notification_preferences;
DROP TRIGGER IF EXISTS update_notification_channels_updated_at ON notification_channels;

-- Drop notification tables
DROP TABLE IF EXISTS notification_logs;
DROP TABLE IF EXISTS user_notification_preferences;
DROP TABLE IF EXISTS notification_channels;
