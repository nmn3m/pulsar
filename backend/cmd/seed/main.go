package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/nmn3m/pulsar/backend/internal/config"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
	"github.com/nmn3m/pulsar/backend/internal/repository/postgres"
)

// Demo Data Configuration
const (
	DemoPassword = "DemoPass123!" // Password for all demo users
)

// Seed data containers
var (
	demoOrg   *domain.Organization
	demoUsers map[string]*domain.User
	demoTeams map[string]*domain.Team
)

func main() {
	fmt.Println("🚀 Pulsar Demo Data Seeder")
	fmt.Println("==========================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := postgres.NewDB(cfg.Database.URL)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("✅ Connected to database")

	ctx := context.Background()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)
	escalationRepo := postgres.NewEscalationPolicyRepository(db)
	alertRepo := postgres.NewAlertRepository(db)
	incidentRepo := postgres.NewIncidentRepository(db.DB)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db.DB)
	webhookRepo := postgres.NewWebhookRepository(db.DB)
	apiKeyRepo := postgres.NewAPIKeyRepository(db.DB)
	routingRepo := postgres.NewRoutingRuleRepository(db)
	dndRepo := postgres.NewDNDSettingsRepository(db)

	// Initialize data containers
	demoUsers = make(map[string]*domain.User)
	demoTeams = make(map[string]*domain.Team)

	// Run seed functions in order
	if err := seedOrganization(ctx, orgRepo); err != nil {
		fmt.Printf("❌ Failed to seed organization: %v\n", err)
		os.Exit(1)
	}

	if err := seedUsers(ctx, userRepo, orgRepo, emailVerificationRepo); err != nil {
		fmt.Printf("❌ Failed to seed users: %v\n", err)
		os.Exit(1)
	}

	if err := seedTeams(ctx, teamRepo); err != nil {
		fmt.Printf("❌ Failed to seed teams: %v\n", err)
		os.Exit(1)
	}

	schedules, err := seedSchedules(ctx, scheduleRepo)
	if err != nil {
		fmt.Printf("❌ Failed to seed schedules: %v\n", err)
		os.Exit(1)
	}

	policies, err := seedEscalationPolicies(ctx, escalationRepo, schedules)
	if err != nil {
		fmt.Printf("❌ Failed to seed escalation policies: %v\n", err)
		os.Exit(1)
	}

	alerts, err := seedAlerts(ctx, alertRepo, policies)
	if err != nil {
		fmt.Printf("❌ Failed to seed alerts: %v\n", err)
		os.Exit(1)
	}

	if err := seedIncidents(ctx, incidentRepo, alerts); err != nil {
		fmt.Printf("❌ Failed to seed incidents: %v\n", err)
		os.Exit(1)
	}

	if err := seedNotificationChannels(ctx, notificationRepo); err != nil {
		fmt.Printf("❌ Failed to seed notification channels: %v\n", err)
		os.Exit(1)
	}

	if err := seedWebhooks(ctx, webhookRepo); err != nil {
		fmt.Printf("❌ Failed to seed webhooks: %v\n", err)
		os.Exit(1)
	}

	apiKey, err := seedAPIKeys(ctx, apiKeyRepo)
	if err != nil {
		fmt.Printf("❌ Failed to seed API keys: %v\n", err)
		os.Exit(1)
	}

	if err := seedRoutingRules(ctx, routingRepo, policies); err != nil {
		fmt.Printf("❌ Failed to seed routing rules: %v\n", err)
		os.Exit(1)
	}

	if err := seedDNDSettings(ctx, dndRepo); err != nil {
		fmt.Printf("❌ Failed to seed DND settings: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n==========================")
	fmt.Println("✅ Demo data seeded successfully!")
	fmt.Println("\n📋 Login Credentials:")
	fmt.Println("==========================")
	for name, user := range demoUsers {
		fmt.Printf("  %s: %s / %s\n", name, user.Email, DemoPassword)
	}
	if apiKey != "" {
		fmt.Println("\n🔑 Demo API Key:")
		fmt.Println("==========================")
		fmt.Printf("  %s\n", apiKey)
		fmt.Println("  (Save this - it won't be shown again!)")
	}
	fmt.Println("\n🎯 Ready for demo presentation!")
}

func seedOrganization(ctx context.Context, repo *postgres.OrganizationRepository) error {
	fmt.Println("\n📁 Creating organization...")

	org := &domain.Organization{
		ID:   uuid.New(),
		Name: "ACME Corporation",
		Slug: "acme-corp",
	}

	if err := repo.Create(ctx, org); err != nil {
		return fmt.Errorf("create organization: %w", err)
	}

	demoOrg = org
	fmt.Printf("   ✓ Organization: %s (ID: %s)\n", org.Name, org.ID)
	return nil
}

