package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type webhookRepository struct {
	db *sqlx.DB
}

func NewWebhookRepository(db *sqlx.DB) repository.WebhookRepository {
	return &webhookRepository{db: db}
}

// Webhook Endpoints

func (r *webhookRepository) CreateEndpoint(ctx context.Context, endpoint *domain.WebhookEndpoint) error {
	query := `
		INSERT INTO webhook_endpoints (
			id, organization_id, name, url, secret, enabled,
			alert_created, alert_updated, alert_acknowledged, alert_closed, alert_escalated,
			incident_created, incident_updated, incident_resolved,
			headers, timeout_seconds, max_retries, retry_delay_seconds,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14,
			$15, $16, $17, $18,
			$19, $20
		)
	`

	headersJSON, err := json.Marshal(endpoint.Headers)
	if err != nil {
		return err
	}

	now := time.Now()
	endpoint.CreatedAt = now
	endpoint.UpdatedAt = now

	_, err = r.db.ExecContext(ctx, query,
		endpoint.ID,
		endpoint.OrganizationID,
		endpoint.Name,
		endpoint.URL,
		endpoint.Secret,
		endpoint.Enabled,
		endpoint.AlertCreated,
		endpoint.AlertUpdated,
		endpoint.AlertAcknowledged,
		endpoint.AlertClosed,
		endpoint.AlertEscalated,
		endpoint.IncidentCreated,
		endpoint.IncidentUpdated,
		endpoint.IncidentResolved,
		headersJSON,
		endpoint.TimeoutSeconds,
		endpoint.MaxRetries,
		endpoint.RetryDelaySeconds,
		endpoint.CreatedAt,
		endpoint.UpdatedAt,
	)

	return err
}

func (r *webhookRepository) GetEndpointByID(ctx context.Context, id uuid.UUID) (*domain.WebhookEndpoint, error) {
	query := `
		SELECT id, organization_id, name, url, secret, enabled,
			alert_created, alert_updated, alert_acknowledged, alert_closed, alert_escalated,
			incident_created, incident_updated, incident_resolved,
			headers, timeout_seconds, max_retries, retry_delay_seconds,
			created_at, updated_at
		FROM webhook_endpoints
		WHERE id = $1
	`

	var endpoint domain.WebhookEndpoint
	var headersJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&endpoint.ID,
		&endpoint.OrganizationID,
		&endpoint.Name,
		&endpoint.URL,
		&endpoint.Secret,
		&endpoint.Enabled,
		&endpoint.AlertCreated,
		&endpoint.AlertUpdated,
		&endpoint.AlertAcknowledged,
		&endpoint.AlertClosed,
		&endpoint.AlertEscalated,
		&endpoint.IncidentCreated,
		&endpoint.IncidentUpdated,
		&endpoint.IncidentResolved,
		&headersJSON,
		&endpoint.TimeoutSeconds,
		&endpoint.MaxRetries,
		&endpoint.RetryDelaySeconds,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if len(headersJSON) > 0 {
		if err := json.Unmarshal(headersJSON, &endpoint.Headers); err != nil {
			return nil, err
		}
	}

	return &endpoint, nil
}

