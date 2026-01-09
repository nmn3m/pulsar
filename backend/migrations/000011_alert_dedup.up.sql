-- Add deduplication fields to alerts table
ALTER TABLE alerts ADD COLUMN IF NOT EXISTS dedup_key VARCHAR(512);
ALTER TABLE alerts ADD COLUMN IF NOT EXISTS dedup_count INTEGER DEFAULT 1;
ALTER TABLE alerts ADD COLUMN IF NOT EXISTS first_occurrence_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE alerts ADD COLUMN IF NOT EXISTS last_occurrence_at TIMESTAMP WITH TIME ZONE;

-- Index for efficient dedup lookups
CREATE INDEX IF NOT EXISTS idx_alerts_dedup_key ON alerts(organization_id, dedup_key) WHERE dedup_key IS NOT NULL AND status != 'closed';

-- Update existing alerts to set first_occurrence_at
UPDATE alerts SET first_occurrence_at = created_at WHERE first_occurrence_at IS NULL;
UPDATE alerts SET last_occurrence_at = created_at WHERE last_occurrence_at IS NULL;
