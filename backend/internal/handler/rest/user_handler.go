package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
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
