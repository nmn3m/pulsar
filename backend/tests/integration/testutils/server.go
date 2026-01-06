package testutils

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/nmn3m/pulsar/backend/internal/config"
	"github.com/nmn3m/pulsar/backend/internal/handler/rest"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/repository/postgres"
	"github.com/nmn3m/pulsar/backend/internal/service"
	"go.uber.org/zap"
)

// TestServer wraps httptest.Server with all dependencies
type TestServer struct {
	Server *httptest.Server
	Router *gin.Engine
	DB     *TestDB
	Config *config.Config
	Logger *zap.Logger

	// Services exposed for direct manipulation in tests
	AuthService         *service.AuthService
	AlertService        *service.AlertService
	TeamService         *service.TeamService
	ScheduleService     *service.ScheduleService
	EscalationService   *service.EscalationService
	NotificationService *service.NotificationService
	IncidentService     *service.IncidentService
	WebhookService      *service.WebhookService
	UserService         *service.UserService
}

// NewTestServer creates a new test server with all dependencies wired up
func NewTestServer(testDB *TestDB, testCfg *TestConfig) (*TestServer, error) {
	gin.SetMode(gin.TestMode)

	// Build config from test config
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: testCfg.ServerPort,
			Env:  "test",
		},
		Database: config.DatabaseConfig{
			URL: testCfg.DatabaseURL,
		},
		JWT: config.JWTConfig{
			Secret:        testCfg.JWTSecret,
			RefreshSecret: testCfg.RefreshSecret,
			AccessTTL:     15,
			RefreshTTL:    7,
		},
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"*"},
		},
	}

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create postgres DB wrapper that matches the expected type
	db := &postgres.DB{DB: testDB.DB}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	alertRepo := postgres.NewAlertRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)
	escalationRepo := postgres.NewEscalationPolicyRepository(db)
	notificationRepo := postgres.NewNotificationRepository(testDB.DB)
	incidentRepo := postgres.NewIncidentRepository(testDB.DB)
	webhookRepo := postgres.NewWebhookRepository(testDB.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, orgRepo, cfg)
	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(orgRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, userRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	wsService := service.NewWebSocketService(logger)
	incidentService := service.NewIncidentService(incidentRepo, wsService)
	webhookService := service.NewWebhookService(webhookRepo, logger)

	// Initialize alert notifier with dependencies
	alertNotifier := service.NewAlertNotifier(notificationService, userRepo, teamRepo, scheduleService)

	// Initialize alert and escalation services with notifier
	alertService := service.NewAlertService(alertRepo, alertNotifier, wsService, webhookService)
	escalationService := service.NewEscalationService(escalationRepo, alertRepo, alertNotifier)

	// Initialize handlers
	authHandler := rest.NewAuthHandler(authService)
	alertHandler := rest.NewAlertHandler(alertService)
	teamHandler := rest.NewTeamHandler(teamService)
	userHandler := rest.NewUserHandler(userService)
	scheduleHandler := rest.NewScheduleHandler(scheduleService)
	escalationHandler := rest.NewEscalationHandler(escalationService)
	notificationHandler := rest.NewNotificationHandler(notificationService)
	incidentHandler := rest.NewIncidentHandler(incidentService)
	webhookHandler := rest.NewWebhookHandler(webhookService)
	incomingWebhookHandler := rest.NewIncomingWebhookHandler(webhookService, alertService, logger)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	// Setup routes (mirrors main.go)
	setupRoutes(router, authMiddleware, authHandler, alertHandler, teamHandler,
		userHandler, scheduleHandler, escalationHandler, notificationHandler,
		incidentHandler, webhookHandler, incomingWebhookHandler)

	// Create test server
	server := httptest.NewServer(router)

	return &TestServer{
		Server:              server,
		Router:              router,
		DB:                  testDB,
		Config:              cfg,
		Logger:              logger,
		AuthService:         authService,
		AlertService:        alertService,
		TeamService:         teamService,
		ScheduleService:     scheduleService,
		EscalationService:   escalationService,
		NotificationService: notificationService,
		IncidentService:     incidentService,
		WebhookService:      webhookService,
		UserService:         userService,
	}, nil
}

