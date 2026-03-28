package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nmn3m/pulsar/backend/internal/adapter/inbound/rest/middleware"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
	"github.com/nmn3m/pulsar/backend/internal/core/port/inbound"
)

type UserHandler struct {
	userService inbound.UserService
}

func NewUserHandler(userService inbound.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ListOrganizationUsers godoc
// @Summary      List organization users
// @Description  List all users in the organization
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string][]domain.User
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /users [get]
func (h *UserHandler) ListOrganizationUsers(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := h.userService.ListOrganizationUsers(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// UpdateProfile godoc
// @Summary      Update current user's profile
// @Description  Update the authenticated user's profile (full name, phone, timezone)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body dto.UpdateProfileRequest true "Profile update fields"
// @Success      200 {object} domain.User
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /users/me [patch]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