func seedUsers(ctx context.Context, userRepo *postgres.UserRepository, orgRepo *postgres.OrganizationRepository, emailVerificationRepo *postgres.EmailVerificationRepository) error {
	fmt.Println("\n👥 Creating users...")

	// Hash password once for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	users := []struct {
		Key      string
		Email    string
		Username string
		FullName string
		Role     domain.UserRole
	}{
		{"admin", "admin@acme-corp.com", "admin", "Demo Administrator", domain.RoleOwner},
		{"alice", "alice@acme-corp.com", "alice", "Alice Chen", domain.RoleAdmin},
		{"bob", "bob@acme-corp.com", "bob", "Bob Martinez", domain.RoleMember},
		{"carol", "carol@acme-corp.com", "carol", "Carol Williams", domain.RoleMember},
		{"david", "david@acme-corp.com", "david", "David Kim", domain.RoleMember},
		{"emma", "emma@acme-corp.com", "emma", "Emma Johnson", domain.RoleAdmin},
	}

	for _, u := range users {
		fullName := u.FullName
		user := &domain.User{
			ID:            uuid.New(),
			Email:         u.Email,
			Username:      u.Username,
			PasswordHash:  string(hashedPassword),
			FullName:      &fullName,
			Timezone:      "America/New_York",
			IsActive:      true,
			EmailVerified: true, // Pre-verified for demo
		}

		if err := userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("create user %s: %w", u.Username, err)
		}

		// Create email verification record (marked as verified)
		verification := &domain.EmailVerification{
			ID:        uuid.New(),
			UserID:    user.ID,
			Email:     user.Email,
			OTP:       "000000", // Dummy OTP since already verified
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Verified:  true,
		}
		if err := emailVerificationRepo.Create(ctx, verification); err != nil {
			return fmt.Errorf("create email verification for %s: %w", u.Username, err)
		}

		// Add user to organization
		if err := orgRepo.AddUser(ctx, demoOrg.ID, user.ID, u.Role); err != nil {
			return fmt.Errorf("add user %s to org: %w", u.Username, err)
		}

		demoUsers[u.Key] = user
		fmt.Printf("   ✓ User: %s (%s) - Role: %s\n", u.FullName, u.Email, u.Role)
	}

	return nil
}

func seedTeams(ctx context.Context, repo *postgres.TeamRepository) error {
	fmt.Println("\n👨‍👩‍👧‍👦 Creating teams...")

	teams := []struct {
		Key         string
		Name        string
		Description string
		Members     []struct {
			UserKey string
			Role    string
		}
	}{
		{
			Key:         "platform",
			Name:        "Platform Engineering",
			Description: "Infrastructure, deployment pipelines, and platform reliability",
			Members: []struct {
				UserKey string
				Role    string
			}{
				{"alice", "lead"},
				{"bob", "member"},
				{"carol", "member"},
			},
		},
		{
			Key:         "backend",
			Name:        "Backend Services",
			Description: "API development, microservices, and database management",
			Members: []struct {
				UserKey string
				Role    string
			}{
				{"emma", "lead"},
				{"david", "member"},
				{"bob", "member"},
			},
		},
		{
			Key:         "frontend",
			Name:        "Frontend Team",
			Description: "Web and mobile user interfaces",
			Members: []struct {
				UserKey string
				Role    string
			}{
				{"carol", "lead"},
				{"alice", "member"},
			},
		},
		{
			Key:         "oncall",
			Name:        "On-Call Rotation",
			Description: "Cross-functional on-call team for incident response",
			Members: []struct {
				UserKey string
				Role    string
			}{
				{"admin", "lead"},
				{"alice", "member"},
				{"emma", "member"},
				{"bob", "member"},
			},
		},
	}

	for _, t := range teams {
		desc := t.Description
		team := &domain.Team{
			ID:             uuid.New(),
			OrganizationID: demoOrg.ID,
			Name:           t.Name,
			Description:    &desc,
		}

		if err := repo.Create(ctx, team); err != nil {
			return fmt.Errorf("create team %s: %w", t.Name, err)
		}

		// Add members
		for _, m := range t.Members {
			user := demoUsers[m.UserKey]
			if user == nil {
				continue
			}
			role := domain.TeamRoleMember
			if m.Role == "lead" {
				role = domain.TeamRoleLead
			}
			if err := repo.AddMember(ctx, team.ID, user.ID, role); err != nil {
				return fmt.Errorf("add member to team %s: %w", t.Name, err)
			}
		}

		demoTeams[t.Key] = team
		fmt.Printf("   ✓ Team: %s (%d members)\n", t.Name, len(t.Members))
	}

	return nil
}

