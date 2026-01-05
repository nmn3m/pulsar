package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
	"go.uber.org/zap"
)

type WebhookService struct {
	webhookRepo repository.WebhookRepository
	logger      *zap.Logger
	httpClient  *http.Client
}

func NewWebhookService(webhookRepo repository.WebhookRepository, logger *zap.Logger) *WebhookService {
	return &WebhookService{
		webhookRepo: webhookRepo,
		logger:      logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Endpoint Management

func (s *WebhookService) CreateEndpoint(ctx context.Context, orgID uuid.UUID, req *domain.CreateWebhookEndpointRequest) (*domain.WebhookEndpoint, error) {
	// Generate a secure secret for HMAC signing
	secret, err := generateSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	endpoint := &domain.WebhookEndpoint{
		ID:                uuid.New(),
		OrganizationID:    orgID,
		Name:              req.Name,
		URL:               req.URL,
		Secret:            secret,
		Enabled:           req.Enabled,
		AlertCreated:      req.AlertCreated,
		AlertUpdated:      req.AlertUpdated,
		AlertAcknowledged: req.AlertAcknowledged,
		AlertClosed:       req.AlertClosed,
		AlertEscalated:    req.AlertEscalated,
		IncidentCreated:   req.IncidentCreated,
		IncidentUpdated:   req.IncidentUpdated,
		IncidentResolved:  req.IncidentResolved,
		Headers:           req.Headers,
		TimeoutSeconds:    getIntOrDefault(req.TimeoutSeconds, 30),
		MaxRetries:        getIntOrDefault(req.MaxRetries, 3),
		RetryDelaySeconds: getIntOrDefault(req.RetryDelaySeconds, 60),
	}

	if endpoint.Headers == nil {
		endpoint.Headers = make(map[string]string)
	}

	if err := s.webhookRepo.CreateEndpoint(ctx, endpoint); err != nil {
		return nil, err
	}

	return endpoint, nil
}

func (s *WebhookService) GetEndpoint(ctx context.Context, id uuid.UUID) (*domain.WebhookEndpoint, error) {
	return s.webhookRepo.GetEndpointByID(ctx, id)
}

func (s *WebhookService) ListEndpoints(ctx context.Context, orgID uuid.UUID) ([]*domain.WebhookEndpoint, error) {
	return s.webhookRepo.ListEndpoints(ctx, orgID)
}

func (s *WebhookService) UpdateEndpoint(ctx context.Context, id, orgID uuid.UUID, req *domain.UpdateWebhookEndpointRequest) (*domain.WebhookEndpoint, error) {
	endpoint, err := s.webhookRepo.GetEndpointByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if endpoint.OrganizationID != orgID {
		return nil, domain.ErrNotFound
	}

	// Update fields if provided
	if req.Name != nil {
		endpoint.Name = *req.Name
	}
	if req.URL != nil {
		endpoint.URL = *req.URL
	}
	if req.Enabled != nil {
		endpoint.Enabled = *req.Enabled
	}
	if req.AlertCreated != nil {
		endpoint.AlertCreated = *req.AlertCreated
	}
	if req.AlertUpdated != nil {
		endpoint.AlertUpdated = *req.AlertUpdated
	}
	if req.AlertAcknowledged != nil {
		endpoint.AlertAcknowledged = *req.AlertAcknowledged
	}
	if req.AlertClosed != nil {
		endpoint.AlertClosed = *req.AlertClosed
	}
	if req.AlertEscalated != nil {
		endpoint.AlertEscalated = *req.AlertEscalated
	}
	if req.IncidentCreated != nil {
		endpoint.IncidentCreated = *req.IncidentCreated
	}
	if req.IncidentUpdated != nil {
		endpoint.IncidentUpdated = *req.IncidentUpdated
	}
	if req.IncidentResolved != nil {
		endpoint.IncidentResolved = *req.IncidentResolved
	}
	if req.Headers != nil {
		endpoint.Headers = req.Headers
	}
	if req.TimeoutSeconds != nil {
		endpoint.TimeoutSeconds = *req.TimeoutSeconds
	}
	if req.MaxRetries != nil {
		endpoint.MaxRetries = *req.MaxRetries
	}
	if req.RetryDelaySeconds != nil {
		endpoint.RetryDelaySeconds = *req.RetryDelaySeconds
	}

	if err := s.webhookRepo.UpdateEndpoint(ctx, endpoint); err != nil {
		return nil, err
	}

	return endpoint, nil
}

func (s *WebhookService) DeleteEndpoint(ctx context.Context, id, orgID uuid.UUID) error {
	return s.webhookRepo.DeleteEndpoint(ctx, id, orgID)
}

// Webhook Delivery

func (s *WebhookService) TriggerWebhooks(ctx context.Context, orgID uuid.UUID, eventType string, data map[string]interface{}) {
	// Run asynchronously to not block the caller
	go func() {
		endpoints, err := s.webhookRepo.ListEndpoints(context.Background(), orgID)
		if err != nil {
			s.logger.Error("Failed to list webhook endpoints", zap.Error(err))
			return
		}

		for _, endpoint := range endpoints {
			if !endpoint.Enabled {
				continue
			}

			if !endpoint.ShouldTriggerEvent(eventType) {
				continue
			}

			payload := &domain.WebhookPayload{
				EventType:      eventType,
				EventID:        uuid.New().String(),
				OrganizationID: orgID.String(),
				Timestamp:      time.Now(),
				Data:           data,
			}

			delivery := &domain.WebhookDelivery{
				ID:                uuid.New(),
				WebhookEndpointID: endpoint.ID,
				OrganizationID:    orgID,
				EventType:         eventType,
				Payload:           data,
				Status:            domain.WebhookDeliveryPending,
				Attempts:          0,
			}

			if err := s.webhookRepo.CreateDelivery(context.Background(), delivery); err != nil {
				s.logger.Error("Failed to create webhook delivery", zap.Error(err))
				continue
			}

			// Attempt immediate delivery
			s.deliverWebhook(context.Background(), endpoint, delivery, payload)
		}
	}()
}

func (s *WebhookService) deliverWebhook(ctx context.Context, endpoint *domain.WebhookEndpoint, delivery *domain.WebhookDelivery, payload *domain.WebhookPayload) {
	delivery.Attempts++
	now := time.Now()
	delivery.LastAttemptAt = &now

	// Serialize payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("Failed to serialize webhook payload", zap.Error(err))
		s.markDeliveryFailed(ctx, delivery, "Failed to serialize payload: "+err.Error())
		return
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint.URL, bytes.NewReader(payloadBytes))
	if err != nil {
		s.logger.Error("Failed to create HTTP request", zap.Error(err))
		s.markDeliveryFailed(ctx, delivery, "Failed to create request: "+err.Error())
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Pulsar-Webhooks/1.0")

	// Add custom headers
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	// Generate and add HMAC signature
	signature := generateHMACSignature(payloadBytes, endpoint.Secret)
	req.Header.Set("X-Pulsar-Signature", signature)
	req.Header.Set("X-Pulsar-Event", payload.EventType)
	req.Header.Set("X-Pulsar-Delivery", delivery.ID.String())

	// Set custom timeout
	client := &http.Client{
		Timeout: time.Duration(endpoint.TimeoutSeconds) * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Failed to send webhook",
			zap.String("url", endpoint.URL),
			zap.Error(err),
		)
		s.handleDeliveryError(ctx, endpoint, delivery, err.Error())
		return
	}
	defer resp.Body.Close()

	// Read response body (limit to 1MB)
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		s.logger.Warn("Failed to read webhook response body", zap.Error(err))
	}

	responseBodyStr := string(responseBody)
	delivery.ResponseStatus = &resp.StatusCode
	delivery.ResponseBody = &responseBodyStr

	// Check if successful (2xx status code)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		delivery.Status = domain.WebhookDeliverySuccess
		delivery.NextRetryAt = nil
		delivery.ErrorMessage = nil

		if err := s.webhookRepo.UpdateDelivery(ctx, delivery); err != nil {
			s.logger.Error("Failed to update webhook delivery", zap.Error(err))
		}

		s.logger.Info("Webhook delivered successfully",
			zap.String("endpoint", endpoint.Name),
			zap.String("url", endpoint.URL),
			zap.Int("status", resp.StatusCode),
		)
	} else {
		errMsg := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, responseBodyStr)
		s.handleDeliveryError(ctx, endpoint, delivery, errMsg)
	}
}

