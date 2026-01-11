package testutils

import "os"

// TestConfig holds configuration for integration tests
type TestConfig struct {
	DatabaseURL   string
	JWTSecret     string
	RefreshSecret string
	ServerPort    string
}

// LoadTestConfig loads test configuration from environment variables with defaults
func LoadTestConfig() *TestConfig {
	// Check TEST_DATABASE_URL first, then DATABASE_URL, then default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://pulsar_test:pulsar_test_password@localhost:5434/pulsar_test?sslmode=disable"
	}

	// Check TEST_JWT_SECRET first, then JWT_SECRET, then default
	jwtSecret := os.Getenv("TEST_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if jwtSecret == "" {
		jwtSecret = "test_jwt_secret_at_least_32_characters_long_for_testing"
	}

	// Check TEST_JWT_REFRESH_SECRET first, then JWT_REFRESH_SECRET, then default
	refreshSecret := os.Getenv("TEST_JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	}
	if refreshSecret == "" {
		refreshSecret = "test_refresh_secret_at_least_32_characters_long_for_testing"
	}

	return &TestConfig{
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		RefreshSecret: refreshSecret,
		ServerPort:    "0", // Random port for parallel tests
	}
}
