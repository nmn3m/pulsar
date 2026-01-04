package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ejson "encoding/json"
)

// SlackConfig represents the configuration for the Slack provider
type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel,omitempty"`    // Optional: override default channel
	Username   string `json:"username,omitempty"`   // Optional: bot username
	IconEmoji  string `json:"icon_emoji,omitempty"` // Optional: bot icon
}

// SlackMessage represents a Slack message payload
type SlackMessage struct {
	Text      string `json:"text"`
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// SlackProvider implements the NotificationProvider interface for Slack
type SlackProvider struct {
	config SlackConfig
}

// NewSlackProvider creates a new Slack notification provider
func NewSlackProvider(config SlackConfig) *SlackProvider {
	return &SlackProvider{
		config: config,
	}
}

// Send sends a Slack notification
func (p *SlackProvider) Send(recipient string, subject string, message string) error {
	// Build the full message
	fullMessage := message
	if subject != "" {
		fullMessage = fmt.Sprintf("*%s*\n%s", subject, message)
	}

	// Build the Slack message payload
	payload := SlackMessage{
		Text: fullMessage,
	}

	// Use recipient as channel override if provided and it starts with # or @
	if recipient != "" && (recipient[0] == '#' || recipient[0] == '@') {
		payload.Channel = recipient
	} else if p.config.Channel != "" {
		payload.Channel = p.config.Channel
	}

	if p.config.Username != "" {
		payload.Username = p.config.Username
	}

	if p.config.IconEmoji != "" {
		payload.IconEmoji = p.config.IconEmoji
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	// Send the HTTP POST request to the webhook URL
	resp, err := http.Post(p.config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Slack webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Slack API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ValidateConfig validates the Slack provider configuration
func (p *SlackProvider) ValidateConfig(config json.RawMessage) error {
	var slackConfig SlackConfig
	if err := ejson.Unmarshal(config, &slackConfig); err != nil {
		return fmt.Errorf("invalid configuration format: %w", err)
	}

	// Validate required fields
	if slackConfig.WebhookURL == "" {
		return fmt.Errorf("webhook_url is required")
	}

	// Validate webhook URL format
	if len(slackConfig.WebhookURL) < 10 || slackConfig.WebhookURL[:8] != "https://" {
		return fmt.Errorf("webhook_url must be a valid HTTPS URL")
	}

	return nil
}