func (s *WebhookService) handleDeliveryError(ctx context.Context, endpoint *domain.WebhookEndpoint, delivery *domain.WebhookDelivery, errMsg string) {
	if delivery.Attempts >= endpoint.MaxRetries {
		s.markDeliveryFailed(ctx, delivery, errMsg)
	} else {
		// Schedule retry
		nextRetry := time.Now().Add(time.Duration(endpoint.RetryDelaySeconds) * time.Second)
		delivery.NextRetryAt = &nextRetry
		delivery.ErrorMessage = &errMsg

		if err := s.webhookRepo.UpdateDelivery(ctx, delivery); err != nil {
			s.logger.Error("Failed to update webhook delivery", zap.Error(err))
		}

		s.logger.Warn("Webhook delivery failed, will retry",
			zap.String("endpoint", endpoint.Name),
			zap.Int("attempt", delivery.Attempts),
			zap.Int("max_retries", endpoint.MaxRetries),
			zap.Time("next_retry", nextRetry),
			zap.String("error", errMsg),
		)
	}
}

func (s *WebhookService) markDeliveryFailed(ctx context.Context, delivery *domain.WebhookDelivery, errMsg string) {
	delivery.Status = domain.WebhookDeliveryFailed
	delivery.NextRetryAt = nil
	delivery.ErrorMessage = &errMsg

	if err := s.webhookRepo.UpdateDelivery(ctx, delivery); err != nil {
		s.logger.Error("Failed to update webhook delivery", zap.Error(err))
	}

	s.logger.Error("Webhook delivery permanently failed",
		zap.String("delivery_id", delivery.ID.String()),
		zap.Int("attempts", delivery.Attempts),
		zap.String("error", errMsg),
	)
}

