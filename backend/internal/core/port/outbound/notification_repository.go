package outbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
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
