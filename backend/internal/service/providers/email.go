package providers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/resend/resend-go/v2"
)

// EmailConfig represents the configuration for the email provider
// Supports both SMTP (for development) and Resend (for production)
type EmailConfig struct {
	// Provider type: "smtp" or "resend"
	Provider string `json:"provider"`

	// SMTP settings (used when provider is "smtp")
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	UseTLS       bool   `json:"use_tls"`

	// Resend settings (used when provider is "resend")
	ResendAPIKey string `json:"resend_api_key"`

	// Common settings
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
}

// EmailProvider implements the NotificationProvider interface for email
type EmailProvider struct {
	config EmailConfig
}

// NewEmailProvider creates a new email notification provider
func NewEmailProvider(config *EmailConfig) *EmailProvider {
	return &EmailProvider{
		config: *config,
	}
}

// Send sends an email notification
func (p *EmailProvider) Send(recipient, subject, message string) error {
	// Validate recipient is a valid email
	if !strings.Contains(recipient, "@") {
		return fmt.Errorf("invalid email address: %s", recipient)
	}

	// Route to the appropriate sender based on provider
	if p.config.Provider == "resend" {
		return p.sendViaResend(recipient, subject, message)
	}

	// Default to SMTP
	return p.sendViaSMTP(recipient, subject, message)
}

// sendViaResend sends email using the Resend API
func (p *EmailProvider) sendViaResend(recipient, subject, message string) error {
	client := resend.NewClient(p.config.ResendAPIKey)

	// Build the from address
	from := p.config.FromAddress
	if p.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", p.config.FromName, p.config.FromAddress)
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{recipient},
		Subject: subject,
		Text:    message,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}

	return nil
}

// sendViaSMTP sends email using SMTP (for development with Mailpit)
func (p *EmailProvider) sendViaSMTP(recipient, subject, message string) error {
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

	addr := fmt.Sprintf("%s:%d", p.config.SMTPHost, p.config.SMTPPort)

	// Setup authentication (optional for Mailpit)
	var auth smtp.Auth
	if p.config.SMTPUsername != "" && p.config.SMTPPassword != "" {
		auth = smtp.PlainAuth("", p.config.SMTPUsername, p.config.SMTPPassword, p.config.SMTPHost)
	}

	// Send with or without TLS
	if p.config.UseTLS {
		return p.sendWithTLS(addr, auth, p.config.FromAddress, []string{recipient}, []byte(emailMessage))
	}

	return smtp.SendMail(addr, auth, p.config.FromAddress, []string{recipient}, []byte(emailMessage))
}

func (p *EmailProvider) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Try direct TLS connection first
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: p.config.SMTPHost,
	})
	if err != nil {
		// Fallback to STARTTLS
		return p.sendWithSTARTTLS(addr, auth, from, to, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP mail command failed: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("SMTP rcpt command failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP data command failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}

func (p *EmailProvider) sendWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Try STARTTLS if available
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: p.config.SMTPHost}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP mail command failed: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("SMTP rcpt command failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP data command failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}

// ValidateConfig validates the email provider configuration
func (p *EmailProvider) ValidateConfig(config json.RawMessage) error {
	var emailConfig EmailConfig
	if err := json.Unmarshal(config, &emailConfig); err != nil {
		return fmt.Errorf("invalid configuration format: %w", err)
	}

	// Validate common required fields
	if emailConfig.FromAddress == "" {
		return fmt.Errorf("from_address is required")
	}

	if !strings.Contains(emailConfig.FromAddress, "@") {
		return fmt.Errorf("from_address must be a valid email address")
	}

	// Validate provider-specific fields
	provider := emailConfig.Provider
	if provider == "" {
		provider = "smtp" // Default to SMTP
	}

	if provider == "resend" {
		if emailConfig.ResendAPIKey == "" {
			return fmt.Errorf("resend_api_key is required when using Resend provider")
		}
	} else {
		// SMTP validation
		if emailConfig.SMTPHost == "" {
			return fmt.Errorf("smtp_host is required when using SMTP provider")
		}

		if emailConfig.SMTPPort == 0 {
			return fmt.Errorf("smtp_port is required when using SMTP provider")
		}
	}

	return nil
}
