package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// TelnyxSMSConfig represents the configuration for the Telnyx SMS provider
type TelnyxSMSConfig struct {
	APIKey              string `json:"api_key"`
	FromNumber          string `json:"from_number"`
	MessagingProfileID  string `json:"messaging_profile_id,omitempty"`
}

// TelnyxSMSProvider implements the NotificationProvider interface for Telnyx SMS
type TelnyxSMSProvider struct {
	config TelnyxSMSConfig
}

// NewTelnyxSMSProvider creates a new Telnyx SMS notification provider
func NewTelnyxSMSProvider(config TelnyxSMSConfig) *TelnyxSMSProvider {
	return &TelnyxSMSProvider{
		config: config,
	}
}

// telnyxSMSRequest represents the request body for the Telnyx Messages API
type telnyxSMSRequest struct {
	From                string `json:"from"`
	To                  string `json:"to"`
	Text                string `json:"text"`
	MessagingProfileID  string `json:"messaging_profile_id,omitempty"`
}

// Send sends an SMS notification via Telnyx
func (p *TelnyxSMSProvider) Send(recipient, subject, message string) error {
	if !isValidE164(recipient) {
		return fmt.Errorf("invalid phone number format: must be E.164 (e.g. +1234567890)")
	}

	// Build the SMS body: prepend subject if provided
	body := message
	if subject != "" {
		body = fmt.Sprintf("[%s] %s", subject, message)
	}

	// Truncate at 1600 characters (Telnyx concatenated SMS limit)
	if len(body) > 1600 {
		body = body[:1600]
	}

	// Build the API request
	smsReq := telnyxSMSRequest{
		From: p.config.FromNumber,
		To:   recipient,
		Text: body,
	}
	if p.config.MessagingProfileID != "" {
		smsReq.MessagingProfileID = p.config.MessagingProfileID
	}

	jsonData, err := json.Marshal(smsReq)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.telnyx.com/v2/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telnyx API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// ValidateConfig validates the Telnyx SMS provider configuration
func (p *TelnyxSMSProvider) ValidateConfig(config json.RawMessage) error {
	var cfg TelnyxSMSConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return fmt.Errorf("invalid configuration format: %w", err)
	}

	if cfg.APIKey == "" {
		return fmt.Errorf("api_key is required")
	}

	if cfg.FromNumber == "" {
		return fmt.Errorf("from_number is required")
	}

	if !isValidE164(cfg.FromNumber) {
		return fmt.Errorf("from_number must be in E.164 format (e.g. +1234567890)")
	}

	return nil
}

// isValidE164 checks if a phone number is in E.164 format
func isValidE164(phone string) bool {
	re := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}
