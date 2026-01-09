package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type MetricsService struct {
	metricsRepo repository.MetricsRepository
}

func NewMetricsService(metricsRepo repository.MetricsRepository) *MetricsService {
	return &MetricsService{
		metricsRepo: metricsRepo,
	}
}

// GetDashboardMetrics returns aggregated metrics for the dashboard
func (s *MetricsService) GetDashboardMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.DashboardMetrics, error) {
	alertMetrics, err := s.metricsRepo.GetAlertMetrics(ctx, orgID, filter)
	if err != nil {
		return nil, err
	}

	incidentMetrics, err := s.metricsRepo.GetIncidentMetrics(ctx, orgID, filter)
	if err != nil {
		return nil, err
	}

	notificationMetrics, err := s.metricsRepo.GetNotificationMetrics(ctx, orgID, filter)
	if err != nil {
		return nil, err
	}

	alertTrend, err := s.metricsRepo.GetAlertTrend(ctx, orgID, filter)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardMetrics{
		Alerts:        alertMetrics,
		Incidents:     incidentMetrics,
		Notifications: notificationMetrics,
		AlertTrend:    alertTrend,
	}, nil
}

// GetAlertMetrics returns alert-specific metrics
func (s *MetricsService) GetAlertMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.AlertMetrics, error) {
	return s.metricsRepo.GetAlertMetrics(ctx, orgID, filter)
}

// GetIncidentMetrics returns incident-specific metrics
func (s *MetricsService) GetIncidentMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.IncidentMetrics, error) {
	return s.metricsRepo.GetIncidentMetrics(ctx, orgID, filter)
}

// GetNotificationMetrics returns notification-specific metrics
func (s *MetricsService) GetNotificationMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.NotificationMetrics, error) {
	return s.metricsRepo.GetNotificationMetrics(ctx, orgID, filter)
}

// GetAlertTrend returns time-series data for alerts
func (s *MetricsService) GetAlertTrend(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.AlertTrend, error) {
	return s.metricsRepo.GetAlertTrend(ctx, orgID, filter)
}

// GetTeamMetrics returns team performance metrics
func (s *MetricsService) GetTeamMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) ([]domain.TeamMetrics, error) {
	return s.metricsRepo.GetTeamMetrics(ctx, orgID, filter)
}
