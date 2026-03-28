package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type NotificationService interface {
	CreateChannel(ctx context.Context, orgID uuid.UUID, req *dto.CreateNotificationChannelRequest) (*domain.NotificationChannel, error)
	GetChannel(ctx context.Context, id uuid.UUID) (*domain.NotificationChannel, error)
	ListChannels(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationChannel, error)
	UpdateChannel(ctx context.Context, id uuid.UUID, req *dto.UpdateNotificationChannelRequest) (*domain.NotificationChannel, error)
	DeleteChannel(ctx context.Context, id uuid.UUID) error
	CreatePreference(ctx context.Context, userID uuid.UUID, req *dto.CreateUserNotificationPreferenceRequest) (*domain.UserNotificationPreference, error)
	GetPreference(ctx context.Context, id uuid.UUID) (*domain.UserNotificationPreference, error)
	ListUserPreferences(ctx context.Context, userID uuid.UUID) ([]domain.UserNotificationPreference, error)
	UpdatePreference(ctx context.Context, id uuid.UUID, req *dto.UpdateUserNotificationPreferenceRequest) (*domain.UserNotificationPreference, error)
	DeletePreference(ctx context.Context, id uuid.UUID) error
	SendNotification(ctx context.Context, orgID uuid.UUID, req *dto.SendNotificationRequest) (*domain.NotificationLog, error)
	ProcessPendingNotifications(ctx context.Context, limit int) error
	GetLog(ctx context.Context, id uuid.UUID) (*domain.NotificationLog, error)
	ListLogs(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error)
	ListLogsByAlert(ctx context.Context, alertID uuid.UUID) ([]domain.NotificationLog, error)
	ListLogsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error)
}
