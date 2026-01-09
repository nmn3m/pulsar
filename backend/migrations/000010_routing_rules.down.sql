-- Drop indexes
DROP INDEX IF EXISTS idx_routing_rules_priority;
DROP INDEX IF EXISTS idx_routing_rules_enabled;
DROP INDEX IF EXISTS idx_routing_rules_org;

-- Drop table
DROP TABLE IF EXISTS alert_routing_rules;
