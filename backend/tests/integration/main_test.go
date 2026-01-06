package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/nmn3m/pulsar/backend/tests/integration/testutils"
)

var (
	testServer   *testutils.TestServer
	testDB       *testutils.TestDB
	testFixtures *testutils.TestFixtures
)

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Load test configuration
	cfg := testutils.LoadTestConfig()

	// Setup test database
	var err error
	testDB, err = testutils.NewTestDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v\n"+
			"Make sure the test database is running: docker-compose -f docker-compose.test.yml up -d", err)
	}

	// Reset database (drop all tables and re-run migrations)
	if err := testDB.Reset(); err != nil {
		log.Fatalf("Failed to reset database: %v", err)
	}

	// Setup test server
	testServer, err = testutils.NewTestServer(testDB, cfg)
	if err != nil {
		log.Fatalf("Failed to create test server: %v", err)
	}

	// Create fixtures helper
	testFixtures = testutils.NewTestFixtures(testServer)

	// Run tests
	code := m.Run()

	// Cleanup
	testServer.Close()
	testDB.Close()

	os.Exit(code)
}

// newTestClient creates a new test client for the current test
func newTestClient(t *testing.T) *testutils.TestClient {
	return testutils.NewTestClient(t, testServer.URL())
}

// cleanDatabase truncates all tables before a test
func cleanDatabase(t *testing.T) {
	ctx := context.Background()
	if err := testDB.TruncateAll(ctx); err != nil {
		t.Fatalf("Failed to clean database: %v", err)
	}
}
