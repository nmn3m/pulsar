ALTER TABLE notification_channels DROP CONSTRAINT IF EXISTS valid_channel_type;
ALTER TABLE notification_channels ADD CONSTRAINT valid_channel_type
    CHECK (channel_type IN ('email', 'slack', 'teams', 'webhook', 'sms'));
