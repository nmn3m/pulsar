package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/pulsar/backend/internal/domain"
	"github.com/yourusername/pulsar/backend/internal/service/providers"
)

type NotificationRepository interface {
	CreateChannel(ctx context.Context, channel *domain.NotificationChannel) error
	GetChannelByID(ctx context.Context, id uuid.UUID) (*domain.NotificationChannel, error)
	ListChannels(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationChannel, error)
	UpdateChannel(ctx context.Context, channel *domain.NotificationChannel) error
	DeleteChannel(ctx context.Context, id uuid.UUID) error

	CreatePreference(ctx context.Context, pref *domain.UserNotificationPreference) error
	GetPreferenceByID(ctx context.Context, id uuid.UUID) (*domain.UserNotificationPreference, error)
	GetPreferenceByUserAndChannel(ctx context.Context, userID, channelID uuid.UUID) (*domain.UserNotificationPreference, error)
	ListPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserNotificationPreference, error)
	UpdatePreference(ctx context.Context, pref *domain.UserNotificationPreference) error
	DeletePreference(ctx context.Context, id uuid.UUID) error

	CreateLog(ctx context.Context, log *domain.NotificationLog) error
	GetLogByID(ctx context.Context, id uuid.UUID) (*domain.NotificationLog, error)
	ListLogs(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error)
	ListLogsByAlert(ctx context.Context, alertID uuid.UUID) ([]domain.NotificationLog, error)
	ListLogsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error)
	GetPendingNotifications(ctx context.Context, limit int) ([]domain.NotificationLog, error)
	UpdateLogStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus, errorMsg *string) error
	IsUserInDND(ctx context.Context, userID, channelID uuid.UUID) (bool, error)
}

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

// createProviderFromChannel creates a provider instance from a channel's configuration
func (s *NotificationService) createProviderFromChannel(channel *domain.NotificationChannel) (domain.NotificationProvider, error) {
	switch channel.ChannelType {
	case domain.ChannelTypeEmail:
		var config providers.EmailConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal email config: %w", err)
		}
		return providers.NewEmailProvider(config), nil

	case domain.ChannelTypeSlack:
		var config providers.SlackConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal slack config: %w", err)
		}
		return providers.NewSlackProvider(config), nil

	case domain.ChannelTypeTeams:
		var config providers.TeamsConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal teams config: %w", err)
		}
		return providers.NewTeamsProvider(config), nil

	case domain.ChannelTypeWebhook:
		var config providers.WebhookConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal webhook config: %w", err)
		}
		return providers.NewWebhookProvider(config), nil

	default:
		return nil, fmt.Errorf("unsupported channel type: %s", channel.ChannelType)
	}
}

// validateChannelConfig validates a channel configuration for a specific type
func (s *NotificationService) validateChannelConfig(channelType domain.ChannelType, config json.RawMessage) error {
	switch channelType {
	case domain.ChannelTypeEmail:
		provider := &providers.EmailProvider{}
		return provider.ValidateConfig(config)

	case domain.ChannelTypeSlack:
		provider := &providers.SlackProvider{}
		return provider.ValidateConfig(config)

	case domain.ChannelTypeTeams:
		provider := &providers.TeamsProvider{}
		return provider.ValidateConfig(config)

	case domain.ChannelTypeWebhook:
		provider := &providers.WebhookProvider{}
		return provider.ValidateConfig(config)

	default:
		return fmt.Errorf("unsupported channel type: %s", channelType)
	}
}

// ==================== Channel Management ====================

