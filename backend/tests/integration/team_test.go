package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/nmn3m/pulsar/backend/internal/service"
)

// ============================================================================
// POST /api/v1/teams
// ============================================================================

func TestTeams_Create_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name":        "Engineering Team",
		"description": "The engineering team",
	}

	resp := client.Post("/api/v1/teams", reqBody)
	client.AssertStatus(resp, http.StatusCreated)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["id"] == nil {
		t.Error("Expected id in response")
	}
	if result["name"] != "Engineering Team" {
		t.Errorf("Expected name 'Engineering Team', got %v", result["name"])
	}
}

func TestTeams_Create_Unauthorized(t *testing.T) {
	cleanDatabase(t)
	client := newTestClient(t)

	reqBody := map[string]interface{}{
		"name": "Test Team",
	}

	resp := client.Post("/api/v1/teams", reqBody)
	client.ExpectStatus(resp, http.StatusUnauthorized)
}

func TestTeams_Create_MissingName(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"description": "A team without a name",
	}

	resp := client.Post("/api/v1/teams", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// GET /api/v1/teams
// ============================================================================

func TestTeams_List_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	// Create some teams
	for i := 0; i < 3; i++ {
		testFixtures.CreateTeam(ctx, user.Organization.ID, fmt.Sprintf("Team %d", i))
	}

	resp := client.Get("/api/v1/teams")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	teams := result["teams"].([]interface{})
	if len(teams) != 3 {
		t.Errorf("Expected 3 teams, got %d", len(teams))
	}
}

func TestTeams_List_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/teams")
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Handle nil or empty teams array
	teamsRaw := result["teams"]
	if teamsRaw != nil {
		teams := teamsRaw.([]interface{})
		if len(teams) != 0 {
			t.Errorf("Expected 0 teams, got %d", len(teams))
		}
	}
}

// ============================================================================
// GET /api/v1/teams/:id
// ============================================================================

func TestTeams_Get_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	resp := client.Get(fmt.Sprintf("/api/v1/teams/%s", team.ID))
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// API returns team directly (not wrapped in "team" key)
	if result["id"] != team.ID.String() {
		t.Errorf("Expected team ID %s, got %v", team.ID, result["id"])
	}
}

func TestTeams_Get_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Get("/api/v1/teams/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusNotFound)
}

// ============================================================================
// PATCH /api/v1/teams/:id
// ============================================================================

func TestTeams_Update_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Original Name")

	reqBody := map[string]interface{}{
		"name":        "Updated Name",
		"description": "Updated description",
	}

	resp := client.Patch(fmt.Sprintf("/api/v1/teams/%s", team.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	if result["name"] != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %v", result["name"])
	}
}

func TestTeams_Update_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"name": "Updated Name",
	}

	resp := client.Patch("/api/v1/teams/00000000-0000-0000-0000-000000000000", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// DELETE /api/v1/teams/:id
// ============================================================================

func TestTeams_Delete_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	resp := client.Delete(fmt.Sprintf("/api/v1/teams/%s", team.ID))
	client.AssertStatus(resp, http.StatusOK)

	// Verify it's deleted
	resp = client.Get(fmt.Sprintf("/api/v1/teams/%s", team.ID))
	client.ExpectStatus(resp, http.StatusNotFound)
}

func TestTeams_Delete_NotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	resp := client.Delete("/api/v1/teams/00000000-0000-0000-0000-000000000000")
	client.ExpectStatus(resp, http.StatusBadRequest) // API returns 400 for not found errors
}

// ============================================================================
// POST /api/v1/teams/:id/members
// ============================================================================

func TestTeams_AddMember_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	reqBody := map[string]interface{}{
		"user_id": user.User.ID.String(),
		"role":    "member",
	}

	resp := client.Post(fmt.Sprintf("/api/v1/teams/%s/members", team.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)
}

func TestTeams_AddMember_TeamNotFound(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	reqBody := map[string]interface{}{
		"user_id": user.User.ID.String(),
		"role":    "member",
	}

	resp := client.Post("/api/v1/teams/00000000-0000-0000-0000-000000000000/members", reqBody)
	client.ExpectStatus(resp, http.StatusBadRequest)
}

// ============================================================================
// GET /api/v1/teams/:id/members
// ============================================================================

func TestTeams_ListMembers_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	// Add member
	testServer.TeamService.AddMember(ctx, team.ID, &service.AddTeamMemberRequest{
		UserID: &user.User.ID,
		Role:   "member",
	})

	resp := client.Get(fmt.Sprintf("/api/v1/teams/%s/members", team.ID))
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	members := result["members"].([]interface{})
	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}
}

func TestTeams_ListMembers_Empty(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	resp := client.Get(fmt.Sprintf("/api/v1/teams/%s/members", team.ID))
	client.AssertStatus(resp, http.StatusOK)

	var result map[string]interface{}
	client.ParseJSON(resp, &result)

	// Handle nil or empty members array
	membersRaw := result["members"]
	if membersRaw != nil {
		members := membersRaw.([]interface{})
		if len(members) != 0 {
			t.Errorf("Expected 0 members, got %d", len(members))
		}
	}
}

// ============================================================================
// DELETE /api/v1/teams/:id/members/:userId
// ============================================================================

func TestTeams_RemoveMember_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	// Add member first
	testServer.TeamService.AddMember(ctx, team.ID, &service.AddTeamMemberRequest{
		UserID: &user.User.ID,
		Role:   "member",
	})

	resp := client.Delete(fmt.Sprintf("/api/v1/teams/%s/members/%s", team.ID, user.User.ID))
	client.AssertStatus(resp, http.StatusOK)
}

// ============================================================================
// PATCH /api/v1/teams/:id/members/:userId
// ============================================================================

func TestTeams_UpdateMemberRole_Success(t *testing.T) {
	cleanDatabase(t)
	ctx := context.Background()
	client := newTestClient(t)

	user, _ := testFixtures.CreateUniqueUser(ctx)
	client.SetAuthToken(user.AccessToken)

	team, _ := testFixtures.CreateTeam(ctx, user.Organization.ID, "Test Team")

	// Add member first
	testServer.TeamService.AddMember(ctx, team.ID, &service.AddTeamMemberRequest{
		UserID: &user.User.ID,
		Role:   "member",
	})

	reqBody := map[string]interface{}{
		"role": "admin",
	}

	resp := client.Patch(fmt.Sprintf("/api/v1/teams/%s/members/%s", team.ID, user.User.ID), reqBody)
	client.AssertStatus(resp, http.StatusOK)
}
