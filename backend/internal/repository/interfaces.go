package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Organization, error)
	Update(ctx context.Context, org *domain.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Organization, error)

	// Organization user methods
	AddUser(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error
	RemoveUser(ctx context.Context, orgID, userID uuid.UUID) error
	GetUserRole(ctx context.Context, orgID, userID uuid.UUID) (domain.UserRole, error)
	UpdateUserRole(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error
	ListUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error)
	ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]*domain.Organization, error)
}

type AlertRepository interface {
	Create(ctx context.Context, alert *domain.Alert) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error)
	Update(ctx context.Context, alert *domain.Alert) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *domain.AlertFilter) ([]*domain.Alert, int, error)

	// Alert actions
	Acknowledge(ctx context.Context, id, userID uuid.UUID) error
	Close(ctx context.Context, id, userID uuid.UUID, reason string) error
	Snooze(ctx context.Context, id uuid.UUID, until time.Time) error
	Assign(ctx context.Context, id uuid.UUID, userID, teamID *uuid.UUID) error

	// Deduplication
	FindByDedupKey(ctx context.Context, orgID uuid.UUID, dedupKey string) (*domain.Alert, error)
	IncrementDedupCount(ctx context.Context, id uuid.UUID) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Team, error)

	// Team member methods
	AddMember(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]*domain.UserWithTeamRole, error)
	ListUserTeams(ctx context.Context, userID uuid.UUID) ([]*domain.Team, error)
}

type ScheduleRepository interface {
	// Schedule CRUD
	Create(ctx context.Context, schedule *domain.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error)
	Update(ctx context.Context, schedule *domain.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Schedule, error)
	GetWithRotations(ctx context.Context, id uuid.UUID) (*domain.ScheduleWithRotations, error)

	// Rotation CRUD
	CreateRotation(ctx context.Context, rotation *domain.ScheduleRotation) error
	GetRotation(ctx context.Context, id uuid.UUID) (*domain.ScheduleRotation, error)
	UpdateRotation(ctx context.Context, rotation *domain.ScheduleRotation) error
	DeleteRotation(ctx context.Context, id uuid.UUID) error
	ListRotations(ctx context.Context, scheduleID uuid.UUID) ([]*domain.ScheduleRotation, error)

	// Rotation participants
	AddParticipant(ctx context.Context, participant *domain.ScheduleRotationParticipant) error
	RemoveParticipant(ctx context.Context, rotationID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, rotationID uuid.UUID) ([]*domain.ParticipantWithUser, error)
	ReorderParticipants(ctx context.Context, rotationID uuid.UUID, userIDs []uuid.UUID) error

	// Overrides
	CreateOverride(ctx context.Context, override *domain.ScheduleOverride) error
	GetOverride(ctx context.Context, id uuid.UUID) (*domain.ScheduleOverride, error)
	UpdateOverride(ctx context.Context, override *domain.ScheduleOverride) error
	DeleteOverride(ctx context.Context, id uuid.UUID) error
	ListOverrides(ctx context.Context, scheduleID uuid.UUID, start, end time.Time) ([]*domain.ScheduleOverride, error)

	// On-call calculation
	GetOnCallUser(ctx context.Context, scheduleID uuid.UUID, at time.Time) (*domain.OnCallUser, error)
}

type EscalationPolicyRepository interface {
	// Policy CRUD
	Create(ctx context.Context, policy *domain.EscalationPolicy) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicy, error)
	Update(ctx context.Context, policy *domain.EscalationPolicy) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.EscalationPolicy, error)
	GetWithRules(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicyWithRules, error)

	// Rule CRUD
	CreateRule(ctx context.Context, rule *domain.EscalationRule) error
	GetRule(ctx context.Context, id uuid.UUID) (*domain.EscalationRule, error)
	UpdateRule(ctx context.Context, rule *domain.EscalationRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error
	ListRules(ctx context.Context, policyID uuid.UUID) ([]*domain.EscalationRule, error)

	// Target CRUD
	AddTarget(ctx context.Context, target *domain.EscalationTarget) error
	RemoveTarget(ctx context.Context, id uuid.UUID) error
	ListTargets(ctx context.Context, ruleID uuid.UUID) ([]*domain.EscalationTarget, error)

	// Escalation events
	CreateEvent(ctx context.Context, event *domain.AlertEscalationEvent) error
	GetLatestEvent(ctx context.Context, alertID uuid.UUID) (*domain.AlertEscalationEvent, error)
	UpdateEvent(ctx context.Context, event *domain.AlertEscalationEvent) error
	ListPendingEscalations(ctx context.Context, before time.Time) ([]*domain.AlertEscalationEvent, error)
}

type RoutingRuleRepository interface {
	Create(ctx context.Context, rule *domain.AlertRoutingRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AlertRoutingRule, error)
	Update(ctx context.Context, rule *domain.AlertRoutingRule) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.AlertRoutingRule, error)
	ListEnabled(ctx context.Context, orgID uuid.UUID) ([]*domain.AlertRoutingRule, error)
	Reorder(ctx context.Context, orgID uuid.UUID, ruleIDs []uuid.UUID) error
}

type DNDSettingsRepository interface {
	Create(ctx context.Context, settings *domain.UserDNDSettings) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserDNDSettings, error)
	Update(ctx context.Context, settings *domain.UserDNDSettings) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Upsert(ctx context.Context, settings *domain.UserDNDSettings) error
}
