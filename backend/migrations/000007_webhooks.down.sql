-- Drop triggers
DROP TRIGGER IF EXISTS update_incoming_webhook_tokens_updated_at ON incoming_webhook_tokens;
DROP TRIGGER IF EXISTS update_webhook_deliveries_updated_at ON webhook_deliveries;
DROP TRIGGER IF EXISTS update_webhook_endpoints_updated_at ON webhook_endpoints;

-- Drop tables
DROP TABLE IF EXISTS incoming_webhook_tokens;
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS webhook_endpoints;
