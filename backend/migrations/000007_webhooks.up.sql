-- Webhook endpoints table (outgoing webhooks)
CREATE TABLE IF NOT EXISTS webhook_endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    secret VARCHAR(255) NOT NULL, -- For HMAC signature
    enabled BOOLEAN DEFAULT true,

    -- Event filters
    alert_created BOOLEAN DEFAULT false,
    alert_updated BOOLEAN DEFAULT false,
    alert_acknowledged BOOLEAN DEFAULT false,
    alert_closed BOOLEAN DEFAULT false,
    alert_escalated BOOLEAN DEFAULT false,
    incident_created BOOLEAN DEFAULT false,
    incident_updated BOOLEAN DEFAULT false,
    incident_resolved BOOLEAN DEFAULT false,

    -- HTTP configuration
    headers JSONB DEFAULT '{}', -- Custom headers
    timeout_seconds INTEGER DEFAULT 30,

    -- Retry configuration
    max_retries INTEGER DEFAULT 3,
    retry_delay_seconds INTEGER DEFAULT 60,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_webhook_endpoints_org_id ON webhook_endpoints(organization_id);
CREATE INDEX idx_webhook_endpoints_enabled ON webhook_endpoints(enabled);

-- Webhook delivery logs
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_endpoint_id UUID NOT NULL REFERENCES webhook_endpoints(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,

    -- Delivery tracking
    status VARCHAR(20) NOT NULL, -- pending, success, failed
    attempts INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    next_retry_at TIMESTAMP WITH TIME ZONE,

    -- Response tracking
    response_status_code INTEGER,
    response_body TEXT,
    error_message TEXT,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_webhook_deliveries_endpoint_id ON webhook_deliveries(webhook_endpoint_id);
CREATE INDEX idx_webhook_deliveries_org_id ON webhook_deliveries(organization_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at);
CREATE INDEX idx_webhook_deliveries_created_at ON webhook_deliveries(created_at);

-- Incoming webhook tokens (for receiving webhooks from external sources)
CREATE TABLE IF NOT EXISTS incoming_webhook_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    enabled BOOLEAN DEFAULT true,

    -- Integration type
    integration_type VARCHAR(50) NOT NULL, -- generic, prometheus, grafana, datadog, etc.

    -- Default settings for created alerts
    default_priority VARCHAR(10) DEFAULT 'P3',
    default_tags JSONB DEFAULT '[]',

    -- Usage tracking
    last_used_at TIMESTAMP WITH TIME ZONE,
    request_count INTEGER DEFAULT 0,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_incoming_webhook_tokens_org_id ON incoming_webhook_tokens(organization_id);
CREATE INDEX idx_incoming_webhook_tokens_token ON incoming_webhook_tokens(token);
CREATE INDEX idx_incoming_webhook_tokens_enabled ON incoming_webhook_tokens(enabled);

-- Trigger for updated_at
CREATE TRIGGER update_webhook_endpoints_updated_at
    BEFORE UPDATE ON webhook_endpoints
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhook_deliveries_updated_at
    BEFORE UPDATE ON webhook_deliveries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_incoming_webhook_tokens_updated_at
    BEFORE UPDATE ON incoming_webhook_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
