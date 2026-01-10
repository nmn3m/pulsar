package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// Create godoc
// @Summary      Create a new team
// @Description  Create a new team in the organization
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.CreateTeamRequest true "Create team request"
// @Success      201 {object} domain.Team
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /teams [post]
func (h *TeamHandler) Create(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req service.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamService.CreateTeam(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// Get godoc
// @Summary      Get a team
// @Description  Get a team by ID with its members
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Success      200 {object} service.TeamWithMembers
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /teams/{id} [get]
func (h *TeamHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	team, err := h.teamService.GetTeamWithMembers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// Update godoc
// @Summary      Update a team
// @Description  Update a team by ID
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        request body service.UpdateTeamRequest true "Update team request"
// @Success      200 {object} domain.Team
// @Failure      400 {object} map[string]string
// @Router       /teams/{id} [patch]
func (h *TeamHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var req service.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamService.UpdateTeam(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// Delete godoc
// @Summary      Delete a team
// @Description  Delete a team by ID
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /teams/{id} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	if err := h.teamService.DeleteTeam(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "team deleted successfully"})
}

// List godoc
// @Summary      List teams
// @Description  List all teams in the organization
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Page size" default(20)
// @Success      200 {object} map[string][]domain.Team
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /teams [get]
func (h *TeamHandler) List(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	teams, err := h.teamService.ListTeams(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// AddMember godoc
// @Summary      Add a member to a team
// @Description  Add a user as a member to a team
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        request body service.AddTeamMemberRequest true "Add member request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/members [post]
func (h *TeamHandler) AddMember(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var req service.AddTeamMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.teamService.AddMember(c.Request.Context(), teamID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}

// RemoveMember godoc
// @Summary      Remove a member from a team
// @Description  Remove a user from a team
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        userId path string true "User ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.teamService.RemoveMember(c.Request.Context(), teamID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed successfully"})
}

// UpdateMemberRole godoc
// @Summary      Update a team member's role
// @Description  Update the role of a team member
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        userId path string true "User ID" format(uuid)
// @Param        request body service.UpdateTeamMemberRoleRequest true "Update role request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/members/{userId} [patch]
func (h *TeamHandler) UpdateMemberRole(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req service.UpdateTeamMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.teamService.UpdateMemberRole(c.Request.Context(), teamID, userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member role updated successfully"})
}

// ListMembers godoc
// @Summary      List team members
// @Description  List all members of a team
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Success      200 {object} map[string][]domain.UserWithTeamRole
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /teams/{id}/members [get]
func (h *TeamHandler) ListMembers(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	members, err := h.teamService.ListMembers(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// InviteMember godoc
// @Summary      Invite a member to a team
// @Description  Invite a user by email. If user exists, adds them directly. If not, sends an invitation email.
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        request body service.InviteMemberRequest true "Invite member request"
// @Success      200 {object} service.InvitationResponse
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/invite [post]
func (h *TeamHandler) InviteMember(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var req service.InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.teamService.AddMemberOrInvite(c.Request.Context(), teamID, orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListInvitations godoc
// @Summary      List team invitations
// @Description  List all pending invitations for a team
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Success      200 {object} map[string][]domain.TeamInvitation
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/invitations [get]
func (h *TeamHandler) ListInvitations(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	invitations, err := h.teamService.ListTeamInvitations(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invitations": invitations})
}

// CancelInvitation godoc
// @Summary      Cancel a team invitation
// @Description  Cancel a pending invitation
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        invitationId path string true "Invitation ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/invitations/{invitationId} [delete]
func (h *TeamHandler) CancelInvitation(c *gin.Context) {
	invitationIDStr := c.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invitation ID"})
		return
	}

	if err := h.teamService.CancelInvitation(c.Request.Context(), invitationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation cancelled"})
}

// ResendInvitation godoc
// @Summary      Resend a team invitation
// @Description  Resend the invitation email
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Team ID" format(uuid)
// @Param        invitationId path string true "Invitation ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /teams/{id}/invitations/{invitationId}/resend [post]
func (h *TeamHandler) ResendInvitation(c *gin.Context) {
	invitationIDStr := c.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invitation ID"})
		return
	}

	if err := h.teamService.ResendInvitation(c.Request.Context(), invitationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation resent"})
}
