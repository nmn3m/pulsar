package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

// ============================================================================
// POST /api/v1/webhooks/endpoints
// ============================================================================

func TestWebhooks_CreateEndpoint_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name":        "Slack Notifications",
		"url":         "https://hooks.slack.com/services/xxx/yyy/zzz",
		"event_types": []string{"alert.created", "alert.closed"},
		"is_enabled":  true,
	}

	resp := client.Post("/api/v1/webhooks/endpoints", reqBody)
	client.AssertStatus(resp, http.StatusCreated)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] == nil {
		t.Error("Expected id in response")
	}
	if result["name"] != "Slack Notifications" {
		t.Errorf("Expected name 'Slack Notifications', got %v", result["name"])
	}
}

func TestWebhooks_CreateEndpoint_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]interface{}{
		"name": "Test Webhook",
		"url":  "https://example.com/webhook",
	}

	resp := client.Post("/api/v1/webhooks/endpoints", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestWebhooks_CreateEndpoint_MissingURL(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name": "Test Webhook",
	}

	resp := client.Post("/api/v1/webhooks/endpoints", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// GET /api/v1/webhooks/endpoints
// ============================================================================

func TestWebhooks_ListEndpoints_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create endpoints
	for i := 0; i < 3; i++ {
		testFixtures.CreateWebhookEndpoint(ctx, user.Organization.ID,
			fmt.Sprintf("Webhook %d", i),
			fmt.Sprintf("https://example.com/webhook%d", i))
	}

	resp := client.Get("/api/v1/webhooks/endpoints")
	client.AssertStatus(resp, http.StatusOK)

	var result []interface{}
	client.ParseJSON(resp, &result)

	if len(result) != 3 {
		t.Errorf("Expected 3 endpoints, got %d", len(result))
	}
}

func TestWebhooks_ListEndpoints_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/webhooks/endpoints")
	client.AssertStatus(resp, http.StatusOK)

	var result []interface{}
	client.ParseJSON(resp, &result)

	if len(result) != 0 {
		t.Errorf("Expected 0 endpoints, got %d", len(result))
	}
}

// ============================================================================
// GET /api/v1/webhooks/endpoints/:id
// ============================================================================

func TestWebhooks_GetEndpoint_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	endpoint, _ := testFixtures.CreateWebhookEndpoint(ctx, user.Organization.ID, "Test Webhook", "https://example.com/webhook")

	resp := client.Get(fmt.Sprintf("/api/v1/webhooks/endpoints/%s", endpoint.ID))
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] != endpoint.ID.String() {
		t.Errorf("Expected endpoint ID %s, got %v", endpoint.ID, result["id"])
	}
}

func TestWebhooks_GetEndpoint_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/webhooks/endpoints/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// PATCH /api/v1/webhooks/endpoints/:id
// ============================================================================

func TestWebhooks_UpdateEndpoint_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	endpoint, _ := testFixtures.CreateWebhookEndpoint(ctx, user.Organization.ID, "Original Webhook", "https://example.com/webhook")

	reqBody := map[string]interface{}{
		"name":       "Updated Webhook",
		"is_enabled": false,
	}

	resp := client.Patch(fmt.Sprintf("/api/v1/webhooks/endpoints/%s", endpoint.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["name"] != "Updated Webhook" {
		t.Errorf("Expected name 'Updated Webhook', got %v", result["name"])
	}
}

func TestWebhooks_UpdateEndpoint_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name": "Updated Webhook",
	}

	resp := client.Patch("/api/v1/webhooks/endpoints/00000000-0000-0000-0000-000000000000", reqBody)
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// DELETE /api/v1/webhooks/endpoints/:id
// ============================================================================

func TestWebhooks_DeleteEndpoint_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	endpoint, _ := testFixtures.CreateWebhookEndpoint(ctx, user.Organization.ID, "Test Webhook", "https://example.com/webhook")

	resp := client.Delete(fmt.Sprintf("/api/v1/webhooks/endpoints/%s", endpoint.ID))
	client.AssertStatus(resp, http.StatusOK)

	// Verify it's deleted
	resp = client.Get(fmt.Sprintf("/api/v1/webhooks/endpoints/%s", endpoint.ID))
	client.ExpectStatus(resp, http.StatusNotFound)
}

