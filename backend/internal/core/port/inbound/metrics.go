package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type MetricsService interface {
	GetDashboardMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.DashboardMetrics, error)
	GetAlertMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.AlertMetrics, error)
	GetIncidentMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.IncidentMetrics, error)
	GetNotificationMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.NotificationMetrics, error)
	GetAlertTrend(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.AlertTrend, error)
	GetTeamMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) ([]domain.TeamMetrics, error)
}