func (s *NotificationService) CreateChannel(ctx context.Context, orgID uuid.UUID, req *domain.CreateNotificationChannelRequest) (*domain.NotificationChannel, error) {
	// Validate the provider configuration
	if err := s.validateChannelConfig(req.ChannelType, req.Config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	channel := &domain.NotificationChannel{
		OrganizationID: orgID,
		Name:           req.Name,
		ChannelType:    req.ChannelType,
		IsEnabled:      req.IsEnabled,
		Config:         req.Config,
	}

	if err := s.repo.CreateChannel(ctx, channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *NotificationService) GetChannel(ctx context.Context, id uuid.UUID) (*domain.NotificationChannel, error) {
	return s.repo.GetChannelByID(ctx, id)
}

func (s *NotificationService) ListChannels(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationChannel, error) {
	return s.repo.ListChannels(ctx, orgID)
}

func (s *NotificationService) UpdateChannel(ctx context.Context, id uuid.UUID, req *domain.UpdateNotificationChannelRequest) (*domain.NotificationChannel, error) {
	channel, err := s.repo.GetChannelByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		channel.Name = *req.Name
	}

	if req.ChannelType != nil {
		channel.ChannelType = *req.ChannelType
	}

	if req.IsEnabled != nil {
		channel.IsEnabled = *req.IsEnabled
	}

	if req.Config != nil {
		// Validate the new configuration
		if err := s.validateChannelConfig(channel.ChannelType, req.Config); err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}
		channel.Config = req.Config
	}

	if err := s.repo.UpdateChannel(ctx, channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *NotificationService) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteChannel(ctx, id)
}

// ==================== User Preference Management ====================

func (s *NotificationService) CreatePreference(ctx context.Context, userID uuid.UUID, req *domain.CreateUserNotificationPreferenceRequest) (*domain.UserNotificationPreference, error) {
	// Validate channel exists
	if _, err := s.repo.GetChannelByID(ctx, req.ChannelID); err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	pref := &domain.UserNotificationPreference{
		UserID:      userID,
		ChannelID:   req.ChannelID,
		IsEnabled:   req.IsEnabled,
		DNDEnabled:  req.DNDEnabled,
		MinPriority: req.MinPriority,
	}

	// Parse DND times if provided
	if req.DNDStartTime != nil {
		t, err := time.Parse("15:04:05", *req.DNDStartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid dnd_start_time format: %w", err)
		}
		pref.DNDStartTime = &t
	}

	if req.DNDEndTime != nil {
		t, err := time.Parse("15:04:05", *req.DNDEndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid dnd_end_time format: %w", err)
		}
		pref.DNDEndTime = &t
	}

	if err := s.repo.CreatePreference(ctx, pref); err != nil {
		return nil, err
	}

	return pref, nil
}

func (s *NotificationService) GetPreference(ctx context.Context, id uuid.UUID) (*domain.UserNotificationPreference, error) {
	return s.repo.GetPreferenceByID(ctx, id)
}

func (s *NotificationService) ListUserPreferences(ctx context.Context, userID uuid.UUID) ([]domain.UserNotificationPreference, error) {
	return s.repo.ListPreferencesByUser(ctx, userID)
}

func (s *NotificationService) UpdatePreference(ctx context.Context, id uuid.UUID, req *domain.UpdateUserNotificationPreferenceRequest) (*domain.UserNotificationPreference, error) {
	pref, err := s.repo.GetPreferenceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.IsEnabled != nil {
		pref.IsEnabled = *req.IsEnabled
	}

	if req.DNDEnabled != nil {
		pref.DNDEnabled = *req.DNDEnabled
	}

	if req.MinPriority != nil {
		pref.MinPriority = req.MinPriority
	}

	if req.DNDStartTime != nil {
		t, err := time.Parse("15:04:05", *req.DNDStartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid dnd_start_time format: %w", err)
		}
		pref.DNDStartTime = &t
	}

	if req.DNDEndTime != nil {
		t, err := time.Parse("15:04:05", *req.DNDEndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid dnd_end_time format: %w", err)
		}
		pref.DNDEndTime = &t
	}

	if err := s.repo.UpdatePreference(ctx, pref); err != nil {
		return nil, err
	}

	return pref, nil
}

func (s *NotificationService) DeletePreference(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePreference(ctx, id)
}

// ==================== Notification Sending ====================

