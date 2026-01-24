package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	SMTP      SMTPConfig
	Email     EmailConfig
	Telemetry TelemetryConfig
}

// TelemetryConfig holds OpenTelemetry configuration
type TelemetryConfig struct {
	Enabled      bool
	ServiceName  string
	OTLPEndpoint string // e.g., "localhost:4317" for gRPC or "localhost:4318" for HTTP
	OTLPProtocol string // "grpc" or "http"
	Environment  string // e.g., "development", "production"
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string // Sender email address
	FromName string // Sender display name
	Enabled  bool
	UseTLS   bool
}

// EmailConfig holds email provider configuration
// Provider can be "smtp" (for development with Mailpit) or "resend" (for production)
type EmailConfig struct {
	Provider     string // "smtp" or "resend"
	Enabled      bool
	From         string // Sender email address
	FromName     string // Sender display name
	ResendAPIKey string // Resend API key (used when Provider is "resend")
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret        string
	RefreshSecret string
	AccessTTL     int // in minutes
	RefreshTTL    int // in days
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", ""),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AccessTTL:     getEnvInt("JWT_ACCESS_TTL", 60), // 60 minutes default
			RefreshTTL:    getEnvInt("JWT_REFRESH_TTL", 7), // 7 days
		},
		CORS: CORSConfig{
			AllowedOrigins: parseAllowedOrigins(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnvInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@pulsar.local"),
			FromName: getEnv("SMTP_FROM_NAME", "Pulsar"),
			Enabled:  getEnv("SMTP_ENABLED", "false") == "true",
			UseTLS:   getEnv("SMTP_USE_TLS", "true") == "true",
		},
		Email: EmailConfig{
			Provider:     getEnv("EMAIL_PROVIDER", "smtp"), // "smtp" for dev (Mailpit), "resend" for production
			Enabled:      getEnv("EMAIL_ENABLED", "false") == "true",
			From:         getEnv("EMAIL_FROM", "noreply@pulsar.local"),
			FromName:     getEnv("EMAIL_FROM_NAME", "Pulsar"),
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		},
		Telemetry: TelemetryConfig{
			Enabled:      getEnv("OTEL_ENABLED", "false") == "true",
			ServiceName:  getEnv("OTEL_SERVICE_NAME", "pulsar-backend"),
			OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			OTLPProtocol: getEnv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc"), // "grpc" or "http"
			Environment:  getEnv("OTEL_ENVIRONMENT", "development"),
		},
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	if len(c.JWT.RefreshSecret) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}
