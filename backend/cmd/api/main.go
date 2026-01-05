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
	"github.com/nmn3m/pulsar/backend/internal/config"
	"github.com/nmn3m/pulsar/backend/internal/handler/rest"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/pkg/logger"
	"github.com/nmn3m/pulsar/backend/internal/repository/postgres"
	"github.com/nmn3m/pulsar/backend/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Server.Env)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Pulsar API server",
		zap.String("env", cfg.Server.Env),
		zap.String("port", cfg.Server.Port),
	)

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
	notificationRepo := postgres.NewNotificationRepository(db)
	incidentRepo := postgres.NewIncidentRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, orgRepo, cfg)
	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(orgRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, userRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	incidentService := service.NewIncidentService(incidentRepo)

	// Initialize alert notifier with dependencies
	alertNotifier := service.NewAlertNotifier(notificationService, userRepo, teamRepo, scheduleService)

	// Initialize alert and escalation services with notifier
	alertService := service.NewAlertService(alertRepo, alertNotifier)
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

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Setup router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

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

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop escalation worker
	escalationWorkerQuit <- true

	log.Info("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server stopped")
}
