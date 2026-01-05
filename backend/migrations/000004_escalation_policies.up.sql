-- Escalation policies table
CREATE TABLE IF NOT EXISTS escalation_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    repeat_enabled BOOLEAN NOT NULL DEFAULT false,
    repeat_count INTEGER, -- NULL = infinite, number = max repeats
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Escalation rules - steps in an escalation policy
CREATE TABLE IF NOT EXISTS escalation_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL REFERENCES escalation_policies(id) ON DELETE CASCADE,
    position INTEGER NOT NULL, -- order of execution (0-based)
    escalation_delay INTEGER NOT NULL DEFAULT 0, -- minutes before escalating
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(policy_id, position)
);

-- Escalation targets - who to notify in each rule
CREATE TABLE IF NOT EXISTS escalation_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id UUID NOT NULL REFERENCES escalation_rules(id) ON DELETE CASCADE,
    target_type VARCHAR(50) NOT NULL, -- user, team, schedule
    target_id UUID NOT NULL, -- references users, teams, or schedules
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_target_type CHECK (target_type IN ('user', 'team', 'schedule'))
);

-- Escalation events - tracks escalation state for alerts
CREATE TABLE IF NOT EXISTS alert_escalation_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id UUID NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    policy_id UUID NOT NULL REFERENCES escalation_policies(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES escalation_rules(id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL, -- triggered, acknowledged, completed, stopped
    current_level INTEGER NOT NULL DEFAULT 0, -- which rule position we're at
    repeat_count INTEGER NOT NULL DEFAULT 0, -- how many times we've repeated
    next_escalation_at TIMESTAMP, -- when to escalate next
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_event_type CHECK (event_type IN ('triggered', 'acknowledged', 'completed', 'stopped'))
);

-- Add escalation_policy_id to alerts table
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'alerts' AND column_name = 'escalation_policy_id'
    ) THEN
        ALTER TABLE alerts ADD COLUMN escalation_policy_id UUID REFERENCES escalation_policies(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_escalation_policies_organization_id ON escalation_policies(organization_id);
CREATE INDEX IF NOT EXISTS idx_escalation_rules_policy_id ON escalation_rules(policy_id);
CREATE INDEX IF NOT EXISTS idx_escalation_rules_position ON escalation_rules(policy_id, position);
CREATE INDEX IF NOT EXISTS idx_escalation_targets_rule_id ON escalation_targets(rule_id);
CREATE INDEX IF NOT EXISTS idx_escalation_targets_target ON escalation_targets(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_alert_escalation_events_alert_id ON alert_escalation_events(alert_id);
CREATE INDEX IF NOT EXISTS idx_alert_escalation_events_next_escalation ON alert_escalation_events(next_escalation_at) WHERE next_escalation_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_alerts_escalation_policy_id ON alerts(escalation_policy_id);

-- Trigger for updated_at on escalation_policies
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'update_escalation_policies_updated_at'
    ) THEN
        CREATE TRIGGER update_escalation_policies_updated_at
            BEFORE UPDATE ON escalation_policies
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Trigger for updated_at on escalation_rules
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'update_escalation_rules_updated_at'
    ) THEN
        CREATE TRIGGER update_escalation_rules_updated_at
            BEFORE UPDATE ON escalation_rules
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
