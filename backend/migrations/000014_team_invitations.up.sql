-- Team invitations table
CREATE TABLE IF NOT EXISTS team_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    token VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    invited_by_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for looking up invitations by email
CREATE INDEX IF NOT EXISTS idx_team_invitations_email ON team_invitations(email);

-- Index for looking up invitations by team
CREATE INDEX IF NOT EXISTS idx_team_invitations_team_id ON team_invitations(team_id);

-- Index for looking up invitations by token
CREATE INDEX IF NOT EXISTS idx_team_invitations_token ON team_invitations(token);

-- Index for looking up pending invitations
CREATE INDEX IF NOT EXISTS idx_team_invitations_status ON team_invitations(status) WHERE status = 'pending';

-- Unique constraint: one pending invitation per email per team
CREATE UNIQUE INDEX IF NOT EXISTS idx_team_invitations_unique_pending
ON team_invitations(team_id, email) WHERE status = 'pending';
