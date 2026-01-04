package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	ejson "encoding/json"
)

// WebhookConfig represents the configuration for the generic webhook provider
type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method,omitempty"` // HTTP method (default: POST)
	Headers map[string]string `json:"headers,omitempty"`
	Timeout int               `json:"timeout,omitempty"` // Timeout in seconds (default: 30)
}

// WebhookPayload represents the payload sent to the webhook
type WebhookPayload struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// WebhookProvider implements the NotificationProvider interface for generic webhooks
type WebhookProvider struct {
	config WebhookConfig
}

// NewWebhookProvider creates a new generic webhook notification provider
func NewWebhookProvider(config WebhookConfig) *WebhookProvider {
	return &WebhookProvider{
		config: config,
	}
}

// Send sends a notification to a generic webhook
func (p *WebhookProvider) Send(recipient string, subject string, message string) error {
	// Build the webhook payload
	payload := WebhookPayload{
		Recipient: recipient,
		Subject:   subject,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// Determine HTTP method (default to POST)
	method := p.config.Method
	if method == "" {
		method = "POST"
	}

	// Create HTTP request
	req, err := http.NewRequest(method, p.config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set default content type
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range p.config.Headers {
		req.Header.Set(key, value)
	}

	// Set timeout (default to 30 seconds)
	timeout := time.Duration(p.config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	client := &http.Client{
		Timeout: timeout,
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ValidateConfig validates the generic webhook provider configuration
func (p *WebhookProvider) ValidateConfig(config json.RawMessage) error {
	var webhookConfig WebhookConfig
	if err := ejson.Unmarshal(config, &webhookConfig); err != nil {
		return fmt.Errorf("invalid configuration format: %w", err)
	}

	// Validate required fields
	if webhookConfig.URL == "" {
		return fmt.Errorf("url is required")
	}

	// Validate URL format (should be HTTP or HTTPS)
	if len(webhookConfig.URL) < 8 {
		return fmt.Errorf("url must be a valid HTTP or HTTPS URL")
	}

	if webhookConfig.URL[:7] != "http://" && webhookConfig.URL[:8] != "https://" {
		return fmt.Errorf("url must start with http:// or https://")
	}

	// Validate HTTP method if provided
	if webhookConfig.Method != "" {
		validMethods := map[string]bool{
			"GET":     true,
			"POST":    true,
			"PUT":     true,
			"PATCH":   true,
			"DELETE":  true,
			"HEAD":    true,
			"OPTIONS": true,
		}
		if !validMethods[webhookConfig.Method] {
			return fmt.Errorf("invalid HTTP method: %s", webhookConfig.Method)
		}
	}

	// Validate timeout if provided
	if webhookConfig.Timeout < 0 || webhookConfig.Timeout > 300 {
		return fmt.Errorf("timeout must be between 0 and 300 seconds")
	}

	return nil
}
