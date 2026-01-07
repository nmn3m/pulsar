package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// ============================================================================
// GET /api/v1/metrics/dashboard
// ============================================================================

func TestMetrics_GetDashboard_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, err := testFixtures.CreateUniqueUser(ctx)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	client.SetAuthToken(user.AccessToken)

	// Create test data
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert 1")
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert 2")
	testFixtures.CreateIncident(ctx, user.Organization.ID, user.User.ID, "Test Incident")

	resp := client.Get("/api/v1/metrics/dashboard")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify alerts section exists
	if result["alerts"] == nil {
		t.Error("Expected alerts section in dashboard metrics")
	}
	alerts := result["alerts"].(map[string]interface{})
	if alerts["total"].(float64) < 2 {
		t.Errorf("Expected at least 2 alerts, got %v", alerts["total"])
	}

	// Verify incidents section exists
	if result["incidents"] == nil {
		t.Error("Expected incidents section in dashboard metrics")
	}

	// Verify notifications section exists
	if result["notifications"] == nil {
		t.Error("Expected notifications section in dashboard metrics")
	}

	// Verify updated_at exists
	if result["updated_at"] == nil {
		t.Error("Expected updated_at in dashboard metrics")
	}
}

func TestMetrics_GetDashboard_WithTimeFilter(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create test alert
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	// Use time filter from last hour
	startTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	endTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

	resp := client.GetWithQuery("/api/v1/metrics/dashboard", map[string]string{
		"start_time": startTime,
		"end_time":   endTime,
	})
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	alerts := result["alerts"].(map[string]interface{})
	if alerts["total"].(float64) < 1 {
		t.Errorf("Expected at least 1 alert within time range, got %v", alerts["total"])
	}
}

func TestMetrics_GetDashboard_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	// Make request without auth
	resp := client.Get("/api/v1/metrics/dashboard")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/metrics/alerts
// ============================================================================

func TestMetrics_GetAlertMetrics_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create alerts with different priorities
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Critical Alert")
	testFixtures.CreateAlert(ctx, user.Organization.ID, "High Alert")
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Medium Alert")

	resp := client.Get("/api/v1/metrics/alerts")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify total count
	if result["total"].(float64) < 3 {
		t.Errorf("Expected at least 3 alerts, got %v", result["total"])
	}

	// Verify by_priority map exists
	if result["by_priority"] == nil {
		t.Error("Expected by_priority in alert metrics")
	}

	// Verify status counts exist
	if result["open"] == nil {
		t.Error("Expected open count in alert metrics")
	}
}

func TestMetrics_GetAlertMetrics_WithAcknowledged(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create and acknowledge an alert
	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	// Acknowledge the alert
	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/acknowledge", alert.ID), nil)
	client.AssertStatus(resp, http.StatusOK)

	resp = client.Get("/api/v1/metrics/alerts")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify acknowledged count
	if result["acknowledged"].(float64) < 1 {
		t.Errorf("Expected at least 1 acknowledged alert, got %v", result["acknowledged"])
	}
}

func TestMetrics_GetAlertMetrics_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/metrics/alerts")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/metrics/incidents
// ============================================================================

func TestMetrics_GetIncidentMetrics_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create incidents with different severities
	testFixtures.CreateIncident(ctx, user.Organization.ID, user.User.ID, "Critical Incident")
	testFixtures.CreateIncident(ctx, user.Organization.ID, user.User.ID, "High Incident")

	resp := client.Get("/api/v1/metrics/incidents")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify total count
	if result["total"].(float64) < 2 {
		t.Errorf("Expected at least 2 incidents, got %v", result["total"])
	}

	// Verify by_severity map exists
	if result["by_severity"] == nil {
		t.Error("Expected by_severity in incident metrics")
	}

	// Verify status counts exist
	if result["open"] == nil {
		t.Error("Expected open count in incident metrics")
	}
}

