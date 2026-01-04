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
	"github.com/pulsar/backend/internal/config"
	"github.com/pulsar/backend/internal/handler/rest"
	"github.com/pulsar/backend/internal/middleware"
	"github.com/pulsar/backend/internal/pkg/logger"
	"github.com/pulsar/backend/internal/repository/postgres"
	"github.com/pulsar/backend/internal/service"
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

	// Initialize services
	authService := service.NewAuthService(userRepo, orgRepo, cfg)
	alertService := service.NewAlertService(alertRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(orgRepo)

	// Initialize handlers
	authHandler := rest.NewAuthHandler(authService)
	alertHandler := rest.NewAlertHandler(alertService)
	teamHandler := rest.NewTeamHandler(teamService)
	userHandler := rest.NewUserHandler(userService)

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

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server stopped")
}
