package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TeamsConfig represents the configuration for the Microsoft Teams provider
type TeamsConfig struct {
	WebhookURL string `json:"webhook_url"`
	ThemeColor string `json:"theme_color,omitempty"` // Optional: hex color for the message card
}

// TeamsMessageCard represents a Microsoft Teams MessageCard payload
// Using the MessageCard format for compatibility
type TeamsMessageCard struct {
	Type       string `json:"@type"`
	Context    string `json:"@context"`
	Summary    string `json:"summary,omitempty"`
	ThemeColor string `json:"themeColor,omitempty"`
	Title      string `json:"title,omitempty"`
	Text       string `json:"text"`
}

// TeamsProvider implements the NotificationProvider interface for Microsoft Teams
type TeamsProvider struct {
	config TeamsConfig
}

// NewTeamsProvider creates a new Microsoft Teams notification provider
func NewTeamsProvider(config TeamsConfig) *TeamsProvider {
	return &TeamsProvider{
		config: config,
	}
}

// Send sends a Microsoft Teams notification
func (p *TeamsProvider) Send(recipient, subject, message string) error {
	// Build the Teams MessageCard payload
	payload := TeamsMessageCard{
		Type:    "MessageCard",
		Context: "https://schema.org/extensions",
		Text:    message,
	}

	// Add subject as title if provided
	if subject != "" {
		payload.Title = subject
		payload.Summary = subject
	} else {
		payload.Summary = "Notification from Pulsar"
	}

	// Add theme color if configured
	if p.config.ThemeColor != "" {
		payload.ThemeColor = p.config.ThemeColor
	} else {
		// Default to a blue color
		payload.ThemeColor = "0078D4"
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Teams payload: %w", err)
	}

	// Send the HTTP POST request to the webhook URL
	resp, err := http.Post(p.config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send Teams webhook: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("teams API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ValidateConfig validates the Microsoft Teams provider configuration
func (p *TeamsProvider) ValidateConfig(config json.RawMessage) error {
	var teamsConfig TeamsConfig
	if err := json.Unmarshal(config, &teamsConfig); err != nil {
		return fmt.Errorf("invalid configuration format: %w", err)
	}

	// Validate required fields
	if teamsConfig.WebhookURL == "" {
		return fmt.Errorf("webhook_url is required")
	}

	// Validate webhook URL format
	if len(teamsConfig.WebhookURL) < 10 || teamsConfig.WebhookURL[:8] != "https://" {
		return fmt.Errorf("webhook_url must be a valid HTTPS URL")
	}

	// Validate theme color format if provided (should be hex color without #)
	if teamsConfig.ThemeColor != "" {
		if len(teamsConfig.ThemeColor) != 6 {
			return fmt.Errorf("theme_color must be a 6-digit hex color (without #)")
		}
		// Basic validation that it contains only hex characters
		for _, c := range teamsConfig.ThemeColor {
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
				return fmt.Errorf("theme_color must be a valid hex color")
			}
		}
	}

	return nil
}
