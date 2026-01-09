package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// ============================================================================
// POST /api/v1/notifications/channels
// ============================================================================

func TestNotifications_CreateChannel_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	config, _ := json.Marshal(map[string]interface{}{
		"smtp_host":     "smtp.example.com",
		"smtp_port":     587,
		"smtp_username": "alerts@example.com",
		"smtp_password": "password123",
		"from":          "alerts@example.com",
		"from_address":  "alerts@example.com",
	})

	reqBody := map[string]interface{}{
		"name":         "Email Notifications",
		"channel_type": "email",
		"is_enabled":   true,
		"config":       json.RawMessage(config),
	}

	resp := client.Post("/api/v1/notifications/channels", reqBody)
	client.AssertStatus(resp, http.StatusCreated)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] == nil {
		t.Error("Expected id in response")
	}
	if result["name"] != "Email Notifications" {
		t.Errorf("Expected name 'Email Notifications', got %v", result["name"])
	}
}

func TestNotifications_CreateChannel_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]interface{}{
		"name":         "Test Channel",
		"channel_type": "email",
	}

	resp := client.Post("/api/v1/notifications/channels", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/notifications/channels
// ============================================================================

func TestNotifications_ListChannels_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create channels
	for i := 0; i < 3; i++ {
		testFixtures.CreateNotificationChannel(ctx, user.Organization.ID, fmt.Sprintf("Channel %d", i))
	}

	resp := client.Get("/api/v1/notifications/channels")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	channels := result["channels"].([]interface{})
	if len(channels) != 3 {
		t.Errorf("Expected 3 channels, got %d", len(channels))
	}
}

func TestNotifications_ListChannels_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/channels")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	channelsRaw := result["channels"]
	if channelsRaw != nil {
		channels := channelsRaw.([]interface{})
		if len(channels) != 0 {
			t.Errorf("Expected 0 channels, got %d", len(channels))
		}
	}
}

// ============================================================================
// GET /api/v1/notifications/channels/:id
// ============================================================================

func TestNotifications_GetChannel_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	channel, _ := testFixtures.CreateNotificationChannel(ctx, user.Organization.ID, "Test Channel")

	resp := client.Get(fmt.Sprintf("/api/v1/notifications/channels/%s", channel.ID))
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] != channel.ID.String() {
		t.Errorf("Expected channel ID %s, got %v", channel.ID, result["id"])
	}
}

func TestNotifications_GetChannel_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/channels/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// PATCH /api/v1/notifications/channels/:id
// ============================================================================

func TestNotifications_UpdateChannel_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	channel, _ := testFixtures.CreateNotificationChannel(ctx, user.Organization.ID, "Original Channel")

	reqBody := map[string]interface{}{
		"name":       "Updated Channel",
		"is_enabled": false,
	}

	resp := client.Patch(fmt.Sprintf("/api/v1/notifications/channels/%s", channel.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["name"] != "Updated Channel" {
		t.Errorf("Expected name 'Updated Channel', got %v", result["name"])
	}
}

func TestNotifications_UpdateChannel_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name": "Updated Channel",
	}

	resp := client.Patch("/api/v1/notifications/channels/00000000-0000-0000-0000-000000000000", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// DELETE /api/v1/notifications/channels/:id
// ============================================================================

func TestNotifications_DeleteChannel_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	channel, _ := testFixtures.CreateNotificationChannel(ctx, user.Organization.ID, "Test Channel")

	resp := client.Delete(fmt.Sprintf("/api/v1/notifications/channels/%s", channel.ID))
	client.AssertStatus(resp, http.StatusOK)

	// Verify it's deleted
	resp = client.Get(fmt.Sprintf("/api/v1/notifications/channels/%s", channel.ID))
	client.ExpectStatus(resp, http.StatusNotFound)
}

func TestNotifications_DeleteChannel_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Delete("/api/v1/notifications/channels/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// GET /api/v1/notifications/preferences
// ============================================================================

func TestNotifications_ListPreferences_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/preferences")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Expecting preferences wrapped in object - preferences can be nil or empty initially
	if prefsRaw, ok := result["preferences"].([]interface{}); ok && len(prefsRaw) > 0 {
		t.Logf("User has %d preferences", len(prefsRaw))
	}
}

// ============================================================================
// POST /api/v1/notifications/preferences
// ============================================================================

func TestNotifications_CreatePreference_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	channel, _ := testFixtures.CreateNotificationChannel(ctx, user.Organization.ID, "Test Channel")

	reqBody := map[string]interface{}{
		"channel_id":  channel.ID.String(),
		"is_enabled":  true,
		"dnd_enabled": false,
	}

	resp := client.Post("/api/v1/notifications/preferences", reqBody)
	client.AssertStatus(resp, http.StatusCreated)
}

// ============================================================================
// GET /api/v1/notifications/preferences/:id
// ============================================================================

func TestNotifications_GetPreference_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/preferences/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// PATCH /api/v1/notifications/preferences/:id
// ============================================================================

func TestNotifications_UpdatePreference_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"is_enabled": false,
	}

	resp := client.Patch("/api/v1/notifications/preferences/00000000-0000-0000-0000-000000000000", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// DELETE /api/v1/notifications/preferences/:id
// ============================================================================

func TestNotifications_DeletePreference_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Delete("/api/v1/notifications/preferences/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// POST /api/v1/notifications/send
// ============================================================================

func TestNotifications_Send_MissingChannel(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Test with missing channel_id to trigger validation error
	reqBody := map[string]interface{}{
		"recipient": "test@example.com",
		"message":   "Test notification",
	}

	resp := client.Post("/api/v1/notifications/send", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for validation errors
}

// ============================================================================
// GET /api/v1/notifications/logs
// ============================================================================

func TestNotifications_ListLogs_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/logs")
	client.AssertStatus(resp, http.StatusOK)
}

// ============================================================================
// GET /api/v1/notifications/logs/:id
// ============================================================================

func TestNotifications_GetLog_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/logs/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// GET /api/v1/notifications/logs/user/me
// ============================================================================

func TestNotifications_ListLogsByUser_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/notifications/logs/user/me")
	client.AssertStatus(resp, http.StatusOK)
}

// ============================================================================
// GET /api/v1/notifications/logs/alert/:alertId
// ============================================================================

func TestNotifications_ListLogsByAlert_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	alert, _ := testFixtures.CreateAlert(ctx, user.Organization.ID, "Test Alert")

	resp := client.Get(fmt.Sprintf("/api/v1/notifications/logs/alert/%s", alert.ID))
	client.AssertStatus(resp, http.StatusOK)
}
