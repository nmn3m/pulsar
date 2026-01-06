package integration

import (
	"context"
	"net/http"
	"testing"
)

// ============================================================================
// POST /api/v1/auth/register
// ============================================================================

func TestAuth_Register_Success(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]string{
		"email":             "newuser@example.com",
		"username":          "newuser",
		"password":          "SecurePassword123!",
		"full_name":         "New User",
		"organization_name": "New Organization",
	}

	resp := client.Post("/api/v1/auth/register", reqBody)
	client.AssertStatus(resp, http.StatusCreated)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["access_token"] == nil {
		t.Error("Expected access_token in response")
	}
	if result["refresh_token"] == nil {
		t.Error("Expected refresh_token in response")
	}
	if result["user"] == nil {
		t.Error("Expected user in response")
	}
	if result["organization"] == nil {
		t.Error("Expected organization in response")
	}

	user := result["user"].(map[string]interface{})
	if user["email"] != "newuser@example.com" {
		t.Errorf("Expected email 'newuser@example.com', got %v", user["email"])
	}
}

func TestAuth_Register_DuplicateEmail(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	// Create existing user
	_, err := testFixtures.CreateUser(ctx, "existing@example.com", "existinguser", "Existing Org")
	if err != nil {
		t.Fatalf("Failed to create existing user: %v", err)
	}

	// Try to register with same email
	reqBody := map[string]string{
		"email":             "existing@example.com",
		"username":          "newuser",
		"password":          "SecurePassword123!",
		"full_name":         "New User",
		"organization_name": "New Org",
	}

	resp := client.Post("/api/v1/auth/register", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

func TestAuth_Register_MissingRequiredFields(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	// Missing email
	reqBody := map[string]string{
		"username":          "newuser",
		"password":          "SecurePassword123!",
		"organization_name": "New Org",
	}

	resp := client.Post("/api/v1/auth/register", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

func TestAuth_Register_InvalidEmail(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]string{
		"email":             "invalid-email",
		"username":          "newuser",
		"password":          "SecurePassword123!",
		"full_name":         "New User",
		"organization_name": "New Org",
	}

	resp := client.Post("/api/v1/auth/register", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// POST /api/v1/auth/login
// ============================================================================

func TestAuth_Login_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	// Create user first
	_, err := testFixtures.CreateUser(ctx, "login@example.com", "loginuser", "Login Org")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	reqBody := map[string]string{
		"email":    "login@example.com",
		"password": "TestPassword123!",
	}

	resp := client.Post("/api/v1/auth/login", reqBody)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["access_token"] == nil {
		t.Error("Expected access_token in response")
	}
	if result["refresh_token"] == nil {
		t.Error("Expected refresh_token in response")
	}
}

func TestAuth_Login_InvalidEmail(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "AnyPassword123!",
	}

	resp := client.Post("/api/v1/auth/login", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestAuth_Login_InvalidPassword(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	// Create user first
	_, err := testFixtures.CreateUser(ctx, "login@example.com", "loginuser", "Login Org")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	reqBody := map[string]string{
		"email":    "login@example.com",
		"password": "WrongPassword123!",
	}

	resp := client.Post("/api/v1/auth/login", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestAuth_Login_MissingCredentials(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	// Missing password
	reqBody := map[string]string{
		"email": "test@example.com",
	}

	resp := client.Post("/api/v1/auth/login", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// POST /api/v1/auth/refresh
// ============================================================================

func TestAuth_RefreshToken_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, err := testFixtures.CreateUniqueUser(ctx)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	reqBody := map[string]string{
		"refresh_token": user.RefreshToken,
	}

	resp := client.Post("/api/v1/auth/refresh", reqBody)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["access_token"] == nil {
		t.Error("Expected access_token in response")
	}
}

func TestAuth_RefreshToken_InvalidToken(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]string{
		"refresh_token": "invalid_token",
	}

	resp := client.Post("/api/v1/auth/refresh", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestAuth_RefreshToken_MissingToken(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]string{}

	resp := client.Post("/api/v1/auth/refresh", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// GET /api/v1/auth/me
// ============================================================================

func TestAuth_GetMe_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, err := testFixtures.CreateUniqueUser(ctx)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	client.SetAuthToken(user.AccessToken)
	resp := client.Get("/api/v1/auth/me")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["email"] != user.User.Email {
		t.Errorf("Expected email %s, got %v", user.User.Email, result["email"])
	}
	if result["id"] != user.User.ID.String() {
		t.Errorf("Expected user ID %s, got %v", user.User.ID, result["id"])
	}
}

func TestAuth_GetMe_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	// No auth token set
	resp := client.Get("/api/v1/auth/me")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestAuth_GetMe_InvalidToken(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	client.SetAuthToken("invalid_token")
	resp := client.Get("/api/v1/auth/me")
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

// ============================================================================
// POST /api/v1/auth/logout
// ============================================================================

func TestAuth_Logout_Success(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	resp := client.Post("/api/v1/auth/logout", nil)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["message"] == nil {
		t.Error("Expected message in response")
	}
}
