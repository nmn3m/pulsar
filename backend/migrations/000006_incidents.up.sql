-- Create incidents table
CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('investigating', 'identified', 'monitoring', 'resolved')),
    priority VARCHAR(10) NOT NULL CHECK (priority IN ('P1', 'P2', 'P3', 'P4', 'P5')),
    created_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_to_team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on organization_id for filtering
CREATE INDEX idx_incidents_organization_id ON incidents(organization_id);

-- Create index on status for filtering
CREATE INDEX idx_incidents_status ON incidents(status);

-- Create index on severity for filtering
CREATE INDEX idx_incidents_severity ON incidents(severity);

-- Create index on assigned_to_team_id for filtering
CREATE INDEX idx_incidents_assigned_to_team_id ON incidents(assigned_to_team_id);

-- Create incident_responders table
CREATE TABLE IF NOT EXISTS incident_responders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('incident_commander', 'responder')),
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(incident_id, user_id)
);

-- Create index on incident_id for lookups
CREATE INDEX idx_incident_responders_incident_id ON incident_responders(incident_id);

-- Create index on user_id for lookups
CREATE INDEX idx_incident_responders_user_id ON incident_responders(user_id);

-- Create incident_timeline table
CREATE TABLE IF NOT EXISTS incident_timeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN (
        'created',
        'status_changed',
        'severity_changed',
        'responder_added',
        'responder_removed',
        'note_added',
        'alert_linked',
        'alert_unlinked',
        'resolved'
    )),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    description TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on incident_id for timeline queries
CREATE INDEX idx_incident_timeline_incident_id ON incident_timeline(incident_id);

-- Create index on created_at for chronological ordering
CREATE INDEX idx_incident_timeline_created_at ON incident_timeline(created_at DESC);

-- Create incident_alerts table (link alerts to incidents)
CREATE TABLE IF NOT EXISTS incident_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    alert_id UUID NOT NULL REFERENCES alerts(id) ON DELETE CASCADE,
    linked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    linked_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(incident_id, alert_id)
);

-- Create index on incident_id for lookups
CREATE INDEX idx_incident_alerts_incident_id ON incident_alerts(incident_id);

-- Create index on alert_id for lookups
CREATE INDEX idx_incident_alerts_alert_id ON incident_alerts(alert_id);

-- Create trigger to update updated_at on incidents
CREATE OR REPLACE FUNCTION update_incidents_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER incidents_updated_at
    BEFORE UPDATE ON incidents
    FOR EACH ROW
    EXECUTE FUNCTION update_incidents_updated_at();
