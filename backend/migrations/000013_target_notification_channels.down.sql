-- Remove notification_channels from escalation_targets table
ALTER TABLE escalation_targets
DROP COLUMN IF EXISTS notification_channels;
