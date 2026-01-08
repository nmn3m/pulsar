export type WebhookDeliveryStatus = 'pending' | 'success' | 'failed';

export type IncomingWebhookIntegrationType = 'generic' | 'prometheus' | 'grafana' | 'datadog';

export interface WebhookEndpoint {
  id: string;
  organization_id: string;
  name: string;
  url: string;
  enabled: boolean;

  // Event filters
  alert_created: boolean;
  alert_updated: boolean;
  alert_acknowledged: boolean;
  alert_closed: boolean;
  alert_escalated: boolean;
  incident_created: boolean;
  incident_updated: boolean;
  incident_resolved: boolean;

  // HTTP configuration
  headers: Record<string, string>;
  timeout_seconds: number;

  // Retry configuration
  max_retries: number;
  retry_delay_seconds: number;

  created_at: string;
  updated_at: string;
}

export interface WebhookDelivery {
  id: string;
  webhook_endpoint_id: string;
  organization_id: string;
  event_type: string;
  payload: Record<string, any>;
  status: WebhookDeliveryStatus;
  attempts: number;
  last_attempt_at?: string;
  next_retry_at?: string;
  response_status_code?: number;
  response_body?: string;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

export interface IncomingWebhookToken {
  id: string;
  organization_id: string;
  name: string;
  token: string;
  enabled: boolean;
  integration_type: IncomingWebhookIntegrationType;
  default_priority: string;
  default_tags: string[];
  last_used_at?: string;
  request_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateWebhookEndpointRequest {
  name: string;
  url: string;
  enabled?: boolean;
  alert_created?: boolean;
  alert_updated?: boolean;
  alert_acknowledged?: boolean;
  alert_closed?: boolean;
  alert_escalated?: boolean;
  incident_created?: boolean;
  incident_updated?: boolean;
  incident_resolved?: boolean;
  headers?: Record<string, string>;
  timeout_seconds?: number;
  max_retries?: number;
  retry_delay_seconds?: number;
}

export interface UpdateWebhookEndpointRequest {
  name?: string;
  url?: string;
  enabled?: boolean;
  alert_created?: boolean;
  alert_updated?: boolean;
  alert_acknowledged?: boolean;
  alert_closed?: boolean;
  alert_escalated?: boolean;
  incident_created?: boolean;
  incident_updated?: boolean;
  incident_resolved?: boolean;
  headers?: Record<string, string>;
  timeout_seconds?: number;
  max_retries?: number;
  retry_delay_seconds?: number;
}

export interface CreateIncomingWebhookTokenRequest {
  name: string;
  integration_type: IncomingWebhookIntegrationType;
  default_priority?: string;
  default_tags?: string[];
}

export interface ListWebhookDeliveriesResponse {
  deliveries: WebhookDelivery[];
  limit: number;
  offset: number;
}
