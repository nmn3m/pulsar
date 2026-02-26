package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *domain.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error)
	Update(ctx context.Context, schedule *domain.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Schedule, error)
	GetWithRotations(ctx context.Context, id uuid.UUID) (*domain.ScheduleWithRotations, error)
	CreateRotation(ctx context.Context, rotation *domain.ScheduleRotation) error
	GetRotation(ctx context.Context, id uuid.UUID) (*domain.ScheduleRotation, error)
	UpdateRotation(ctx context.Context, rotation *domain.ScheduleRotation) error
	DeleteRotation(ctx context.Context, id uuid.UUID) error
	ListRotations(ctx context.Context, scheduleID uuid.UUID) ([]*domain.ScheduleRotation, error)
	AddParticipant(ctx context.Context, participant *domain.ScheduleRotationParticipant) error
	RemoveParticipant(ctx context.Context, rotationID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, rotationID uuid.UUID) ([]*domain.ParticipantWithUser, error)
	ReorderParticipants(ctx context.Context, rotationID uuid.UUID, userIDs []uuid.UUID) error
	CreateOverride(ctx context.Context, override *domain.ScheduleOverride) error
	GetOverride(ctx context.Context, id uuid.UUID) (*domain.ScheduleOverride, error)
	UpdateOverride(ctx context.Context, override *domain.ScheduleOverride) error
	DeleteOverride(ctx context.Context, id uuid.UUID) error
	ListOverrides(ctx context.Context, scheduleID uuid.UUID, start, end time.Time) ([]*domain.ScheduleOverride, error)
	GetOnCallUser(ctx context.Context, scheduleID uuid.UUID, at time.Time) (*domain.OnCallUser, error)
}
