package postgres

import (
	"github.com/nmn3m/pulsar/backend/internal/usecase/repository"
)

// Compile-time interface checks
var (
	_ repository.AlertRepository          = (*AlertRepository)(nil)
	_ repository.APIKeyRepository         = (*apiKeyRepository)(nil)
	_ repository.DNDSettingsRepository     = (*DNDSettingsRepository)(nil)
	_ repository.EmailVerificationRepository = (*EmailVerificationRepository)(nil)
	_ repository.EscalationPolicyRepository = (*EscalationPolicyRepository)(nil)
	_ repository.IncidentRepository       = (*incidentRepository)(nil)
	_ repository.TeamInvitationRepository = (*TeamInvitationRepo)(nil)
	_ repository.MetricsRepository        = (*metricsRepository)(nil)
	_ repository.NotificationRepository   = (*NotificationRepository)(nil)
	_ repository.OrganizationRepository   = (*OrganizationRepository)(nil)
	_ repository.RoutingRuleRepository    = (*RoutingRuleRepository)(nil)
	_ repository.ScheduleRepository       = (*ScheduleRepository)(nil)
	_ repository.TeamRepository           = (*TeamRepository)(nil)
	_ repository.UserRepository           = (*UserRepository)(nil)
	_ repository.WebhookRepository        = (*webhookRepository)(nil)
)