// setupRoutes configures all API routes (mirrors main.go)
func setupRoutes(
	router *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	authHandler *rest.AuthHandler,
	alertHandler *rest.AlertHandler,
	teamHandler *rest.TeamHandler,
	userHandler *rest.UserHandler,
	scheduleHandler *rest.ScheduleHandler,
	escalationHandler *rest.EscalationHandler,
	notificationHandler *rest.NotificationHandler,
	incidentHandler *rest.IncidentHandler,
	webhookHandler *rest.WebhookHandler,
	incomingWebhookHandler *rest.IncomingWebhookHandler,
) {
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/auth/me", authHandler.GetMe)

			// User routes
			protected.GET("/users", userHandler.ListOrganizationUsers)

			// Alert routes
			alerts := protected.Group("/alerts")
			{
				alerts.GET("", alertHandler.List)
				alerts.POST("", alertHandler.Create)
				alerts.GET("/:id", alertHandler.Get)
				alerts.PATCH("/:id", alertHandler.Update)
				alerts.DELETE("/:id", alertHandler.Delete)
				alerts.POST("/:id/acknowledge", alertHandler.Acknowledge)
				alerts.POST("/:id/close", alertHandler.Close)
				alerts.POST("/:id/snooze", alertHandler.Snooze)
				alerts.POST("/:id/assign", alertHandler.Assign)
			}

			// Team routes
			teams := protected.Group("/teams")
			{
				teams.GET("", teamHandler.List)
				teams.POST("", teamHandler.Create)
				teams.GET("/:id", teamHandler.Get)
				teams.PATCH("/:id", teamHandler.Update)
				teams.DELETE("/:id", teamHandler.Delete)
				teams.POST("/:id/members", teamHandler.AddMember)
				teams.GET("/:id/members", teamHandler.ListMembers)
				teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
				teams.PATCH("/:id/members/:userId", teamHandler.UpdateMemberRole)
			}

			// Schedule routes
			schedules := protected.Group("/schedules")
			{
				schedules.GET("", scheduleHandler.List)
				schedules.POST("", scheduleHandler.Create)
				schedules.GET("/:id", scheduleHandler.Get)
				schedules.PATCH("/:id", scheduleHandler.Update)
				schedules.DELETE("/:id", scheduleHandler.Delete)
				schedules.GET("/:id/oncall", scheduleHandler.GetOnCall)

				// Rotation routes
				schedules.GET("/:id/rotations", scheduleHandler.ListRotations)
				schedules.POST("/:id/rotations", scheduleHandler.CreateRotation)
				schedules.GET("/:id/rotations/:rotationId", scheduleHandler.GetRotation)
				schedules.PATCH("/:id/rotations/:rotationId", scheduleHandler.UpdateRotation)
				schedules.DELETE("/:id/rotations/:rotationId", scheduleHandler.DeleteRotation)

				// Participant routes
				schedules.GET("/:id/rotations/:rotationId/participants", scheduleHandler.ListParticipants)
				schedules.POST("/:id/rotations/:rotationId/participants", scheduleHandler.AddParticipant)
				schedules.DELETE("/:id/rotations/:rotationId/participants/:userId", scheduleHandler.RemoveParticipant)
				schedules.PUT("/:id/rotations/:rotationId/participants/reorder", scheduleHandler.ReorderParticipants)

				// Override routes
				schedules.GET("/:id/overrides", scheduleHandler.ListOverrides)
				schedules.POST("/:id/overrides", scheduleHandler.CreateOverride)
				schedules.GET("/:id/overrides/:overrideId", scheduleHandler.GetOverride)
				schedules.PATCH("/:id/overrides/:overrideId", scheduleHandler.UpdateOverride)
				schedules.DELETE("/:id/overrides/:overrideId", scheduleHandler.DeleteOverride)
			}

			// Escalation policy routes
			escalations := protected.Group("/escalation-policies")
			{
				escalations.GET("", escalationHandler.List)
				escalations.POST("", escalationHandler.Create)
				escalations.GET("/:id", escalationHandler.Get)
				escalations.PATCH("/:id", escalationHandler.Update)
				escalations.DELETE("/:id", escalationHandler.Delete)

				// Rule routes
				escalations.GET("/:id/rules", escalationHandler.ListRules)
				escalations.POST("/:id/rules", escalationHandler.CreateRule)
				escalations.GET("/:id/rules/:ruleId", escalationHandler.GetRule)
				escalations.PATCH("/:id/rules/:ruleId", escalationHandler.UpdateRule)
				escalations.DELETE("/:id/rules/:ruleId", escalationHandler.DeleteRule)

				// Target routes
				escalations.GET("/:id/rules/:ruleId/targets", escalationHandler.ListTargets)
				escalations.POST("/:id/rules/:ruleId/targets", escalationHandler.AddTarget)
				escalations.DELETE("/:id/rules/:ruleId/targets/:targetId", escalationHandler.RemoveTarget)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				// Channel routes
				notifications.GET("/channels", notificationHandler.ListChannels)
				notifications.POST("/channels", notificationHandler.CreateChannel)
				notifications.GET("/channels/:id", notificationHandler.GetChannel)
				notifications.PATCH("/channels/:id", notificationHandler.UpdateChannel)
				notifications.DELETE("/channels/:id", notificationHandler.DeleteChannel)

				// User preference routes
				notifications.GET("/preferences", notificationHandler.ListUserPreferences)
				notifications.POST("/preferences", notificationHandler.CreatePreference)
				notifications.GET("/preferences/:id", notificationHandler.GetPreference)
				notifications.PATCH("/preferences/:id", notificationHandler.UpdatePreference)
				notifications.DELETE("/preferences/:id", notificationHandler.DeletePreference)

				// Sending notifications
				notifications.POST("/send", notificationHandler.SendNotification)

				// Notification logs
				notifications.GET("/logs", notificationHandler.ListLogs)
				notifications.GET("/logs/:id", notificationHandler.GetLog)
				notifications.GET("/logs/user/me", notificationHandler.ListLogsByUser)
				notifications.GET("/logs/alert/:alertId", notificationHandler.ListLogsByAlert)
			}

			// Incident routes
			incidents := protected.Group("/incidents")
			{
				incidents.GET("", incidentHandler.List)
				incidents.POST("", incidentHandler.Create)
				incidents.GET("/:id", incidentHandler.GetWithDetails)
				incidents.PATCH("/:id", incidentHandler.Update)
				incidents.DELETE("/:id", incidentHandler.Delete)

				// Responder routes
				incidents.GET("/:id/responders", incidentHandler.ListResponders)
				incidents.POST("/:id/responders", incidentHandler.AddResponder)
				incidents.DELETE("/:id/responders/:responderId", incidentHandler.RemoveResponder)
				incidents.PATCH("/:id/responders/:responderId", incidentHandler.UpdateResponderRole)

				// Timeline routes
				incidents.GET("/:id/timeline", incidentHandler.GetTimeline)
				incidents.POST("/:id/notes", incidentHandler.AddNote)

				// Alert linking routes
				incidents.GET("/:id/alerts", incidentHandler.ListAlerts)
				incidents.POST("/:id/alerts", incidentHandler.LinkAlert)
				incidents.DELETE("/:id/alerts/:alertId", incidentHandler.UnlinkAlert)
			}

			// Webhook routes
			webhooks := protected.Group("/webhooks")
			{
				webhooks.GET("/endpoints", webhookHandler.ListEndpoints)
				webhooks.POST("/endpoints", webhookHandler.CreateEndpoint)
				webhooks.GET("/endpoints/:id", webhookHandler.GetEndpoint)
				webhooks.PATCH("/endpoints/:id", webhookHandler.UpdateEndpoint)
				webhooks.DELETE("/endpoints/:id", webhookHandler.DeleteEndpoint)

				webhooks.GET("/deliveries", webhookHandler.ListDeliveries)

				webhooks.GET("/incoming", webhookHandler.ListIncomingTokens)
				webhooks.POST("/incoming", webhookHandler.CreateIncomingToken)
				webhooks.DELETE("/incoming/:id", webhookHandler.DeleteIncomingToken)
			}
		}

		// Public incoming webhook route (no auth required)
		v1.POST("/webhook/:token", incomingWebhookHandler.ReceiveWebhook)
	}
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	ts.Server.Close()
	ts.Logger.Sync()
}

// URL returns the base URL of the test server
func (ts *TestServer) URL() string {
	return ts.Server.URL
}
