package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// ============================================================================
// POST /api/v1/alerts
// ============================================================================

func TestAlerts_Create_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, err := testFixtures.CreateUniqueUser(ctx)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"source":      "api-test",
		"priority":    "P2",
		"message":     "Test alert message",
		"description": "Test alert description",
		"tags":        []string{"test", "api"},
	}

	resp := client.Post("/api/v1/alerts", reqBody)
	client.AssertStatus(resp, http.StatusCreated)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] == nil {
		t.Error("Expected id in response")
	}
	if result["message"] != "Test alert message" {
		t.Errorf("Expected message 'Test alert message', got %v", result["message"])
	}
	if result["priority"] != "P2" {
		t.Errorf("Expected priority 'P2', got %v", result["priority"])
	}
	if result["status"] != "open" {
		t.Errorf("Expected status 'open', got %v", result["status"])
	}
}

func TestAlerts_Create_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]interface{}{
		"source":   "api",
		"priority": "P2",
		"message":  "Test alert",
	}

	resp := client.Post("/api/v1/alerts", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestAlerts_Create_MissingRequiredFields(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Missing message
	reqBody := map[string]interface{}{
		"source":   "api",
		"priority": "P2",
	}

	resp := client.Post("/api/v1/alerts", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// GET /api/v1/alerts
// ============================================================================

func TestAlerts_List_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create some alerts
	for i := 0; i < 3; i++ {
		_, err := testFixtures.CreateAlert(ctx, user.Organization.ID, fmt.Sprintf("Alert %d", i))
		if err != nil {
			t.Fatalf("Failed to create alert: %v", err)
		}
	}

	resp := client.Get("/api/v1/alerts")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	alerts := result["alerts"].([]interface{})
	if len(alerts) != 3 {
		t.Errorf("Expected 3 alerts, got %d", len(alerts))
	}
}

func TestAlerts_List_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/alerts")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Handle nil or empty alerts array
	alertsRaw := result["alerts"]
	if alertsRaw != nil {
		alerts := alertsRaw.([]interface{})
		if len(alerts) != 0 {
			t.Errorf("Expected 0 alerts, got %d", len(alerts))
		}
	}
}

func TestAlerts_List_WithPagination(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create 5 alerts
	for i := 0; i < 5; i++ {
		testFixtures.CreateAlert(ctx, user.Organization.ID, fmt.Sprintf("Alert %d", i))
	}

	resp := client.GetWithQuery("/api/v1/alerts", map[string]string{
		"page":      "1",
		"page_size": "2",
	})
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	alerts := result["alerts"].([]interface{})
	if len(alerts) != 2 {
		t.Errorf("Expected 2 alerts on page 1, got %d", len(alerts))
	}
}

func TestAlerts_List_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/alerts")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/alerts/:id
// ============================================================================

func TestAlerts_Get_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, err := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")
	if err != nil {
		t.Fatalf("Failed to create alert: %v", err)
	}

	resp := client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] != alert.ID.String() {
		t.Errorf("Expected alert ID %s, got %v", alert.ID, result["id"])
	}
}

func TestAlerts_Get_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/alerts/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

func TestAlerts_Get_InvalidUUID(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/alerts/invalid-uuid")
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// PATCH /api/v1/alerts/:id
// ============================================================================

func TestAlerts_Update_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	reqBody := map[string]interface{}{
		"priority": "P1",
	}

	resp := client.Patch(fmt.Sprintf("/api/v1/alerts/%s", alert.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["priority"] != "P1" {
		t.Errorf("Expected priority 'P1', got %v", result["priority"])
	}
}

func TestAlerts_Update_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"priority": "P1",
	}

	resp := client.Patch("/api/v1/alerts/00000000-0000-0000-0000-000000000000", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// DELETE /api/v1/alerts/:id
// ============================================================================

func TestAlerts_Delete_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	resp := client.Delete(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	client.AssertStatus(resp, http.StatusOK)

	// Verify it's deleted
	resp = client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	client.ExpectStatus(resp, http.StatusNotFound)
}

func TestAlerts_Delete_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Delete("/api/v1/alerts/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// POST /api/v1/alerts/:id/acknowledge
// ============================================================================

func TestAlerts_Acknowledge_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/acknowledge", alert.ID), nil)
	client.AssertStatus(resp, http.StatusOK)

	// Verify status changed
	resp = client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["status"] != "acknowledged" {
		t.Errorf("Expected status 'acknowledged', got %v", result["status"])
	}
}

func TestAlerts_Acknowledge_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Post("/api/v1/alerts/00000000-0000-0000-0000-000000000000/acknowledge", nil)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// POST /api/v1/alerts/:id/close
// ============================================================================

func TestAlerts_Close_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	reqBody := map[string]string{
		"reason": "Resolved by test",
	}

	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/close", alert.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	// Verify status changed
	resp = client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["status"] != "closed" {
		t.Errorf("Expected status 'closed', got %v", result["status"])
	}
}

func TestAlerts_Close_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]string{
		"reason": "Test reason",
	}

	resp := client.Post("/api/v1/alerts/00000000-0000-0000-0000-000000000000/close", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// POST /api/v1/alerts/:id/snooze
// ============================================================================

func TestAlerts_Snooze_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	snoozeUntil := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	reqBody := map[string]string{
		"until": snoozeUntil,
	}

	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/snooze", alert.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	// Verify status changed
	resp = client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["status"] != "snoozed" {
		t.Errorf("Expected status 'snoozed', got %v", result["status"])
	}
}

func TestAlerts_Snooze_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]string{
		"until": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	}

	resp := client.Post("/api/v1/alerts/00000000-0000-0000-0000-000000000000/snooze", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// POST /api/v1/alerts/:id/assign
// ============================================================================

func TestAlerts_Assign_ToUser_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	reqBody := map[string]string{
		"user_id": user.User.ID.String(),
	}

	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/assign", alert.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	// Verify assignment
	resp = client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["assigned_to_user_id"] != user.User.ID.String() {
		t.Errorf("Expected assigned_to_user_id %s, got %v", user.User.ID, result["assigned_to_user_id"])
	}
}

func TestAlerts_Assign_ToTeam_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")
	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	reqBody := map[string]string{
		"team_id": team.ID.String(),
	}

	resp := client.Post(fmt.Sprintf("/api/v1/alerts/%s/assign", alert.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	// Verify assignment
	resp = client.Get(fmt.Sprintf("/api/v1/alerts/%s", alert.ID))
	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["assigned_to_team_id"] != team.ID.String() {
		t.Errorf("Expected assigned_to_team_id %s, got %v", team.ID, result["assigned_to_team_id"])
	}
}

func TestAlerts_Assign_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]string{
		"user_id": user.User.ID.String(),
	}

	resp := client.Post("/api/v1/alerts/00000000-0000-0000-0000-000000000000/assign", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}
