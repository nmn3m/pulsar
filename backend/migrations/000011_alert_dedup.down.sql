-- Drop index
DROP INDEX IF EXISTS idx_alerts_dedup_key;

-- Remove deduplication columns
ALTER TABLE alerts DROP COLUMN IF EXISTS last_occurrence_at;
ALTER TABLE alerts DROP COLUMN IF EXISTS first_occurrence_at;
ALTER TABLE alerts DROP COLUMN IF EXISTS dedup_count;
ALTER TABLE alerts DROP COLUMN IF EXISTS dedup_key;
