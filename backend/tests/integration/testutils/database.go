package testutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TestDB wraps sqlx.DB with test utilities
type TestDB struct {
	*sqlx.DB
	migrationsPath string
}

// NewTestDB creates a new test database connection
func NewTestDB(databaseURL string) (*TestDB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Get migrations path relative to this file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}
	migrationsPath := filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrations")

	return &TestDB{
		DB:             db,
		migrationsPath: migrationsPath,
	}, nil
}

// Migrate runs all database migrations
func (tdb *TestDB) Migrate() error {
	// Read and execute migration files in order
	files, err := filepath.Glob(filepath.Join(tdb.migrationsPath, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := tdb.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

// TruncateAll truncates all tables in the correct order (respecting foreign keys)
func (tdb *TestDB) TruncateAll(ctx context.Context) error {
	// Tables in reverse dependency order to avoid FK violations
	tables := []string{
		"incident_alerts",
		"incident_timeline",
		"incident_responders",
		"incidents",
		"webhook_deliveries",
		"webhook_endpoints",
		"incoming_webhook_tokens",
		"notification_logs",
		"user_notification_preferences",
		"notification_channels",
		"alert_escalation_events",
		"escalation_targets",
		"escalation_rules",
		"escalation_policies",
		"schedule_overrides",
		"schedule_rotation_participants",
		"schedule_rotations",
		"schedules",
		"alert_routing_rules",
		"api_keys",
		"email_verifications",
		"user_dnd_settings",
		"team_invitations",
		"alerts",
		"team_members",
		"teams",
		"organization_users",
		"users",
		"organizations",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		if _, err := tdb.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}

// Reset drops all tables and re-runs migrations
func (tdb *TestDB) Reset() error {
	// Drop all tables
	tables := []string{
		"incident_alerts",
		"incident_timeline",
		"incident_responders",
		"incidents",
		"webhook_deliveries",
		"webhook_endpoints",
		"incoming_webhook_tokens",
		"notification_logs",
		"user_notification_preferences",
		"notification_channels",
		"alert_escalation_events",
		"escalation_targets",
		"escalation_rules",
		"escalation_policies",
		"schedule_overrides",
		"schedule_rotation_participants",
		"schedule_rotations",
		"schedules",
		"alert_routing_rules",
		"api_keys",
		"email_verifications",
		"user_dnd_settings",
		"team_invitations",
		"alerts",
		"team_members",
		"teams",
		"organization_users",
		"users",
		"organizations",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		if _, err := tdb.Exec(query); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	// Re-run migrations
	return tdb.Migrate()
}

// Close closes the database connection
func (tdb *TestDB) Close() error {
	return tdb.DB.Close()
}
