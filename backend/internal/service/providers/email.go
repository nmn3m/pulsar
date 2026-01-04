package providers

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"

	ejson "encoding/json"
)

// EmailConfig represents the configuration for the email provider
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	FromAddress  string `json:"from_address"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"`
}

// EmailProvider implements the NotificationProvider interface for email
type EmailProvider struct {
	config EmailConfig
}

// NewEmailProvider creates a new email notification provider
func NewEmailProvider(config EmailConfig) *EmailProvider {
	return &EmailProvider{
		config: config,
	}
}

// Send sends an email notification
func (p *EmailProvider) Send(recipient string, subject string, message string) error {
	// Validate recipient is a valid email
	if !strings.Contains(recipient, "@") {
		return fmt.Errorf("invalid email address: %s", recipient)
	}

	// Build the email message
	from := p.config.FromAddress
	if p.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", p.config.FromName, p.config.FromAddress)
	}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = recipient
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	// Build the message
	emailMessage := ""
	for k, v := range headers {
		emailMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	emailMessage += "\r\n" + message

	// Setup authentication
	auth := smtp.PlainAuth("", p.config.SMTPUsername, p.config.SMTPPassword, p.config.SMTPHost)

	// Send the email
	addr := fmt.Sprintf("%s:%d", p.config.SMTPHost, p.config.SMTPPort)
	err := smtp.SendMail(addr, auth, p.config.FromAddress, []string{recipient}, []byte(emailMessage))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// ValidateConfig validates the email provider configuration
func (p *EmailProvider) ValidateConfig(config json.RawMessage) error {
	var emailConfig EmailConfig
	if err := ejson.Unmarshal(config, &emailConfig); err != nil {
		return fmt.Errorf("invalid configuration format: %w", err)
	}

	// Validate required fields
	if emailConfig.SMTPHost == "" {
		return fmt.Errorf("smtp_host is required")
	}

	if emailConfig.SMTPPort == 0 {
		return fmt.Errorf("smtp_port is required")
	}

	if emailConfig.SMTPUsername == "" {
		return fmt.Errorf("smtp_username is required")
	}

	if emailConfig.SMTPPassword == "" {
		return fmt.Errorf("smtp_password is required")
	}

	if emailConfig.FromAddress == "" {
		return fmt.Errorf("from_address is required")
	}

	if !strings.Contains(emailConfig.FromAddress, "@") {
		return fmt.Errorf("from_address must be a valid email address")
	}

	return nil
}
