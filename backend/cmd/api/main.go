// @title           Pulsar API
// @version         1.0
// @description     Pulsar is an incident management and alerting platform. This API provides endpoints for managing alerts, incidents, teams, schedules, escalation policies, notifications, and webhooks.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/nmn3m/pulsar
// @contact.email  support@pulsar.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey APIKeyAuth
// @in header
// @name X-API-Key
// @description API key for programmatic access. Generate keys at /api/v1/api-keys.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"

	_ "github.com/nmn3m/pulsar/backend/docs"
	"github.com/nmn3m/pulsar/backend/internal/config"
	"github.com/nmn3m/pulsar/backend/internal/handler/rest"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/pkg/logger"
	"github.com/nmn3m/pulsar/backend/internal/pkg/telemetry"
	"github.com/nmn3m/pulsar/backend/internal/repository/postgres"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize OpenTelemetry if enabled (must be before logger for OTEL logs)
	var otelTelemetry *telemetry.Telemetry
	if cfg.Telemetry.Enabled {
		otelTelemetry, err = telemetry.Initialize(context.Background(), telemetry.Config{
			ServiceName:  cfg.Telemetry.ServiceName,
			Environment:  cfg.Telemetry.Environment,
			OTLPEndpoint: cfg.Telemetry.OTLPEndpoint,
			OTLPProtocol: cfg.Telemetry.OTLPProtocol,
		})
		if err != nil {
			fmt.Printf("Failed to initialize telemetry: %v\n", err)
			os.Exit(1)
		}
		defer otelTelemetry.Shutdown(context.Background())
	}

	// Initialize logger (with OTEL integration if telemetry is enabled)
	var log *zap.Logger
	if cfg.Telemetry.Enabled {
		log, err = logger.NewWithOTEL(cfg.Server.Env)
	} else {
		log, err = logger.New(cfg.Server.Env)
	}
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Pulsar API server",
		zap.String("env", cfg.Server.Env),
		zap.String("port", cfg.Server.Port),
	)

	if cfg.Telemetry.Enabled {
		log.Info("OpenTelemetry initialized",
			zap.String("service", cfg.Telemetry.ServiceName),
			zap.String("endpoint", cfg.Telemetry.OTLPEndpoint),
			zap.String("protocol", cfg.Telemetry.OTLPProtocol),
		)
	} else {
		log.Info("OpenTelemetry disabled")
	}

	// Connect to database
	db, err := postgres.NewDB(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	log.Info("Connected to database")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	alertRepo := postgres.NewAlertRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)
	escalationRepo := postgres.NewEscalationPolicyRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db.DB)
	incidentRepo := postgres.NewIncidentRepository(db.DB)
	webhookRepo := postgres.NewWebhookRepository(db.DB)
	apiKeyRepo := postgres.NewAPIKeyRepository(db.DB)
	metricsRepo := postgres.NewMetricsRepository(db.DB)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(db)
	routingRepo := postgres.NewRoutingRuleRepository(db)
	dndRepo := postgres.NewDNDSettingsRepository(db)
	invitationRepo := postgres.NewTeamInvitationRepo(db)

	// Initialize email service (for OTP verification and team invitations)
	var emailService *service.EmailService
	var emailVerificationService *service.EmailVerificationService
	if cfg.Email.Enabled || cfg.SMTP.Enabled {
		emailService = service.NewEmailService(&cfg.Email, &cfg.SMTP)
		emailVerificationService = service.NewEmailVerificationService(emailVerificationRepo, userRepo, emailService)
		if cfg.Email.Provider == "resend" {
			log.Info("Email service enabled with Resend provider",
				zap.String("from", cfg.Email.From),
			)
		} else {
			log.Info("Email service enabled with SMTP provider",
				zap.String("smtp_host", cfg.SMTP.Host),
				zap.Int("smtp_port", cfg.SMTP.Port),
			)
		}
	} else {
		log.Info("Email service disabled (not configured)")
	}

	// Initialize services
	authService := service.NewAuthService(userRepo, orgRepo, cfg, emailVerificationService)
	teamService := service.NewTeamService(teamRepo, userRepo)
	teamService.SetInvitationRepo(invitationRepo)
	if emailService != nil {
		teamService.SetEmailService(emailService)
	}
	userService := service.NewUserService(orgRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, userRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	wsService := service.NewWebSocketService(log)
	incidentService := service.NewIncidentService(incidentRepo, wsService)
	webhookService := service.NewWebhookService(webhookRepo, log)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)
	metricsService := service.NewMetricsService(metricsRepo)

	// Initialize DND and routing services
	dndService := service.NewDNDService(dndRepo)
	routingService := service.NewRoutingService(routingRepo)

	// Initialize alert notifier with dependencies (including DND service for quiet hours)
	alertNotifier := service.NewAlertNotifier(notificationService, userRepo, teamRepo, scheduleService, dndService)

	// Initialize alert and escalation services with notifier
	alertService := service.NewAlertService(alertRepo, alertNotifier, wsService, webhookService)
	escalationService := service.NewEscalationService(escalationRepo, alertRepo, alertNotifier)

	// Initialize handlers
	authHandler := rest.NewAuthHandler(authService, emailVerificationService)
	alertHandler := rest.NewAlertHandler(alertService)
	teamHandler := rest.NewTeamHandler(teamService)
	userHandler := rest.NewUserHandler(userService)
	scheduleHandler := rest.NewScheduleHandler(scheduleService)
	escalationHandler := rest.NewEscalationHandler(escalationService)
	notificationHandler := rest.NewNotificationHandler(notificationService)
	incidentHandler := rest.NewIncidentHandler(incidentService)
	wsHandler := rest.NewWebSocketHandler(wsService, log)
	webhookHandler := rest.NewWebhookHandler(webhookService)
	incomingWebhookHandler := rest.NewIncomingWebhookHandler(webhookService, alertService, log)
	apiKeyHandler := rest.NewAPIKeyHandler(apiKeyService)
	metricsHandler := rest.NewMetricsHandler(metricsService)
	routingHandler := rest.NewRoutingHandler(routingService)
	dndHandler := rest.NewDNDHandler(dndService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(apiKeyService)
	combinedAuth := middleware.NewCombinedAuthMiddleware(authMiddleware, apiKeyMiddleware)

	// Setup router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Add OpenTelemetry middleware if enabled
	if cfg.Telemetry.Enabled {
		router.Use(middleware.OTelMiddleware(cfg.Telemetry.ServiceName))
		router.Use(middleware.OTelMetricsMiddleware())
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/resend-otp", authHandler.ResendOTP)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/auth/me", authHandler.GetMe)

			// API Key routes
			apiKeys := protected.Group("/api-keys")
			{
				apiKeys.GET("/scopes", apiKeyHandler.GetScopes)
				apiKeys.GET("", apiKeyHandler.List)
				apiKeys.POST("", apiKeyHandler.Create)
				apiKeys.GET("/all", apiKeyHandler.ListAll)
				apiKeys.GET("/:id", apiKeyHandler.Get)
				apiKeys.PATCH("/:id", apiKeyHandler.Update)
				apiKeys.DELETE("/:id", apiKeyHandler.Delete)
				apiKeys.POST("/:id/revoke", apiKeyHandler.Revoke)
			}

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
				teams.POST("/:id/invite", teamHandler.InviteMember)
				teams.GET("/:id/invitations", teamHandler.ListInvitations)
				teams.DELETE("/:id/invitations/:invitationId", teamHandler.CancelInvitation)
				teams.POST("/:id/invitations/:invitationId/resend", teamHandler.ResendInvitation)
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

			// Routing rules routes
			routing := protected.Group("/routing-rules")
			{
				routing.GET("", routingHandler.List)
				routing.POST("", routingHandler.Create)
				routing.PUT("/reorder", routingHandler.Reorder)
				routing.GET("/:id", routingHandler.Get)
				routing.PATCH("/:id", routingHandler.Update)
				routing.DELETE("/:id", routingHandler.Delete)
			}

			// User DND (Do Not Disturb) routes
			usersDND := protected.Group("/users/me/dnd")
			{
				usersDND.GET("", dndHandler.GetDNDSettings)
				usersDND.PUT("", dndHandler.UpdateDNDSettings)
				usersDND.DELETE("", dndHandler.DeleteDNDSettings)
				usersDND.GET("/status", dndHandler.CheckDNDStatus)
				usersDND.POST("/overrides", dndHandler.AddDNDOverride)
				usersDND.DELETE("/overrides/:index", dndHandler.RemoveDNDOverride)
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

			// Metrics routes
			metrics := protected.Group("/metrics")
			{
				metrics.GET("/dashboard", metricsHandler.GetDashboard)
				metrics.GET("/alerts", metricsHandler.GetAlertMetrics)
				metrics.GET("/alerts/trend", metricsHandler.GetAlertTrend)
				metrics.GET("/incidents", metricsHandler.GetIncidentMetrics)
				metrics.GET("/notifications", metricsHandler.GetNotificationMetrics)
				metrics.GET("/teams", metricsHandler.GetTeamMetrics)
			}

			// WebSocket route
			protected.GET("/ws", wsHandler.HandleWebSocket)
			protected.GET("/ws/stats", wsHandler.GetStats)

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

		// API-key or JWT authenticated routes (for programmatic access)
		// These routes accept either Bearer token or X-API-Key header
		apiAuth := v1.Group("")
		apiAuth.Use(combinedAuth.RequireAuth())
		{
			// Alert management via API key
			apiAlerts := apiAuth.Group("/v2/alerts")
			{
				apiAlerts.GET("", alertHandler.List)
				apiAlerts.POST("", alertHandler.Create)
				apiAlerts.GET("/:id", alertHandler.Get)
				apiAlerts.PATCH("/:id", alertHandler.Update)
				apiAlerts.POST("/:id/acknowledge", alertHandler.Acknowledge)
				apiAlerts.POST("/:id/close", alertHandler.Close)
			}

			// Incident management via API key
			apiIncidents := apiAuth.Group("/v2/incidents")
			{
				apiIncidents.GET("", incidentHandler.List)
				apiIncidents.POST("", incidentHandler.Create)
				apiIncidents.GET("/:id", incidentHandler.GetWithDetails)
				apiIncidents.PATCH("/:id", incidentHandler.Update)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start WebSocket hub
	go wsService.Run()
	log.Info("WebSocket hub started")

	// Start server in a goroutine
	go func() {
		log.Info("Server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Start background worker for processing escalations
	escalationWorkerQuit := make(chan bool)
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Process escalations every 30 seconds
		defer ticker.Stop()

		log.Info("Escalation worker started")

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				if err := escalationService.ProcessPendingEscalations(ctx); err != nil {
					log.Error("Failed to process pending escalations", zap.Error(err))
				}
			case <-escalationWorkerQuit:
				log.Info("Escalation worker stopped")
				return
			}
		}
	}()

	// Start background worker for processing webhook deliveries
	webhookWorkerQuit := make(chan bool)
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Process webhook deliveries every 30 seconds
		defer ticker.Stop()

		log.Info("Webhook delivery worker started")

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				if err := webhookService.ProcessPendingDeliveries(ctx); err != nil {
					log.Error("Failed to process pending webhook deliveries", zap.Error(err))
				}
			case <-webhookWorkerQuit:
				log.Info("Webhook delivery worker stopped")
				return
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop background workers
	escalationWorkerQuit <- true
	webhookWorkerQuit <- true

	log.Info("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server stopped")
}