func TestWebhooks_DeleteEndpoint_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Delete("/api/v1/webhooks/endpoints/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// GET /api/v1/webhooks/deliveries
// ============================================================================

func TestWebhooks_ListDeliveries_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/webhooks/deliveries")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Handle nil or empty deliveries array
	deliveriesRaw := result["deliveries"]
	if deliveriesRaw != nil {
		deliveries := deliveriesRaw.([]interface{})
		// Just verify it's an array, can be empty initially
		_ = deliveries
	}
}

func TestWebhooks_ListDeliveries_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/webhooks/deliveries")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// POST /api/v1/webhooks/incoming
// ============================================================================

func TestWebhooks_CreateIncomingToken_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name":             "Prometheus Alerts",
		"integration_type": "generic",
		"description":      "Token for receiving Prometheus alerts",
	}

	resp := client.Post("/api/v1/webhooks/incoming", reqBody)
	client.AssertStatus(resp, http.StatusCreated)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] == nil {
		t.Error("Expected id in response")
	}
	if result["token"] == nil {
		t.Error("Expected token in response")
	}
	if result["name"] != "Prometheus Alerts" {
		t.Errorf("Expected name 'Prometheus Alerts', got %v", result["name"])
	}
}

func TestWebhooks_CreateIncomingToken_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]interface{}{
		"name":             "Test Token",
		"integration_type": "generic",
	}

	resp := client.Post("/api/v1/webhooks/incoming", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// GET /api/v1/webhooks/incoming
// ============================================================================

func TestWebhooks_ListIncomingTokens_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create tokens
	for i := 0; i < 3; i++ {
		testFixtures.CreateIncomingWebhookToken(ctx, user.Organization.ID, fmt.Sprintf("Token %d", i))
	}

	resp := client.Get("/api/v1/webhooks/incoming")
	client.AssertStatus(resp, http.StatusOK)

	var result []interface{}
	client.ParseJSON(resp, &result)

	if len(result) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(result))
	}
}

func TestWebhooks_ListIncomingTokens_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/webhooks/incoming")
	client.AssertStatus(resp, http.StatusOK)

	var result []interface{}
	client.ParseJSON(resp, &result)

	if len(result) != 0 {
		t.Errorf("Expected 0 tokens, got %d", len(result))
	}
}

// ============================================================================
// DELETE /api/v1/webhooks/incoming/:id
// ============================================================================

func TestWebhooks_DeleteIncomingToken_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	token, _ := testFixtures.CreateIncomingWebhookToken(ctx, user.Organization.ID, "Test Token")

	resp := client.Delete(fmt.Sprintf("/api/v1/webhooks/incoming/%s", token.ID))
	client.AssertStatus(resp, http.StatusOK)
}

func TestWebhooks_DeleteIncomingToken_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Delete("/api/v1/webhooks/incoming/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// POST /api/v1/webhook/:token (Public incoming webhook)
// ============================================================================

func TestWebhooks_ReceiveWebhook_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	incomingToken, _ := testFixtures.CreateIncomingWebhookToken(ctx, user.Organization.ID, "Test Token")

	// Clear auth token - this is a public endpoint
	client.ClearAuthToken()

	// Send generic webhook payload
	reqBody := map[string]interface{}{
		"source":   "test-system",
		"message":  "Test alert from webhook",
		"priority": "P3",
	}

	resp := client.Post(fmt.Sprintf("/api/v1/webhook/%s", incomingToken.Token), reqBody)
	client.AssertStatus(resp, http.StatusCreated) // API returns 201 for created alerts
}

func TestWebhooks_ReceiveWebhook_InvalidToken(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]interface{}{
		"message": "Test alert",
	}

	resp := client.Post("/api/v1/webhook/invalid-token", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized) // API returns 401 for invalid tokens
}