// Background worker to process pending deliveries
func (s *WebhookService) ProcessPendingDeliveries(ctx context.Context) error {
	deliveries, err := s.webhookRepo.GetPendingDeliveries(ctx, 100)
	if err != nil {
		return err
	}

	s.logger.Debug("Processing pending webhook deliveries", zap.Int("count", len(deliveries)))

	for _, delivery := range deliveries {
		endpoint, err := s.webhookRepo.GetEndpointByID(ctx, delivery.WebhookEndpointID)
		if err != nil {
			s.logger.Error("Failed to get webhook endpoint",
				zap.String("endpoint_id", delivery.WebhookEndpointID.String()),
				zap.Error(err),
			)
			continue
		}

		if !endpoint.Enabled {
			s.logger.Debug("Skipping delivery for disabled endpoint",
				zap.String("endpoint", endpoint.Name),
			)
			continue
		}

		payload := &domain.WebhookPayload{
			EventType:      delivery.EventType,
			EventID:        delivery.ID.String(),
			OrganizationID: delivery.OrganizationID.String(),
			Timestamp:      delivery.CreatedAt,
			Data:           delivery.Payload,
		}

		s.deliverWebhook(ctx, endpoint, delivery, payload)
	}

	return nil
}

// Delivery logs

func (s *WebhookService) ListDeliveries(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.WebhookDelivery, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.webhookRepo.ListDeliveries(ctx, orgID, limit, offset)
}

// Incoming Webhooks

func (s *WebhookService) CreateIncomingToken(ctx context.Context, orgID uuid.UUID, req *domain.CreateIncomingWebhookTokenRequest) (*domain.IncomingWebhookToken, error) {
	// Generate a secure token
	tokenStr, err := generateSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	defaultPriority := "P3"
	if req.DefaultPriority != nil {
		defaultPriority = *req.DefaultPriority
	}

	token := &domain.IncomingWebhookToken{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		Name:            req.Name,
		Token:           tokenStr,
		Enabled:         true,
		IntegrationType: req.IntegrationType,
		DefaultPriority: defaultPriority,
		DefaultTags:     req.DefaultTags,
		RequestCount:    0,
	}

	if token.DefaultTags == nil {
		token.DefaultTags = []string{}
	}

	if err := s.webhookRepo.CreateIncomingToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *WebhookService) GetIncomingTokenByToken(ctx context.Context, token string) (*domain.IncomingWebhookToken, error) {
	return s.webhookRepo.GetIncomingTokenByToken(ctx, token)
}

func (s *WebhookService) ListIncomingTokens(ctx context.Context, orgID uuid.UUID) ([]*domain.IncomingWebhookToken, error) {
	return s.webhookRepo.ListIncomingTokens(ctx, orgID)
}

func (s *WebhookService) UpdateIncomingTokenUsage(ctx context.Context, tokenID uuid.UUID) error {
	return s.webhookRepo.UpdateIncomingTokenUsage(ctx, tokenID)
}

func (s *WebhookService) DeleteIncomingToken(ctx context.Context, id, orgID uuid.UUID) error {
	return s.webhookRepo.DeleteIncomingToken(ctx, id, orgID)
}

// Helper functions

func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateHMACSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

func getIntOrDefault(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}

// VerifyWebhookSignature verifies the HMAC signature of an incoming webhook
func VerifyWebhookSignature(payload []byte, signature, secret string) bool {
	expectedSignature := generateHMACSignature(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
