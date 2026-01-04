-- On-call schedules table
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Schedule rotations - defines rotation patterns
CREATE TABLE IF NOT EXISTS schedule_rotations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    rotation_type VARCHAR(50) NOT NULL, -- daily, weekly, custom
    rotation_length INTEGER NOT NULL DEFAULT 1, -- how many days/weeks for each person
    start_date DATE NOT NULL, -- when this rotation starts
    start_time TIME NOT NULL DEFAULT '00:00:00', -- time of day rotation starts
    end_time TIME, -- time of day rotation ends (NULL = 24/7)
    handoff_day INTEGER, -- day of week for handoff (0-6, Sunday=0)
    handoff_time TIME NOT NULL DEFAULT '09:00:00', -- time of handoff
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Rotation participants - who is in each rotation
CREATE TABLE IF NOT EXISTS schedule_rotation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rotation_id UUID NOT NULL REFERENCES schedule_rotations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER NOT NULL, -- order in rotation (0-based)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(rotation_id, user_id),
    UNIQUE(rotation_id, position)
);

-- Schedule overrides - temporary changes to who is on-call
CREATE TABLE IF NOT EXISTS schedule_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_override_time CHECK (end_time > start_time)
);

-- Indexes for performance
CREATE INDEX idx_schedules_organization_id ON schedules(organization_id);
CREATE INDEX idx_schedules_team_id ON schedules(team_id);
CREATE INDEX idx_schedule_rotations_schedule_id ON schedule_rotations(schedule_id);
CREATE INDEX idx_schedule_rotation_participants_rotation_id ON schedule_rotation_participants(rotation_id);
CREATE INDEX idx_schedule_rotation_participants_user_id ON schedule_rotation_participants(user_id);
CREATE INDEX idx_schedule_overrides_schedule_id ON schedule_overrides(schedule_id);
CREATE INDEX idx_schedule_overrides_time_range ON schedule_overrides(start_time, end_time);

-- Trigger for updated_at on schedules
CREATE TRIGGER update_schedules_updated_at
    BEFORE UPDATE ON schedules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for updated_at on schedule_rotations
CREATE TRIGGER update_schedule_rotations_updated_at
    BEFORE UPDATE ON schedule_rotations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for updated_at on schedule_overrides
CREATE TRIGGER update_schedule_overrides_updated_at
    BEFORE UPDATE ON schedule_overrides
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
