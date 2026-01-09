-- Add notification_channels to escalation_targets table
-- This allows per-target notification channel override
-- Format: {"channels": ["email", "slack"], "urgent": true}
ALTER TABLE escalation_targets
ADD COLUMN IF NOT EXISTS notification_channels JSONB DEFAULT NULL;

-- Add comment for documentation
COMMENT ON COLUMN escalation_targets.notification_channels IS 'Override notification channels for this target. NULL = use user preferences. Format: {"channels": ["email", "slack", "sms", "webhook"], "urgent": true}';
