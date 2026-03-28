package postgres

import (
	"github.com/nmn3m/pulsar/backend/internal/core/port/outbound"
)

// Compile-time interface checks
var (
	_ outbound.AlertRepository          = (*AlertRepository)(nil)
	_ outbound.APIKeyRepository         = (*apiKeyRepository)(nil)
	_ outbound.DNDSettingsRepository     = (*DNDSettingsRepository)(nil)
	_ outbound.EmailVerificationRepository = (*EmailVerificationRepository)(nil)
	_ outbound.EscalationPolicyRepository = (*EscalationPolicyRepository)(nil)
	_ outbound.IncidentRepository       = (*incidentRepository)(nil)
	_ outbound.TeamInvitationRepository = (*TeamInvitationRepo)(nil)
	_ outbound.MetricsRepository        = (*metricsRepository)(nil)
	_ outbound.NotificationRepository   = (*NotificationRepository)(nil)
	_ outbound.OrganizationRepository   = (*OrganizationRepository)(nil)
	_ outbound.RoutingRuleRepository    = (*RoutingRuleRepository)(nil)
	_ outbound.ScheduleRepository       = (*ScheduleRepository)(nil)
	_ outbound.TeamRepository           = (*TeamRepository)(nil)
	_ outbound.UserRepository           = (*UserRepository)(nil)
	_ outbound.WebhookRepository        = (*webhookRepository)(nil)
)
