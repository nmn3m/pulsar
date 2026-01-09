-- User DND (Do Not Disturb) Settings
CREATE TABLE IF NOT EXISTS user_dnd_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    enabled BOOLEAN DEFAULT false,
    -- Schedule stored as JSONB for flexibility
    -- Format: {"weekly": [{"day": "monday", "start": "22:00", "end": "08:00"}, ...], "timezone": "America/New_York"}
    schedule JSONB DEFAULT '{}',
    -- One-time overrides (e.g., "I'm on vacation until date X")
    -- Format: [{"start": "2024-01-01T00:00:00Z", "end": "2024-01-07T00:00:00Z", "reason": "Vacation"}]
    overrides JSONB DEFAULT '[]',
    -- Allow high priority (P1) alerts to bypass DND
    allow_p1_override BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Index for efficient lookups
CREATE INDEX IF NOT EXISTS idx_dnd_settings_user ON user_dnd_settings(user_id);
CREATE INDEX IF NOT EXISTS idx_dnd_settings_enabled ON user_dnd_settings(user_id) WHERE enabled = true;