func (r *webhookRepository) ListEndpoints(ctx context.Context, orgID uuid.UUID) ([]*domain.WebhookEndpoint, error) {
	query := `
		SELECT id, organization_id, name, url, secret, enabled,
			alert_created, alert_updated, alert_acknowledged, alert_closed, alert_escalated,
			incident_created, incident_updated, incident_resolved,
			headers, timeout_seconds, max_retries, retry_delay_seconds,
			created_at, updated_at
		FROM webhook_endpoints
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []*domain.WebhookEndpoint
	for rows.Next() {
		var endpoint domain.WebhookEndpoint
		var headersJSON []byte

		err := rows.Scan(
			&endpoint.ID,
			&endpoint.OrganizationID,
			&endpoint.Name,
			&endpoint.URL,
			&endpoint.Secret,
			&endpoint.Enabled,
			&endpoint.AlertCreated,
			&endpoint.AlertUpdated,
			&endpoint.AlertAcknowledged,
			&endpoint.AlertClosed,
			&endpoint.AlertEscalated,
			&endpoint.IncidentCreated,
			&endpoint.IncidentUpdated,
			&endpoint.IncidentResolved,
			&headersJSON,
			&endpoint.TimeoutSeconds,
			&endpoint.MaxRetries,
			&endpoint.RetryDelaySeconds,
			&endpoint.CreatedAt,
			&endpoint.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if len(headersJSON) > 0 {
			if err := json.Unmarshal(headersJSON, &endpoint.Headers); err != nil {
				return nil, err
			}
		}

		endpoints = append(endpoints, &endpoint)
	}

	return endpoints, rows.Err()
}

func (r *webhookRepository) UpdateEndpoint(ctx context.Context, endpoint *domain.WebhookEndpoint) error {
	query := `
		UPDATE webhook_endpoints
		SET name = $1, url = $2, enabled = $3,
			alert_created = $4, alert_updated = $5, alert_acknowledged = $6,
			alert_closed = $7, alert_escalated = $8,
			incident_created = $9, incident_updated = $10, incident_resolved = $11,
			headers = $12, timeout_seconds = $13, max_retries = $14,
			retry_delay_seconds = $15, updated_at = $16
		WHERE id = $17 AND organization_id = $18
	`

	headersJSON, err := json.Marshal(endpoint.Headers)
	if err != nil {
		return err
	}

	endpoint.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		endpoint.Name,
		endpoint.URL,
		endpoint.Enabled,
		endpoint.AlertCreated,
		endpoint.AlertUpdated,
		endpoint.AlertAcknowledged,
		endpoint.AlertClosed,
		endpoint.AlertEscalated,
		endpoint.IncidentCreated,
		endpoint.IncidentUpdated,
		endpoint.IncidentResolved,
		headersJSON,
		endpoint.TimeoutSeconds,
		endpoint.MaxRetries,
		endpoint.RetryDelaySeconds,
		endpoint.UpdatedAt,
		endpoint.ID,
		endpoint.OrganizationID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *webhookRepository) DeleteEndpoint(ctx context.Context, id, orgID uuid.UUID) error {
	query := `DELETE FROM webhook_endpoints WHERE id = $1 AND organization_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, orgID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Webhook Deliveries

func (r *webhookRepository) CreateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
	query := `
		INSERT INTO webhook_deliveries (
			id, webhook_endpoint_id, organization_id, event_type, payload,
			status, attempts, last_attempt_at, next_retry_at,
			response_status_code, response_body, error_message,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	payloadJSON, err := json.Marshal(delivery.Payload)
	if err != nil {
		return err
	}

	now := time.Now()
	delivery.CreatedAt = now
	delivery.UpdatedAt = now

	_, err = r.db.ExecContext(ctx, query,
		delivery.ID,
		delivery.WebhookEndpointID,
		delivery.OrganizationID,
		delivery.EventType,
		payloadJSON,
		delivery.Status,
		delivery.Attempts,
		delivery.LastAttemptAt,
		delivery.NextRetryAt,
		delivery.ResponseStatus,
		delivery.ResponseBody,
		delivery.ErrorMessage,
		delivery.CreatedAt,
		delivery.UpdatedAt,
	)

	return err
}

func (r *webhookRepository) UpdateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error {
	query := `
		UPDATE webhook_deliveries
		SET status = $1, attempts = $2, last_attempt_at = $3, next_retry_at = $4,
			response_status_code = $5, response_body = $6, error_message = $7,
			updated_at = $8
		WHERE id = $9
	`

	delivery.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		delivery.Status,
		delivery.Attempts,
		delivery.LastAttemptAt,
		delivery.NextRetryAt,
		delivery.ResponseStatus,
		delivery.ResponseBody,
		delivery.ErrorMessage,
		delivery.UpdatedAt,
		delivery.ID,
	)

	return err
}

func (r *webhookRepository) GetPendingDeliveries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_endpoint_id, organization_id, event_type, payload,
			status, attempts, last_attempt_at, next_retry_at,
			response_status_code, response_body, error_message,
			created_at, updated_at
		FROM webhook_deliveries
		WHERE status = $1 AND (next_retry_at IS NULL OR next_retry_at <= $2)
		ORDER BY created_at ASC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, domain.WebhookDeliveryPending, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		var delivery domain.WebhookDelivery
		var payloadJSON []byte

		err := rows.Scan(
			&delivery.ID,
			&delivery.WebhookEndpointID,
			&delivery.OrganizationID,
			&delivery.EventType,
			&payloadJSON,
			&delivery.Status,
			&delivery.Attempts,
			&delivery.LastAttemptAt,
			&delivery.NextRetryAt,
			&delivery.ResponseStatus,
			&delivery.ResponseBody,
			&delivery.ErrorMessage,
			&delivery.CreatedAt,
			&delivery.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &delivery.Payload); err != nil {
				return nil, err
			}
		}

		deliveries = append(deliveries, &delivery)
	}

	return deliveries, rows.Err()
}

func (r *webhookRepository) ListDeliveries(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, webhook_endpoint_id, organization_id, event_type, payload,
			status, attempts, last_attempt_at, next_retry_at,
			response_status_code, response_body, error_message,
			created_at, updated_at
		FROM webhook_deliveries
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		var delivery domain.WebhookDelivery
		var payloadJSON []byte

		err := rows.Scan(
			&delivery.ID,
			&delivery.WebhookEndpointID,
			&delivery.OrganizationID,
			&delivery.EventType,
			&payloadJSON,
			&delivery.Status,
			&delivery.Attempts,
			&delivery.LastAttemptAt,
			&delivery.NextRetryAt,
			&delivery.ResponseStatus,
			&delivery.ResponseBody,
			&delivery.ErrorMessage,
			&delivery.CreatedAt,
			&delivery.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &delivery.Payload); err != nil {
				return nil, err
			}
		}

		deliveries = append(deliveries, &delivery)
	}

	return deliveries, rows.Err()
}

// Incoming Webhook Tokens

func (r *webhookRepository) CreateIncomingToken(ctx context.Context, token *domain.IncomingWebhookToken) error {
	query := `
		INSERT INTO incoming_webhook_tokens (
			id, organization_id, name, token, enabled, integration_type,
			default_priority, default_tags, last_used_at, request_count,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	tagsJSON, err := json.Marshal(token.DefaultTags)
	if err != nil {
		return err
	}

	now := time.Now()
	token.CreatedAt = now
	token.UpdatedAt = now

	_, err = r.db.ExecContext(ctx, query,
		token.ID,
		token.OrganizationID,
		token.Name,
		token.Token,
		token.Enabled,
		token.IntegrationType,
		token.DefaultPriority,
		tagsJSON,
		token.LastUsedAt,
		token.RequestCount,
		token.CreatedAt,
		token.UpdatedAt,
	)

	return err
}

func (r *webhookRepository) GetIncomingTokenByToken(ctx context.Context, tokenStr string) (*domain.IncomingWebhookToken, error) {
	query := `
		SELECT id, organization_id, name, token, enabled, integration_type,
			default_priority, default_tags, last_used_at, request_count,
			created_at, updated_at
		FROM incoming_webhook_tokens
		WHERE token = $1
	`

	var token domain.IncomingWebhookToken
	var tagsJSON []byte

	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(
		&token.ID,
		&token.OrganizationID,
		&token.Name,
		&token.Token,
		&token.Enabled,
		&token.IntegrationType,
		&token.DefaultPriority,
		&tagsJSON,
		&token.LastUsedAt,
		&token.RequestCount,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if len(tagsJSON) > 0 {
		if err := json.Unmarshal(tagsJSON, &token.DefaultTags); err != nil {
			return nil, err
		}
	}

	return &token, nil
}

func (r *webhookRepository) ListIncomingTokens(ctx context.Context, orgID uuid.UUID) ([]*domain.IncomingWebhookToken, error) {
	query := `
		SELECT id, organization_id, name, token, enabled, integration_type,
			default_priority, default_tags, last_used_at, request_count,
			created_at, updated_at
		FROM incoming_webhook_tokens
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*domain.IncomingWebhookToken
	for rows.Next() {
		var token domain.IncomingWebhookToken
		var tagsJSON []byte

		err := rows.Scan(
			&token.ID,
			&token.OrganizationID,
			&token.Name,
			&token.Token,
			&token.Enabled,
			&token.IntegrationType,
			&token.DefaultPriority,
			&tagsJSON,
			&token.LastUsedAt,
			&token.RequestCount,
			&token.CreatedAt,
			&token.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &token.DefaultTags); err != nil {
				return nil, err
			}
		}

		tokens = append(tokens, &token)
	}

	return tokens, rows.Err()
}

func (r *webhookRepository) UpdateIncomingTokenUsage(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE incoming_webhook_tokens
		SET last_used_at = $1, request_count = request_count + 1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	return err
}

func (r *webhookRepository) DeleteIncomingToken(ctx context.Context, id, orgID uuid.UUID) error {
	query := `DELETE FROM incoming_webhook_tokens WHERE id = $1 AND organization_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, orgID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
