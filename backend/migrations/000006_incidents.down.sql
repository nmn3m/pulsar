-- Drop trigger
DROP TRIGGER IF EXISTS incidents_updated_at ON incidents;
DROP FUNCTION IF EXISTS update_incidents_updated_at();

-- Drop tables in reverse order
DROP TABLE IF EXISTS incident_alerts;
DROP TABLE IF EXISTS incident_timeline;
DROP TABLE IF EXISTS incident_responders;
DROP TABLE IF EXISTS incidents;