func TestMetrics_GetIncidentMetrics_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/metrics/incidents")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/metrics/notifications
// ============================================================================

func TestMetrics_GetNotificationMetrics_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/metrics/notifications")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify fields exist (counts may be 0 if no notifications sent)
	if result["total"] == nil {
		t.Error("Expected total in notification metrics")
	}
	if result["sent"] == nil {
		t.Error("Expected sent count in notification metrics")
	}
	if result["pending"] == nil {
		t.Error("Expected pending count in notification metrics")
	}
	if result["failed"] == nil {
		t.Error("Expected failed count in notification metrics")
	}
	if result["by_channel"] == nil {
		t.Error("Expected by_channel in notification metrics")
	}
}

func TestMetrics_GetNotificationMetrics_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/metrics/notifications")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/metrics/alerts/trend
// ============================================================================

func TestMetrics_GetAlertTrend_Daily(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create some alerts
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert 1")
	testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert 2")

	resp := client.GetWithQuery("/api/v1/metrics/alerts/trend", map[string]string{
		"period": "daily",
	})
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify period
	if result["period"] != "daily" {
		t.Errorf("Expected period 'daily', got %v", result["period"])
	}

	// Verify created array exists
	if result["created"] == nil {
		t.Error("Expected created array in alert trend")
	}

	// Verify closed array exists
	if result["closed"] == nil {
		t.Error("Expected closed array in alert trend")
	}
}

func TestMetrics_GetAlertTrend_Hourly(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.GetWithQuery("/api/v1/metrics/alerts/trend", map[string]string{
		"period": "hourly",
	})
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["period"] != "hourly" {
		t.Errorf("Expected period 'hourly', got %v", result["period"])
	}
}

func TestMetrics_GetAlertTrend_Weekly(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.GetWithQuery("/api/v1/metrics/alerts/trend", map[string]string{
		"period": "weekly",
	})
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["period"] != "weekly" {
		t.Errorf("Expected period 'weekly', got %v", result["period"])
	}
}

func TestMetrics_GetAlertTrend_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/metrics/alerts/trend")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/metrics/teams
// ============================================================================

func TestMetrics_GetTeamMetrics_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create a team
	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Metrics Test Team")

	// Create an alert and assign to team
	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Team Alert")

	assignReq := map[string]interface{}{
		"team_id": team.ID.String(),
	}
	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/assign", alert.ID), assignReq)
	client.AssertStatus(resp, http.StatusOK)

	resp = client.Get("/api/v1/metrics/teams")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Verify teams array exists
	if result["teams"] == nil {
		t.Error("Expected teams array in response")
	}

	teamsRaw := result["teams"]
	if teamsRaw == nil {
		t.Error("Expected teams array in metrics response")
		return
	}

	teams := teamsRaw.([]interface{})
	if len(teams) == 0 {
		t.Error("Expected at least one team in metrics")
		return
	}

	// Find our team and verify metrics
	found := false
	for _, teamData := range teams {
		tm := teamData.(map[string]interface{})
		if tm["team_id"] == team.ID.String() {
			found = true
			if tm["total_alerts"].(float64) < 1 {
				t.Errorf("Expected at least 1 alert for team, got %v", tm["total_alerts"])
			}
			break
		}
	}

	if !found {
		t.Error("Test team not found in metrics response")
	}
}

func TestMetrics_GetTeamMetrics_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create a team without any alerts assigned
	testFixtures.CreateTeam(ctx, user.Organization.ID, "Empty Team")

	resp := client.Get("/api/v1/metrics/teams")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	teamsRaw := result["teams"]
	if teamsRaw == nil {
		t.Error("Expected teams array in response")
		return
	}

	teams := teamsRaw.([]interface{})
	if len(teams) == 0 {
		t.Error("Expected at least one team in metrics")
	}
}

func TestMetrics_GetTeamMetrics_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/metrics/teams")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}
