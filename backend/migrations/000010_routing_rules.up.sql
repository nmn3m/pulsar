-- Alert Routing Rules
CREATE TABLE IF NOT EXISTS alert_routing_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    priority INTEGER NOT NULL DEFAULT 0,
    conditions JSONB NOT NULL DEFAULT '{}',
    actions JSONB NOT NULL DEFAULT '{}',
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for efficient lookups
CREATE INDEX idx_routing_rules_org ON alert_routing_rules(organization_id);
CREATE INDEX idx_routing_rules_enabled ON alert_routing_rules(organization_id, enabled) WHERE enabled = true;
CREATE INDEX idx_routing_rules_priority ON alert_routing_rules(organization_id, priority);
