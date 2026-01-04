-- Drop triggers
DROP TRIGGER IF EXISTS update_escalation_rules_updated_at ON escalation_rules;
DROP TRIGGER IF EXISTS update_escalation_policies_updated_at ON escalation_policies;

-- Drop indexes
DROP INDEX IF EXISTS idx_alerts_escalation_policy_id;
DROP INDEX IF EXISTS idx_alert_escalation_events_next_escalation;
DROP INDEX IF EXISTS idx_alert_escalation_events_alert_id;
DROP INDEX IF EXISTS idx_escalation_targets_target;
DROP INDEX IF EXISTS idx_escalation_targets_rule_id;
DROP INDEX IF EXISTS idx_escalation_rules_position;
DROP INDEX IF EXISTS idx_escalation_rules_policy_id;
DROP INDEX IF EXISTS idx_escalation_policies_organization_id;

-- Remove column from alerts
ALTER TABLE alerts DROP COLUMN IF EXISTS escalation_policy_id;

-- Drop tables
DROP TABLE IF EXISTS alert_escalation_events;
DROP TABLE IF EXISTS escalation_targets;
DROP TABLE IF EXISTS escalation_rules;
DROP TABLE IF EXISTS escalation_policies;