func (s *NotificationService) SendNotification(ctx context.Context, orgID uuid.UUID, req *domain.SendNotificationRequest) (*domain.NotificationLog, error) {
	// Get the channel
	channel, err := s.repo.GetChannelByID(ctx, req.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %w", err)
	}

	// Check if channel is enabled
	if !channel.IsEnabled {
		return nil, fmt.Errorf("channel is disabled")
	}

	// Check if channel belongs to the organization
	if channel.OrganizationID != orgID {
		return nil, domain.ErrUnauthorized
	}

	// If userID is provided, check user preferences and DND
	if req.UserID != nil {
		// Check if user is in DND mode
		inDND, err := s.repo.IsUserInDND(ctx, *req.UserID, req.ChannelID)
		if err != nil {
			return nil, fmt.Errorf("failed to check DND status: %w", err)
		}
		if inDND {
			return nil, fmt.Errorf("user is in do not disturb mode")
		}

		// Check user preferences
		pref, err := s.repo.GetPreferenceByUserAndChannel(ctx, *req.UserID, req.ChannelID)
		if err == nil && !pref.IsEnabled {
			return nil, fmt.Errorf("user has disabled notifications for this channel")
		}
	}

	// Create the notification log
	log := &domain.NotificationLog{
		OrganizationID: orgID,
		ChannelID:      req.ChannelID,
		UserID:         req.UserID,
		AlertID:        req.AlertID,
		Recipient:      req.Recipient,
		Subject:        req.Subject,
		Message:        req.Message,
		Status:         domain.NotificationStatusPending,
	}

	if err := s.repo.CreateLog(ctx, log); err != nil {
		return nil, fmt.Errorf("failed to create notification log: %w", err)
	}

	// Create provider from channel configuration
	provider, err := s.createProviderFromChannel(channel)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create provider: %v", err)
		s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusFailed, &errMsg)
		return log, fmt.Errorf(errMsg)
	}

	// Send the notification
	subject := ""
	if req.Subject != nil {
		subject = *req.Subject
	}

	err = provider.Send(req.Recipient, subject, req.Message)
	if err != nil {
		errMsg := err.Error()
		s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusFailed, &errMsg)
		return log, fmt.Errorf("failed to send notification: %w", err)
	}

	// Update log status to sent
	if err := s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusSent, nil); err != nil {
		return log, fmt.Errorf("notification sent but failed to update log: %w", err)
	}

	// Refresh the log to get updated status
	log, _ = s.repo.GetLogByID(ctx, log.ID)

	return log, nil
}

// ProcessPendingNotifications processes pending notifications (for background workers)
func (s *NotificationService) ProcessPendingNotifications(ctx context.Context, limit int) error {
	logs, err := s.repo.GetPendingNotifications(ctx, limit)
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %w", err)
	}

	for _, log := range logs {
		// Get the channel
		channel, err := s.repo.GetChannelByID(ctx, log.ChannelID)
		if err != nil {
			errMsg := fmt.Sprintf("channel not found: %v", err)
			s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusFailed, &errMsg)
			continue
		}

		// Skip if channel is disabled
		if !channel.IsEnabled {
			errMsg := "channel is disabled"
			s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusFailed, &errMsg)
			continue
		}

		// Create provider from channel configuration
		provider, err := s.createProviderFromChannel(channel)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create provider: %v", err)
			s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusFailed, &errMsg)
			continue
		}

		// Send the notification
		subject := ""
		if log.Subject != nil {
			subject = *log.Subject
		}

		err = provider.Send(log.Recipient, subject, log.Message)
		if err != nil {
			errMsg := err.Error()
			s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusFailed, &errMsg)
			continue
		}

		// Update status to sent
		s.repo.UpdateLogStatus(ctx, log.ID, domain.NotificationStatusSent, nil)
	}

	return nil
}

// ==================== Notification Logs ====================

func (s *NotificationService) GetLog(ctx context.Context, id uuid.UUID) (*domain.NotificationLog, error) {
	return s.repo.GetLogByID(ctx, id)
}

func (s *NotificationService) ListLogs(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListLogs(ctx, orgID, limit, offset)
}

func (s *NotificationService) ListLogsByAlert(ctx context.Context, alertID uuid.UUID) ([]domain.NotificationLog, error) {
	return s.repo.ListLogsByAlert(ctx, alertID)
}

func (s *NotificationService) ListLogsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListLogsByUser(ctx, userID, limit, offset)
}

// Helper function to unmarshal channel config into a specific struct
func UnmarshalChannelConfig(config json.RawMessage, target interface{}) error {
	return json.Unmarshal(config, target)
}
