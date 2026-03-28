package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/resend/resend-go/v2"

	"github.com/nmn3m/pulsar/backend/internal/config"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(msg *EmailMessage) error
}

// EmailService handles email operations using configurable providers
type EmailService struct {
	config *config.EmailConfig
	smtp   *config.SMTPConfig
	sender EmailSender
}

// NewEmailService creates a new email service with the appropriate provider
func NewEmailService(emailCfg *config.EmailConfig, smtpCfg *config.SMTPConfig) *EmailService {
	svc := &EmailService{
		config: emailCfg,
		smtp:   smtpCfg,
	}

	// Initialize the appropriate sender based on provider
	if emailCfg.Provider == "resend" && emailCfg.ResendAPIKey != "" {
		svc.sender = NewResendSender(emailCfg)
	} else {
		// Default to SMTP (for development with Mailpit)
		svc.sender = NewSMTPSender(smtpCfg)
	}

	return svc
}

func (s *EmailService) IsConfigured() bool {
	if s.config.Provider == "resend" {
		return s.config.Enabled && s.config.ResendAPIKey != ""
	}
	// SMTP mode - check SMTP config
	return s.smtp.Enabled && s.smtp.Host != ""
}

func (s *EmailService) GetProvider() string {
	return s.config.Provider
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func (s *EmailService) Send(msg *EmailMessage) error {
	if !s.IsConfigured() {
		return fmt.Errorf("email service is not configured")
	}
	return s.sender.Send(msg)
}

// =============================================================================
// SMTP Sender (for development with Mailpit)
// =============================================================================

type SMTPSender struct {
	config *config.SMTPConfig
}

func NewSMTPSender(cfg *config.SMTPConfig) *SMTPSender {
	return &SMTPSender{config: cfg}
}

func (s *SMTPSender) Send(msg *EmailMessage) error {
	// Build the email message
	from := s.config.From
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	}

	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(msg.To, ", ")
	headers["Subject"] = msg.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("%s; charset=\"UTF-8\"", contentType)

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(msg.Body)

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	var auth smtp.Auth
	if s.config.Username != "" && s.config.Password != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, s.config.From, msg.To, []byte(message.String()))
	}

	return smtp.SendMail(addr, auth, s.config.From, msg.To, []byte(message.String()))
}

func (s *SMTPSender) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Connect to the server
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: s.config.Host,
	})
	if err != nil {
		// Try STARTTLS instead
		return s.sendWithSTARTTLS(addr, auth, from, to, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
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

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}

func (s *SMTPSender) sendWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Try STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: s.config.Host}
		if err = client.StartTLS(config); err != nil {
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

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close email writer: %w", err)
	}

	return client.Quit()
}

// =============================================================================
// Resend Sender (for production)
// =============================================================================

type ResendSender struct {
	client *resend.Client
	config *config.EmailConfig
}

func NewResendSender(cfg *config.EmailConfig) *ResendSender {
	return &ResendSender{
		client: resend.NewClient(cfg.ResendAPIKey),
		config: cfg,
	}
}

func (s *ResendSender) Send(msg *EmailMessage) error {
	// Build the from address
	from := s.config.From
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.From)
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      msg.To,
		Subject: msg.Subject,
	}

	if msg.IsHTML {
		params.Html = msg.Body
	} else {
		params.Text = msg.Body
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}

	return nil
}

// =============================================================================
// Email Templates (used by auth and team services)
// =============================================================================

// SendTeamInvitation sends a team invitation email
func (s *EmailService) SendTeamInvitation(ctx context.Context, toEmail, teamName, inviterName, inviteToken string) error {
	subject := fmt.Sprintf("You've been invited to join %s - Pulsar", teamName)

	// TODO: Make this URL configurable
	inviteURL := fmt.Sprintf("http://localhost:5173/invitations/accept?token=%s", inviteToken)

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px;">
    <div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow: hidden;">
        <div style="background: linear-gradient(135deg, #6366f1 0%%, #8b5cf6 100%%); padding: 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Pulsar</h1>
        </div>
        <div style="padding: 40px 30px; text-align: center;">
            <h2 style="color: #1f2937; margin: 0 0 10px 0; font-size: 22px;">You're invited!</h2>
            <p style="color: #6b7280; margin: 0 0 30px 0; font-size: 15px;">
                <strong>%s</strong> has invited you to join the team <strong>%s</strong> on Pulsar.
            </p>
            <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #6366f1 0%%, #8b5cf6 100%%); color: #ffffff; text-decoration: none; padding: 14px 32px; border-radius: 8px; font-weight: 600; font-size: 16px;">
                Accept Invitation
            </a>
            <p style="color: #9ca3af; font-size: 13px; margin: 30px 0 0 0;">This invitation expires in 7 days.</p>
        </div>
        <div style="background-color: #f9fafb; padding: 20px; text-align: center; border-top: 1px solid #e5e7eb;">
            <p style="color: #9ca3af; font-size: 12px; margin: 0;">If you didn't expect this invitation, you can safely ignore this email.</p>
        </div>
    </div>
</body>
</html>`, inviterName, teamName, inviteURL)

	return s.Send(&EmailMessage{
		To:      []string{toEmail},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

// SendOTPEmail sends an OTP verification email
func (s *EmailService) SendOTPEmail(to, otp, username string) error {
	subject := "Verify your email - Pulsar"

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px;">
    <div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow: hidden;">
        <div style="background: linear-gradient(135deg, #6366f1 0%%, #8b5cf6 100%%); padding: 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700;">Pulsar</h1>
        </div>
        <div style="padding: 40px 30px; text-align: center;">
            <h2 style="color: #1f2937; margin: 0 0 10px 0; font-size: 22px;">Verify your email</h2>
            <p style="color: #6b7280; margin: 0 0 30px 0; font-size: 15px;">Hi %s, use the code below to verify your email address:</p>
            <div style="background-color: #f3f4f6; border-radius: 8px; padding: 20px; margin-bottom: 30px;">
                <span style="font-size: 36px; font-weight: 700; letter-spacing: 8px; color: #1f2937;">%s</span>
            </div>
            <p style="color: #9ca3af; font-size: 13px; margin: 0;">This code expires in 10 minutes.</p>
        </div>
        <div style="background-color: #f9fafb; padding: 20px; text-align: center; border-top: 1px solid #e5e7eb;">
            <p style="color: #9ca3af; font-size: 12px; margin: 0;">If you didn't request this code, you can safely ignore this email.</p>
        </div>
    </div>
</body>
</html>`, username, otp)

	return s.Send(&EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}