func seedSchedules(ctx context.Context, repo *postgres.ScheduleRepository) (map[string]*domain.Schedule, error) {
	fmt.Println("\n📅 Creating schedules...")

	schedules := make(map[string]*domain.Schedule)

	// Schedule 1: Platform On-Call (Weekly)
	platformDesc := "24/7 on-call rotation for platform engineering"
	platformSchedule := &domain.Schedule{
		ID:             uuid.New(),
		OrganizationID: demoOrg.ID,
		TeamID:         &demoTeams["platform"].ID,
		Name:           "Platform On-Call",
		Description:    &platformDesc,
		Timezone:       "America/New_York",
	}

	if err := repo.Create(ctx, platformSchedule); err != nil {
		return nil, fmt.Errorf("create platform schedule: %w", err)
	}
	schedules["platform"] = platformSchedule

	// Add weekly rotation
	now := time.Now()
	monday := now.AddDate(0, 0, -int(now.Weekday())+1) // Get this Monday
	startDate := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
	handoffTime := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC) // 9 AM handoff

	handoffDay := 1 // Monday
	platformRotation := &domain.ScheduleRotation{
		ID:             uuid.New(),
		ScheduleID:     platformSchedule.ID,
		Name:           "Weekly Rotation",
		RotationType:   domain.RotationTypeWeekly,
		RotationLength: 1,
		StartDate:      startDate,
		StartTime:      handoffTime,
		HandoffDay:     &handoffDay,
		HandoffTime:    handoffTime,
	}

	if err := repo.CreateRotation(ctx, platformRotation); err != nil {
		return nil, fmt.Errorf("create platform rotation: %w", err)
	}

	// Add participants to rotation
	participants := []string{"alice", "bob", "carol"}
	for i, userKey := range participants {
		user := demoUsers[userKey]
		participant := &domain.ScheduleRotationParticipant{
			ID:         uuid.New(),
			RotationID: platformRotation.ID,
			UserID:     user.ID,
			Position:   i + 1,
		}
		if err := repo.AddParticipant(ctx, participant); err != nil {
			return nil, fmt.Errorf("add participant: %w", err)
		}
	}

	fmt.Printf("   ✓ Schedule: %s (Weekly rotation, %d participants)\n", platformSchedule.Name, len(participants))

	// Schedule 2: Backend On-Call (Daily)
	backendDesc := "Daily on-call rotation for backend services"
	backendSchedule := &domain.Schedule{
		ID:             uuid.New(),
		OrganizationID: demoOrg.ID,
		TeamID:         &demoTeams["backend"].ID,
		Name:           "Backend On-Call",
		Description:    &backendDesc,
		Timezone:       "America/New_York",
	}

	if err := repo.Create(ctx, backendSchedule); err != nil {
		return nil, fmt.Errorf("create backend schedule: %w", err)
	}
	schedules["backend"] = backendSchedule

	// Add daily rotation
	backendRotation := &domain.ScheduleRotation{
		ID:             uuid.New(),
		ScheduleID:     backendSchedule.ID,
		Name:           "Daily Rotation",
		RotationType:   domain.RotationTypeDaily,
		RotationLength: 1,
		StartDate:      startDate,
		StartTime:      handoffTime,
		HandoffTime:    handoffTime,
	}

	if err := repo.CreateRotation(ctx, backendRotation); err != nil {
		return nil, fmt.Errorf("create backend rotation: %w", err)
	}

	backendParticipants := []string{"emma", "david"}
	for i, userKey := range backendParticipants {
		user := demoUsers[userKey]
		participant := &domain.ScheduleRotationParticipant{
			ID:         uuid.New(),
			RotationID: backendRotation.ID,
			UserID:     user.ID,
			Position:   i + 1,
		}
		if err := repo.AddParticipant(ctx, participant); err != nil {
			return nil, fmt.Errorf("add participant: %w", err)
		}
	}

	fmt.Printf("   ✓ Schedule: %s (Daily rotation, %d participants)\n", backendSchedule.Name, len(backendParticipants))

	return schedules, nil
}

func seedEscalationPolicies(ctx context.Context, repo *postgres.EscalationPolicyRepository, schedules map[string]*domain.Schedule) (map[string]*domain.EscalationPolicy, error) {
	fmt.Println("\n📈 Creating escalation policies...")

	policies := make(map[string]*domain.EscalationPolicy)

	// Policy 1: Platform Escalation
	platformDesc := "Escalation policy for platform incidents"
	repeatCount := 2
	platformPolicy := &domain.EscalationPolicy{
		ID:             uuid.New(),
		OrganizationID: demoOrg.ID,
		Name:           "Platform Escalation",
		Description:    &platformDesc,
		RepeatEnabled:  true,
		RepeatCount:    &repeatCount,
	}

	if err := repo.Create(ctx, platformPolicy); err != nil {
		return nil, fmt.Errorf("create platform policy: %w", err)
	}
	policies["platform"] = platformPolicy

	// Rule 1: Notify on-call (5 min delay)
	rule1 := &domain.EscalationRule{
		ID:              uuid.New(),
		PolicyID:        platformPolicy.ID,
		Position:        1,
		EscalationDelay: 5, // 5 minutes
	}
	if err := repo.CreateRule(ctx, rule1); err != nil {
		return nil, fmt.Errorf("create rule 1: %w", err)
	}

	// Target: Platform On-Call Schedule
	target1 := &domain.EscalationTarget{
		ID:         uuid.New(),
		RuleID:     rule1.ID,
		TargetType: domain.EscalationTargetTypeSchedule,
		TargetID:   schedules["platform"].ID,
	}
	if err := repo.AddTarget(ctx, target1); err != nil {
		return nil, fmt.Errorf("add target 1: %w", err)
	}

	// Rule 2: Notify entire team (10 min delay)
	rule2 := &domain.EscalationRule{
		ID:              uuid.New(),
		PolicyID:        platformPolicy.ID,
		Position:        2,
		EscalationDelay: 10, // 10 minutes
	}
	if err := repo.CreateRule(ctx, rule2); err != nil {
		return nil, fmt.Errorf("create rule 2: %w", err)
	}

	// Target: Platform Team
	target2 := &domain.EscalationTarget{
		ID:         uuid.New(),
		RuleID:     rule2.ID,
		TargetType: domain.EscalationTargetTypeTeam,
		TargetID:   demoTeams["platform"].ID,
	}
	if err := repo.AddTarget(ctx, target2); err != nil {
		return nil, fmt.Errorf("add target 2: %w", err)
	}

	// Rule 3: Notify admin (15 min delay)
	rule3 := &domain.EscalationRule{
		ID:              uuid.New(),
		PolicyID:        platformPolicy.ID,
		Position:        3,
		EscalationDelay: 15, // 15 minutes
	}
	if err := repo.CreateRule(ctx, rule3); err != nil {
		return nil, fmt.Errorf("create rule 3: %w", err)
	}

	// Target: Admin user
	target3 := &domain.EscalationTarget{
		ID:         uuid.New(),
		RuleID:     rule3.ID,
		TargetType: domain.EscalationTargetTypeUser,
		TargetID:   demoUsers["admin"].ID,
	}
	if err := repo.AddTarget(ctx, target3); err != nil {
		return nil, fmt.Errorf("add target 3: %w", err)
	}

	fmt.Printf("   ✓ Policy: %s (3 rules, repeat %d times)\n", platformPolicy.Name, repeatCount)

	// Policy 2: Backend Escalation (simpler)
	backendDesc := "Escalation policy for backend service incidents"
	backendPolicy := &domain.EscalationPolicy{
		ID:             uuid.New(),
		OrganizationID: demoOrg.ID,
		Name:           "Backend Escalation",
		Description:    &backendDesc,
		RepeatEnabled:  false,
	}

	if err := repo.Create(ctx, backendPolicy); err != nil {
		return nil, fmt.Errorf("create backend policy: %w", err)
	}
	policies["backend"] = backendPolicy

	// Rule 1: Notify backend on-call
	backendRule1 := &domain.EscalationRule{
		ID:              uuid.New(),
		PolicyID:        backendPolicy.ID,
		Position:        1,
		EscalationDelay: 5,
	}
	if err := repo.CreateRule(ctx, backendRule1); err != nil {
		return nil, fmt.Errorf("create backend rule 1: %w", err)
	}

	backendTarget1 := &domain.EscalationTarget{
		ID:         uuid.New(),
		RuleID:     backendRule1.ID,
		TargetType: domain.EscalationTargetTypeSchedule,
		TargetID:   schedules["backend"].ID,
	}
	if err := repo.AddTarget(ctx, backendTarget1); err != nil {
		return nil, fmt.Errorf("add backend target 1: %w", err)
	}

	// Rule 2: Notify backend team lead
	backendRule2 := &domain.EscalationRule{
		ID:              uuid.New(),
		PolicyID:        backendPolicy.ID,
		Position:        2,
		EscalationDelay: 10,
	}
	if err := repo.CreateRule(ctx, backendRule2); err != nil {
		return nil, fmt.Errorf("create backend rule 2: %w", err)
	}

	backendTarget2 := &domain.EscalationTarget{
		ID:         uuid.New(),
		RuleID:     backendRule2.ID,
		TargetType: domain.EscalationTargetTypeUser,
		TargetID:   demoUsers["emma"].ID,
	}
	if err := repo.AddTarget(ctx, backendTarget2); err != nil {
		return nil, fmt.Errorf("add backend target 2: %w", err)
	}

	fmt.Printf("   ✓ Policy: %s (2 rules, no repeat)\n", backendPolicy.Name)

	return policies, nil
}

