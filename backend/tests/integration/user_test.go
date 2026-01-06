package integration

import (
	"context"
	"net/http"
	"testing"
)

// ============================================================================
// GET /api/v1/users
// ============================================================================

func TestUsers_List_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/users")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	users := result["users"].([]interface{})
	// Should have at least the user who made the request
	if len(users) < 1 {
		t.Errorf("Expected at least 1 user, got %d", len(users))
	}
}

func TestUsers_List_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Get("/api/v1/users")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}
