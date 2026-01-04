-- Teams table
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(organization_id, name)
);

-- Team members table
CREATE TABLE team_members (
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (team_id, user_id)
);

-- Alerts table
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    source VARCHAR(100) NOT NULL,
    source_id VARCHAR(255),
    priority VARCHAR(20) NOT NULL DEFAULT 'P3',
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    message TEXT NOT NULL,
    description TEXT,
    tags JSONB DEFAULT '[]',
    custom_fields JSONB DEFAULT '{}',

    -- Assignment
    assigned_to_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_to_team_id UUID REFERENCES teams(id) ON DELETE SET NULL,

    -- Acknowledgment
    acknowledged_by UUID REFERENCES users(id) ON DELETE SET NULL,
    acknowledged_at TIMESTAMP WITH TIME ZONE,

    -- Closure
    closed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    closed_at TIMESTAMP WITH TIME ZONE,
    close_reason TEXT,

    -- Snoozed
    snoozed_until TIMESTAMP WITH TIME ZONE,

    -- Escalation
    escalation_policy_id UUID,
    escalation_level INTEGER DEFAULT 0,
    last_escalated_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for alerts
CREATE INDEX idx_alerts_org_id ON alerts(organization_id);
CREATE INDEX idx_alerts_org_status ON alerts(organization_id, status);
CREATE INDEX idx_alerts_org_priority ON alerts(organization_id, priority);
CREATE INDEX idx_alerts_assigned_user ON alerts(assigned_to_user_id);
CREATE INDEX idx_alerts_assigned_team ON alerts(assigned_to_team_id);
CREATE INDEX idx_alerts_created ON alerts(created_at DESC);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_source ON alerts(source);

-- Create indexes for teams
CREATE INDEX idx_teams_org_id ON teams(organization_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);

-- Add triggers for updated_at
CREATE TRIGGER update_teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alerts_updated_at
    BEFORE UPDATE ON alerts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
