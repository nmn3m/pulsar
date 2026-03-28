package inbound

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, orgID uuid.UUID, req *dto.CreateScheduleRequest) (*domain.Schedule, error)
	GetSchedule(ctx context.Context, id uuid.UUID) (*domain.Schedule, error)
	GetScheduleWithRotations(ctx context.Context, id uuid.UUID) (*domain.ScheduleWithRotations, error)
	UpdateSchedule(ctx context.Context, id uuid.UUID, req *dto.UpdateScheduleRequest) (*domain.Schedule, error)
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
	ListSchedules(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.Schedule, error)
	CreateRotation(ctx context.Context, scheduleID uuid.UUID, req *dto.CreateRotationRequest) (*domain.ScheduleRotation, error)
	GetRotation(ctx context.Context, id uuid.UUID) (*domain.ScheduleRotation, error)
	UpdateRotation(ctx context.Context, id uuid.UUID, req *dto.UpdateRotationRequest) (*domain.ScheduleRotation, error)
	DeleteRotation(ctx context.Context, id uuid.UUID) error
	ListRotations(ctx context.Context, scheduleID uuid.UUID) ([]*domain.ScheduleRotation, error)
	AddParticipant(ctx context.Context, rotationID uuid.UUID, req *dto.AddParticipantRequest) (*domain.ScheduleRotationParticipant, error)
	RemoveParticipant(ctx context.Context, rotationID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, rotationID uuid.UUID) ([]*domain.ParticipantWithUser, error)
	ReorderParticipants(ctx context.Context, rotationID uuid.UUID, req *dto.ReorderParticipantsRequest) error
	CreateOverride(ctx context.Context, scheduleID uuid.UUID, req *dto.CreateOverrideRequest) (*domain.ScheduleOverride, error)
	GetOverride(ctx context.Context, id uuid.UUID) (*domain.ScheduleOverride, error)
	UpdateOverride(ctx context.Context, id uuid.UUID, req *dto.UpdateOverrideRequest) (*domain.ScheduleOverride, error)
	DeleteOverride(ctx context.Context, id uuid.UUID) error
	ListOverrides(ctx context.Context, scheduleID uuid.UUID, start, end time.Time) ([]*domain.ScheduleOverride, error)
	GetOnCallUser(ctx context.Context, scheduleID uuid.UUID, at time.Time) (*domain.OnCallUser, error)
}