func seedAlerts(ctx context.Context, repo *postgres.AlertRepository, policies map[string]*domain.EscalationPolicy) ([]*domain.Alert, error) {
	fmt.Println("\n🚨 Creating sample alerts...")

	var alerts []*domain.Alert
	now := time.Now()

	alertsData := []struct {
		Message     string
		Description string
		Priority    domain.AlertPriority
		Status      domain.AlertStatus
		Source      string
		Tags        []string
		PolicyKey   string
		AgeHours    int
	}{
		{
			Message:     "High CPU usage on prod-api-server-01",
			Description: "CPU usage exceeded 95% for more than 5 minutes. Process: java. PID: 12345.",
			Priority:    domain.PriorityP1,
			Status:      domain.AlertStatusOpen,
			Source:      "Prometheus",
			Tags:        []string{"production", "api", "performance", "critical"},
			PolicyKey:   "platform",
			AgeHours:    0,
		},
		{
			Message:     "Database connection pool exhausted",
			Description: "PostgreSQL connection pool at 100% capacity. New connections are being rejected.",
			Priority:    domain.PriorityP1,
			Status:      domain.AlertStatusAcknowledged,
			Source:      "Datadog",
			Tags:        []string{"production", "database", "postgres"},
			PolicyKey:   "backend",
			AgeHours:    1,
		},
		{
			Message:     "Memory leak detected in payment-service",
			Description: "Memory usage growing steadily. Current: 85%. Rate: +5% per hour.",
			Priority:    domain.PriorityP2,
			Status:      domain.AlertStatusOpen,
			Source:      "NewRelic",
			Tags:        []string{"production", "payment", "memory"},
			PolicyKey:   "backend",
			AgeHours:    2,
		},
		{
			Message:     "SSL certificate expiring in 7 days",
			Description: "Certificate for api.acme-corp.com expires on 2025-01-20. Renewal required.",
			Priority:    domain.PriorityP3,
			Status:      domain.AlertStatusOpen,
			Source:      "CertManager",
			Tags:        []string{"security", "ssl", "maintenance"},
			PolicyKey:   "platform",
			AgeHours:    24,
		},
		{
			Message:     "High error rate on /api/v1/checkout endpoint",
			Description: "Error rate: 15% (threshold: 5%). HTTP 500 responses increasing.",
			Priority:    domain.PriorityP2,
			Status:      domain.AlertStatusClosed,
			Source:      "Prometheus",
			Tags:        []string{"production", "api", "checkout", "errors"},
			PolicyKey:   "backend",
			AgeHours:    48,
		},
		{
			Message:     "Disk space warning on log-aggregator-01",
			Description: "Disk usage at 80%. Estimated time to full: 3 days.",
			Priority:    domain.PriorityP4,
			Status:      domain.AlertStatusOpen,
			Source:      "CloudWatch",
			Tags:        []string{"infrastructure", "disk", "logs"},
			PolicyKey:   "platform",
			AgeHours:    12,
		},
		{
			Message:     "Kubernetes pod crashlooping: notification-worker",
			Description: "Pod has restarted 5 times in the last hour. Last exit code: 137 (OOMKilled).",
			Priority:    domain.PriorityP2,
			Status:      domain.AlertStatusOpen,
			Source:      "Kubernetes",
			Tags:        []string{"production", "k8s", "notification", "oom"},
			PolicyKey:   "platform",
			AgeHours:    3,
		},
		{
			Message:     "Unusual login activity detected",
			Description: "Multiple failed login attempts from IP 192.168.1.100. User: service-account.",
			Priority:    domain.PriorityP3,
			Status:      domain.AlertStatusAcknowledged,
			Source:      "SecurityHub",
			Tags:        []string{"security", "auth", "suspicious"},
			PolicyKey:   "platform",
			AgeHours:    6,
		},
	}

	for _, a := range alertsData {
		desc := a.Description
		policyID := policies[a.PolicyKey].ID
		createdAt := now.Add(-time.Duration(a.AgeHours) * time.Hour)

		alert := &domain.Alert{
			ID:                 uuid.New(),
			OrganizationID:     demoOrg.ID,
			Source:             a.Source,
			Priority:           a.Priority,
			Status:             a.Status,
			Message:            a.Message,
			Description:        &desc,
			Tags:               a.Tags,
			EscalationPolicyID: &policyID,
			CreatedAt:          createdAt,
			UpdatedAt:          createdAt,
		}

		// Set acknowledgment for acknowledged alerts
		if a.Status == domain.AlertStatusAcknowledged {
			ackBy := demoUsers["alice"].ID
			ackAt := createdAt.Add(5 * time.Minute)
			alert.AcknowledgedBy = &ackBy
			alert.AcknowledgedAt = &ackAt
		}

		// Set closure for closed alerts
		if a.Status == domain.AlertStatusClosed {
			closedBy := demoUsers["emma"].ID
			closedAt := createdAt.Add(2 * time.Hour)
			reason := "Issue resolved. Root cause: memory leak in payment service fixed in v2.3.2"
			alert.ClosedBy = &closedBy
			alert.ClosedAt = &closedAt
			alert.CloseReason = &reason
		}

		if err := repo.Create(ctx, alert); err != nil {
			return nil, fmt.Errorf("create alert: %w", err)
		}
		alerts = append(alerts, alert)
		fmt.Printf("   ✓ Alert: [%s] %s (%s)\n", a.Priority, truncate(a.Message, 40), a.Status)
	}

	return alerts, nil
}

