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
	return &TestConfig{
		DatabaseURL: getEnvOrDefault("TEST_DATABASE_URL",
			"postgres://pulsar_test:pulsar_test_password@localhost:5434/pulsar_test?sslmode=disable"),
		JWTSecret: getEnvOrDefault("TEST_JWT_SECRET",
			"test_jwt_secret_at_least_32_characters_long_for_testing"),
		RefreshSecret: getEnvOrDefault("TEST_JWT_REFRESH_SECRET",
			"test_refresh_secret_at_least_32_characters_long_for_testing"),
		ServerPort: "0", // Random port for parallel tests
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
