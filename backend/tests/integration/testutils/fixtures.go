package testutils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

// TestFixtures provides methods to create test data
type TestFixtures struct {
	server *TestServer
}

// NewTestFixtures creates a new test fixtures helper
func NewTestFixtures(server *TestServer) *TestFixtures {
	return &TestFixtures{server: server}
}

// TestUser holds user data with tokens for testing
type TestUser struct {
	User         *domain.User
	Organization *domain.Organization
	AccessToken  string
	RefreshToken string
}

// CreateUser creates a user with the given email, username, and organization name
func (f *TestFixtures) CreateUser(ctx context.Context, email, username, orgName string) (*TestUser, error) {
	req := &service.RegisterRequest{
		Email:            email,
		Username:         username,
		Password:         "TestPassword123!",
		FullName:         "Test User",
		OrganizationName: orgName,
	}

	resp, err := f.server.AuthService.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &TestUser{
		User:         resp.User,
		Organization: resp.Organization,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// CreateUniqueUser creates a user with unique email, username, and organization
func (f *TestFixtures) CreateUniqueUser(ctx context.Context) (*TestUser, error) {
	id := uuid.New().String()[:8]
	return f.CreateUser(ctx,
		fmt.Sprintf("user_%s@test.com", id),
		fmt.Sprintf("user_%s", id),
		fmt.Sprintf("Org_%s", id),
	)
}

// CreateTeam creates a team in the organization
func (f *TestFixtures) CreateTeam(ctx context.Context, orgID uuid.UUID, name string) (*domain.Team, error) {
	desc := "Test team description"
	req := &service.CreateTeamRequest{
		Name:        name,
		Description: &desc,
	}
	return f.server.TeamService.CreateTeam(ctx, orgID, req)
}

// CreateUniqueTeam creates a team with a unique name
func (f *TestFixtures) CreateUniqueTeam(ctx context.Context, orgID uuid.UUID) (*domain.Team, error) {
	id := uuid.New().String()[:8]
	return f.CreateTeam(ctx, orgID, fmt.Sprintf("Team_%s", id))
}

// CreateAlert creates an alert in the organization
func (f *TestFixtures) CreateAlert(ctx context.Context, orgID uuid.UUID, message string) (*domain.Alert, error) {
	req := &service.CreateAlertRequest{
		Source:   "test",
		Priority: "P3",
		Message:  message,
		Tags:     []string{"test"},
	}
	return f.server.AlertService.CreateAlert(ctx, orgID, req)
}

// CreateUniqueAlert creates an alert with a unique message
func (f *TestFixtures) CreateUniqueAlert(ctx context.Context, orgID uuid.UUID) (*domain.Alert, error) {
	id := uuid.New().String()[:8]
	return f.CreateAlert(ctx, orgID, fmt.Sprintf("Alert_%s", id))
}

// CreateSchedule creates a schedule in the organization
func (f *TestFixtures) CreateSchedule(ctx context.Context, orgID uuid.UUID, name string) (*domain.Schedule, error) {
	desc := "Test schedule"
	req := &service.CreateScheduleRequest{
		Name:        name,
		Description: &desc,
		Timezone:    "UTC",
	}
	return f.server.ScheduleService.CreateSchedule(ctx, orgID, req)
}

// CreateUniqueSchedule creates a schedule with a unique name
func (f *TestFixtures) CreateUniqueSchedule(ctx context.Context, orgID uuid.UUID) (*domain.Schedule, error) {
	id := uuid.New().String()[:8]
	return f.CreateSchedule(ctx, orgID, fmt.Sprintf("Schedule_%s", id))
}

// CreateEscalationPolicy creates an escalation policy in the organization
func (f *TestFixtures) CreateEscalationPolicy(ctx context.Context, orgID uuid.UUID, name string) (*domain.EscalationPolicy, error) {
	desc := "Test escalation policy"
	req := &service.CreateEscalationPolicyRequest{
		Name:        name,
		Description: &desc,
	}
	return f.server.EscalationService.CreatePolicy(ctx, orgID, req)
}

// CreateUniqueEscalationPolicy creates an escalation policy with a unique name
func (f *TestFixtures) CreateUniqueEscalationPolicy(ctx context.Context, orgID uuid.UUID) (*domain.EscalationPolicy, error) {
	id := uuid.New().String()[:8]
	return f.CreateEscalationPolicy(ctx, orgID, fmt.Sprintf("Policy_%s", id))
}

// CreateNotificationChannel creates a notification channel in the organization
func (f *TestFixtures) CreateNotificationChannel(ctx context.Context, orgID uuid.UUID, name string) (*domain.NotificationChannel, error) {
	configJSON, _ := json.Marshal(map[string]interface{}{
		"smtp_host":     "smtp.test.com",
		"smtp_port":     587,
		"smtp_username": "test@test.com",
		"smtp_password": "testpassword",
		"from":          "test@test.com",
		"from_address":  "test@test.com",
	})

	req := &domain.CreateNotificationChannelRequest{
		Name:        name,
		ChannelType: domain.ChannelTypeEmail,
		IsEnabled:   true,
		Config:      configJSON,
	}
	return f.server.NotificationService.CreateChannel(ctx, orgID, req)
}

// CreateUniqueNotificationChannel creates a notification channel with a unique name
func (f *TestFixtures) CreateUniqueNotificationChannel(ctx context.Context, orgID uuid.UUID) (*domain.NotificationChannel, error) {
	id := uuid.New().String()[:8]
	return f.CreateNotificationChannel(ctx, orgID, fmt.Sprintf("Channel_%s", id))
}

// CreateIncident creates an incident in the organization
func (f *TestFixtures) CreateIncident(ctx context.Context, orgID, userID uuid.UUID, title string) (*domain.Incident, error) {
	desc := "Test incident description"
	req := &service.CreateIncidentRequest{
		Title:       title,
		Description: &desc,
		Severity:    "medium",
		Priority:    "P3",
	}
	return f.server.IncidentService.CreateIncident(ctx, orgID, userID, req)
}

// CreateUniqueIncident creates an incident with a unique title
func (f *TestFixtures) CreateUniqueIncident(ctx context.Context, orgID, userID uuid.UUID) (*domain.Incident, error) {
	id := uuid.New().String()[:8]
	return f.CreateIncident(ctx, orgID, userID, fmt.Sprintf("Incident_%s", id))
}

// CreateWebhookEndpoint creates a webhook endpoint in the organization
func (f *TestFixtures) CreateWebhookEndpoint(ctx context.Context, orgID uuid.UUID, name, url string) (*domain.WebhookEndpoint, error) {
	req := &domain.CreateWebhookEndpointRequest{
		Name:         name,
		URL:          url,
		Enabled:      true,
		AlertCreated: true,
		AlertUpdated: true,
	}
	return f.server.WebhookService.CreateEndpoint(ctx, orgID, req)
}

// CreateUniqueWebhookEndpoint creates a webhook endpoint with a unique name
func (f *TestFixtures) CreateUniqueWebhookEndpoint(ctx context.Context, orgID uuid.UUID) (*domain.WebhookEndpoint, error) {
	id := uuid.New().String()[:8]
	return f.CreateWebhookEndpoint(ctx, orgID,
		fmt.Sprintf("Webhook_%s", id),
		fmt.Sprintf("https://webhook.test/%s", id),
	)
}

// CreateIncomingWebhookToken creates an incoming webhook token
func (f *TestFixtures) CreateIncomingWebhookToken(ctx context.Context, orgID uuid.UUID, name string) (*domain.IncomingWebhookToken, error) {
	req := &domain.CreateIncomingWebhookTokenRequest{
		Name:            name,
		IntegrationType: domain.IncomingWebhookGeneric,
	}
	return f.server.WebhookService.CreateIncomingToken(ctx, orgID, req)
}

// CreateUniqueIncomingWebhookToken creates an incoming webhook token with a unique name
func (f *TestFixtures) CreateUniqueIncomingWebhookToken(ctx context.Context, orgID uuid.UUID) (*domain.IncomingWebhookToken, error) {
	id := uuid.New().String()[:8]
	return f.CreateIncomingWebhookToken(ctx, orgID, fmt.Sprintf("IncomingWebhook_%s", id))
}