func seedIncidents(ctx context.Context, repo repository.IncidentRepository, alerts []*domain.Alert) error {
	fmt.Println("\n🔥 Creating sample incidents...")

	now := time.Now()

	incidentsData := []struct {
		Title       string
		Description string
		Severity    domain.IncidentSeverity
		Status      domain.IncidentStatus
		Priority    domain.AlertPriority
		AgeHours    int
		AlertIndex  int // Index of alert to link, -1 for none
	}{
		{
			Title:       "Production API Performance Degradation",
			Description: "Multiple API servers experiencing high CPU usage causing elevated response times and timeouts for customers.",
			Severity:    domain.IncidentSeverityCritical,
			Status:      domain.IncidentStatusInvestigating,
			Priority:    domain.PriorityP1,
			AgeHours:    0,
			AlertIndex:  0, // Link to CPU alert
		},
		{
			Title:       "Payment Service Memory Issues",
			Description: "Payment service experiencing memory growth leading to degraded performance and intermittent failures.",
			Severity:    domain.IncidentSeverityHigh,
			Status:      domain.IncidentStatusIdentified,
			Priority:    domain.PriorityP2,
			AgeHours:    2,
			AlertIndex:  2, // Link to memory leak alert
		},
		{
			Title:       "Database Connection Issues",
			Description: "PostgreSQL connection pool exhaustion causing service disruptions across multiple microservices.",
			Severity:    domain.IncidentSeverityCritical,
			Status:      domain.IncidentStatusResolved,
			Priority:    domain.PriorityP1,
			AgeHours:    24,
			AlertIndex:  1, // Link to DB connection alert
		},
	}

	for _, inc := range incidentsData {
		desc := inc.Description
		createdAt := now.Add(-time.Duration(inc.AgeHours) * time.Hour)
		teamID := demoTeams["platform"].ID

		incident := &domain.Incident{
			ID:               uuid.New(),
			OrganizationID:   demoOrg.ID,
			Title:            inc.Title,
			Description:      &desc,
			Severity:         inc.Severity,
			Status:           inc.Status,
			Priority:         inc.Priority,
			CreatedByUserID:  demoUsers["admin"].ID,
			AssignedToTeamID: &teamID,
			StartedAt:        createdAt,
			CreatedAt:        createdAt,
			UpdatedAt:        createdAt,
		}

		// Set resolved time for resolved incidents
		if inc.Status == domain.IncidentStatusResolved {
			resolvedAt := createdAt.Add(4 * time.Hour)
			incident.ResolvedAt = &resolvedAt
		}

		if err := repo.Create(ctx, incident); err != nil {
			return fmt.Errorf("create incident: %w", err)
		}

		// Add responders
		responders := []struct {
			UserKey string
			Role    domain.ResponderRole
		}{
			{"alice", domain.ResponderRoleIncidentCommander},
			{"bob", domain.ResponderRoleResponder},
		}

		for _, r := range responders {
			responder := &domain.IncidentResponder{
				ID:         uuid.New(),
				IncidentID: incident.ID,
				UserID:     demoUsers[r.UserKey].ID,
				Role:       r.Role,
			}
			if err := repo.AddResponder(ctx, responder); err != nil {
				return fmt.Errorf("add responder: %w", err)
			}
		}

		// Link alert if specified
		if inc.AlertIndex >= 0 && inc.AlertIndex < len(alerts) {
			incidentAlert := &domain.IncidentAlert{
				ID:             uuid.New(),
				IncidentID:     incident.ID,
				AlertID:        alerts[inc.AlertIndex].ID,
				LinkedByUserID: &demoUsers["admin"].ID,
			}
			if err := repo.LinkAlert(ctx, incidentAlert); err != nil {
				return fmt.Errorf("link alert: %w", err)
			}
		}

		// Add timeline events
		events := []struct {
			EventType   domain.TimelineEventType
			Description string
			MinutesAgo  int
		}{
			{domain.TimelineEventCreated, "Incident created", inc.AgeHours * 60},
			{domain.TimelineEventNoteAdded, "Initial investigation started. Gathering metrics and logs.", inc.AgeHours*60 - 5},
			{domain.TimelineEventResponderAdded, "Alice Chen joined as Incident Commander", inc.AgeHours*60 - 10},
		}

		if inc.Status == domain.IncidentStatusIdentified || inc.Status == domain.IncidentStatusResolved {
			events = append(events, struct {
				EventType   domain.TimelineEventType
				Description string
				MinutesAgo  int
			}{domain.TimelineEventStatusChanged, "Root cause identified: " + truncate(inc.Description, 50), inc.AgeHours*60 - 30})
		}

		if inc.Status == domain.IncidentStatusResolved {
			events = append(events, struct {
				EventType   domain.TimelineEventType
				Description string
				MinutesAgo  int
			}{domain.TimelineEventResolved, "Incident resolved. Services restored to normal operation.", inc.AgeHours*60 - 240})
		}

		for _, e := range events {
			userID := demoUsers["alice"].ID
			event := &domain.IncidentTimelineEvent{
				ID:          uuid.New(),
				IncidentID:  incident.ID,
				EventType:   e.EventType,
				UserID:      &userID,
				Description: e.Description,
				CreatedAt:   now.Add(-time.Duration(e.MinutesAgo) * time.Minute),
			}
			if err := repo.AddTimelineEvent(ctx, event); err != nil {
				return fmt.Errorf("add timeline event: %w", err)
			}
		}

		fmt.Printf("   ✓ Incident: [%s] %s (%s)\n", inc.Severity, truncate(inc.Title, 35), inc.Status)
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func seedNotificationChannels(ctx context.Context, repo *postgres.NotificationRepository) error {
	fmt.Println("\n📢 Creating notification channels...")

	channels := []struct {
		Name        string
		ChannelType domain.ChannelType
		Config      map[string]interface{}
	}{
		{
			Name:        "Team Email",
			ChannelType: domain.ChannelTypeEmail,
			Config: map[string]interface{}{
				"smtp_host": "mailpit",
				"smtp_port": 1025,
				"from":      "alerts@acme-corp.com",
				"from_name": "ACME Alerts",
			},
		},
		{
			Name:        "Slack #incidents",
			ChannelType: domain.ChannelTypeSlack,
			Config: map[string]interface{}{
				"webhook_url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
				"channel":     "#incidents",
				"username":    "Pulsar Bot",
			},
		},
		{
			Name:        "Slack #alerts",
			ChannelType: domain.ChannelTypeSlack,
			Config: map[string]interface{}{
				"webhook_url": "https://hooks.slack.com/services/T00000000/B00000000/YYYYYYYYYYYYYYYYYYYYYYYY",
				"channel":     "#alerts",
				"username":    "Pulsar Bot",
			},
		},
		{
			Name:        "MS Teams Operations",
			ChannelType: domain.ChannelTypeTeams,
			Config: map[string]interface{}{
				"webhook_url": "https://outlook.office.com/webhook/xxx/IncomingWebhook/yyy/zzz",
			},
		},
		{
			Name:        "PagerDuty Integration",
			ChannelType: domain.ChannelTypeWebhook,
			Config: map[string]interface{}{
				"url":    "https://events.pagerduty.com/v2/enqueue",
				"method": "POST",
				"headers": map[string]string{
					"Content-Type": "application/json",
				},
			},
		},
	}

	for _, ch := range channels {
		configJSON, _ := json.Marshal(ch.Config)
		channel := &domain.NotificationChannel{
			OrganizationID: demoOrg.ID,
			Name:           ch.Name,
			ChannelType:    ch.ChannelType,
			IsEnabled:      true,
			Config:         configJSON,
		}

		if err := repo.CreateChannel(ctx, channel); err != nil {
			return fmt.Errorf("create channel %s: %w", ch.Name, err)
		}
		fmt.Printf("   ✓ Channel: %s (%s)\n", ch.Name, ch.ChannelType)
	}

	return nil
}

func seedWebhooks(ctx context.Context, repo repository.WebhookRepository) error {
	fmt.Println("\n🔗 Creating webhooks...")

	// Create outgoing webhook endpoints
	endpoints := []struct {
		Name    string
		URL     string
		Events  []string
		Headers map[string]string
	}{
		{
			Name:   "DataDog Events",
			URL:    "https://api.datadoghq.com/api/v1/events",
			Events: []string{"alert.created", "alert.closed", "incident.created", "incident.resolved"},
			Headers: map[string]string{
				"DD-API-KEY": "your-datadog-api-key",
			},
		},
		{
			Name:   "Jira Automation",
			URL:    "https://acme-corp.atlassian.net/rest/api/3/issue",
			Events: []string{"incident.created"},
			Headers: map[string]string{
				"Authorization": "Basic base64-encoded-credentials",
			},
		},
		{
			Name:   "Status Page Updates",
			URL:    "https://api.statuspage.io/v1/pages/xxx/incidents",
			Events: []string{"incident.created", "incident.updated", "incident.resolved"},
			Headers: map[string]string{
				"Authorization": "OAuth your-statuspage-token",
			},
		},
	}

	for _, ep := range endpoints {
		// Generate a random secret
		secretBytes := make([]byte, 32)
		rand.Read(secretBytes)
		secret := hex.EncodeToString(secretBytes)

		endpoint := &domain.WebhookEndpoint{
			ID:                uuid.New(),
			OrganizationID:    demoOrg.ID,
			Name:              ep.Name,
			URL:               ep.URL,
			Secret:            secret,
			Enabled:           true,
			AlertCreated:      contains(ep.Events, "alert.created"),
			AlertUpdated:      contains(ep.Events, "alert.updated"),
			AlertAcknowledged: contains(ep.Events, "alert.acknowledged"),
			AlertClosed:       contains(ep.Events, "alert.closed"),
			AlertEscalated:    contains(ep.Events, "alert.escalated"),
			IncidentCreated:   contains(ep.Events, "incident.created"),
			IncidentUpdated:   contains(ep.Events, "incident.updated"),
			IncidentResolved:  contains(ep.Events, "incident.resolved"),
			Headers:           ep.Headers,
			TimeoutSeconds:    30,
			MaxRetries:        3,
			RetryDelaySeconds: 60,
		}

		if err := repo.CreateEndpoint(ctx, endpoint); err != nil {
			return fmt.Errorf("create webhook endpoint %s: %w", ep.Name, err)
		}
		fmt.Printf("   ✓ Webhook Endpoint: %s\n", ep.Name)
	}

	// Create incoming webhook tokens
	tokens := []struct {
		Name            string
		IntegrationType domain.IncomingWebhookIntegrationType
		DefaultPriority string
		DefaultTags     []string
	}{
		{
			Name:            "Prometheus Alerts",
			IntegrationType: domain.IncomingWebhookPrometheus,
			DefaultPriority: "P2",
			DefaultTags:     []string{"prometheus", "monitoring"},
		},
		{
			Name:            "Grafana Alerts",
			IntegrationType: domain.IncomingWebhookGrafana,
			DefaultPriority: "P3",
			DefaultTags:     []string{"grafana", "monitoring"},
		},
		{
			Name:            "Datadog Monitors",
			IntegrationType: domain.IncomingWebhookDatadog,
			DefaultPriority: "P2",
			DefaultTags:     []string{"datadog", "monitoring"},
		},
		{
			Name:            "Generic Webhook",
			IntegrationType: domain.IncomingWebhookGeneric,
			DefaultPriority: "P3",
			DefaultTags:     []string{"external"},
		},
	}

	fmt.Println("\n   📥 Incoming Webhook Tokens:")
	for _, t := range tokens {
		// Generate a random token
		tokenBytes := make([]byte, 24)
		rand.Read(tokenBytes)
		tokenStr := hex.EncodeToString(tokenBytes)

		token := &domain.IncomingWebhookToken{
			ID:              uuid.New(),
			OrganizationID:  demoOrg.ID,
			Name:            t.Name,
			Token:           tokenStr,
			Enabled:         true,
			IntegrationType: t.IntegrationType,
			DefaultPriority: t.DefaultPriority,
			DefaultTags:     t.DefaultTags,
			RequestCount:    0,
		}

		if err := repo.CreateIncomingToken(ctx, token); err != nil {
			return fmt.Errorf("create incoming token %s: %w", t.Name, err)
		}
		fmt.Printf("      ✓ %s: /api/v1/webhook/%s\n", t.Name, tokenStr[:16]+"...")
	}

	return nil
}

func seedAPIKeys(ctx context.Context, repo repository.APIKeyRepository) (string, error) {
	fmt.Println("\n🔑 Creating API keys...")

	// Generate API key
	rawKey, keyPrefix, keyHash, err := domain.GenerateAPIKey()
	if err != nil {
		return "", fmt.Errorf("generate API key: %w", err)
	}

	// Set expiration to 1 year from now
	expiresAt := time.Now().AddDate(1, 0, 0)

	apiKey := &domain.APIKey{
		ID:             uuid.New(),
		OrganizationID: demoOrg.ID,
		UserID:         demoUsers["admin"].ID,
		Name:           "Demo Integration Key",
		KeyPrefix:      keyPrefix,
		KeyHash:        keyHash,
		Scopes:         []string{"*"}, // Full access for demo
		IsActive:       true,
		ExpiresAt:      &expiresAt,
	}

	if err := repo.Create(ctx, apiKey); err != nil {
		return "", fmt.Errorf("create API key: %w", err)
	}
	fmt.Printf("   ✓ API Key: %s (full access, expires in 1 year)\n", apiKey.Name)

	// Create a read-only key
	rawKey2, keyPrefix2, keyHash2, _ := domain.GenerateAPIKey()
	apiKey2 := &domain.APIKey{
		ID:             uuid.New(),
		OrganizationID: demoOrg.ID,
		UserID:         demoUsers["alice"].ID,
		Name:           "Read-Only Dashboard Key",
		KeyPrefix:      keyPrefix2,
		KeyHash:        keyHash2,
		Scopes:         []string{"alerts:read", "incidents:read", "teams:read"},
		IsActive:       true,
		ExpiresAt:      &expiresAt,
	}

	if err := repo.Create(ctx, apiKey2); err != nil {
		return "", fmt.Errorf("create read-only API key: %w", err)
	}
	fmt.Printf("   ✓ API Key: %s (read-only)\n", apiKey2.Name)
	_ = rawKey2 // Suppress unused warning

	return rawKey, nil
}

func seedRoutingRules(ctx context.Context, repo *postgres.RoutingRuleRepository, policies map[string]*domain.EscalationPolicy) error {
	fmt.Println("\n🔀 Creating routing rules...")

	rules := []struct {
		Name        string
		Description string
		Priority    int
		Conditions  domain.RoutingConditions
		Actions     domain.RoutingActions
	}{
		{
			Name:        "Critical Alerts to Platform",
			Description: "Route P1 alerts from infrastructure sources to Platform team",
			Priority:    1,
			Conditions: domain.RoutingConditions{
				Match: "all",
				Conditions: []domain.RoutingCondition{
					{Field: "priority", Operator: "equals", Value: "P1"},
					{Field: "source", Operator: "contains", Value: "prometheus"},
				},
			},
			Actions: domain.RoutingActions{
				AssignTeamID:             &demoTeams["platform"].ID,
				AssignEscalationPolicyID: &policies["platform"].ID,
				AddTags:                  []string{"auto-routed", "infrastructure"},
			},
		},
		{
			Name:        "Database Alerts",
			Description: "Route database-related alerts to Backend team",
			Priority:    2,
			Conditions: domain.RoutingConditions{
				Match: "any",
				Conditions: []domain.RoutingCondition{
					{Field: "tags", Operator: "contains", Value: "database"},
					{Field: "tags", Operator: "contains", Value: "postgres"},
					{Field: "message", Operator: "contains", Value: "database"},
				},
			},
			Actions: domain.RoutingActions{
				AssignTeamID:             &demoTeams["backend"].ID,
				AssignEscalationPolicyID: &policies["backend"].ID,
				AddTags:                  []string{"auto-routed", "database"},
			},
		},
		{
			Name:        "Security Alerts Priority Boost",
			Description: "Upgrade security-related alerts to P2 minimum",
			Priority:    3,
			Conditions: domain.RoutingConditions{
				Match: "any",
				Conditions: []domain.RoutingCondition{
					{Field: "tags", Operator: "contains", Value: "security"},
					{Field: "source", Operator: "equals", Value: "SecurityHub"},
				},
			},
			Actions: domain.RoutingActions{
				SetPriority: strPtr("P2"),
				AddTags:     []string{"security-review"},
			},
		},
		{
			Name:        "Suppress Test Alerts",
			Description: "Suppress alerts from test/staging environments",
			Priority:    0, // Highest priority - check first
			Conditions: domain.RoutingConditions{
				Match: "any",
				Conditions: []domain.RoutingCondition{
					{Field: "tags", Operator: "contains", Value: "test"},
					{Field: "tags", Operator: "contains", Value: "staging"},
					{Field: "source", Operator: "contains", Value: "test"},
				},
			},
			Actions: domain.RoutingActions{
				Suppress: true,
			},
		},
	}

	for _, r := range rules {
		conditionsJSON, _ := json.Marshal(r.Conditions)
		actionsJSON, _ := json.Marshal(r.Actions)

		desc := r.Description
		rule := &domain.AlertRoutingRule{
			ID:             uuid.New(),
			OrganizationID: demoOrg.ID,
			Name:           r.Name,
			Description:    &desc,
			Priority:       r.Priority,
			Conditions:     conditionsJSON,
			Actions:        actionsJSON,
			Enabled:        true,
		}

		if err := repo.Create(ctx, rule); err != nil {
			return fmt.Errorf("create routing rule %s: %w", r.Name, err)
		}
		fmt.Printf("   ✓ Rule: %s (priority: %d)\n", r.Name, r.Priority)
	}

	return nil
}

func seedDNDSettings(ctx context.Context, repo *postgres.DNDSettingsRepository) error {
	fmt.Println("\n🔕 Creating DND settings...")

	// Create DND settings for Alice (weeknight DND)
	aliceSchedule := domain.DNDSchedule{
		Timezone: "America/New_York",
		Weekly: []domain.DNDTimeSlot{
			{Day: "monday", Start: "22:00", End: "08:00"},
			{Day: "tuesday", Start: "22:00", End: "08:00"},
			{Day: "wednesday", Start: "22:00", End: "08:00"},
			{Day: "thursday", Start: "22:00", End: "08:00"},
			{Day: "friday", Start: "22:00", End: "08:00"},
		},
	}
	aliceScheduleJSON, _ := json.Marshal(aliceSchedule)

	aliceDND := &domain.UserDNDSettings{
		ID:              uuid.New(),
		UserID:          demoUsers["alice"].ID,
		Enabled:         true,
		Schedule:        aliceScheduleJSON,
		Overrides:       json.RawMessage("[]"),
		AllowP1Override: true, // P1 alerts will still come through
	}

	if err := repo.Create(ctx, aliceDND); err != nil {
		return fmt.Errorf("create DND for alice: %w", err)
	}
	fmt.Printf("   ✓ DND: Alice (weeknights 10PM-8AM, P1 override enabled)\n")

	// Create DND settings for Bob (weekend DND)
	bobSchedule := domain.DNDSchedule{
		Timezone: "America/New_York",
		Weekly: []domain.DNDTimeSlot{
			{Day: "saturday", Start: "00:00", End: "23:59"},
			{Day: "sunday", Start: "00:00", End: "23:59"},
		},
	}
	bobScheduleJSON, _ := json.Marshal(bobSchedule)

	bobDND := &domain.UserDNDSettings{
		ID:              uuid.New(),
		UserID:          demoUsers["bob"].ID,
		Enabled:         true,
		Schedule:        bobScheduleJSON,
		Overrides:       json.RawMessage("[]"),
		AllowP1Override: true,
	}

	if err := repo.Create(ctx, bobDND); err != nil {
		return fmt.Errorf("create DND for bob: %w", err)
	}
	fmt.Printf("   ✓ DND: Bob (weekends, P1 override enabled)\n")

	return nil
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func strPtr(s string) *string {
	return &s
}
