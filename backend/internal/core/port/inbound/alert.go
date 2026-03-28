package inbound

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type AlertService interface {
	CreateAlert(ctx context.Context, orgID uuid.UUID, req *dto.CreateAlertRequest) (*domain.Alert, error)
	GetAlert(ctx context.Context, id, orgID uuid.UUID) (*domain.Alert, error)
	UpdateAlert(ctx context.Context, id, orgID uuid.UUID, req *dto.UpdateAlertRequest) (*domain.Alert, error)
	DeleteAlert(ctx context.Context, id, orgID uuid.UUID) error
	ListAlerts(ctx context.Context, orgID uuid.UUID, req *dto.ListAlertsRequest) (*dto.ListAlertsResponse, error)
	AcknowledgeAlert(ctx context.Context, id, orgID, userID uuid.UUID) error
	CloseAlert(ctx context.Context, id, orgID, userID uuid.UUID, reason string) error
	SnoozeAlert(ctx context.Context, id, orgID uuid.UUID, until time.Time) error
	AssignAlert(ctx context.Context, id, orgID uuid.UUID, userID, teamID *uuid.UUID) error
}
