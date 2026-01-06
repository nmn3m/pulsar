package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with an organization
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body service.RegisterRequest true "Registration request"
// @Success      201 {object} service.AuthResponse
// @Failure      400 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and return tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body service.LoginRequest true "Login request"
// @Success      200 {object} service.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body object{refresh_token=string} true "Refresh token request"
// @Success      200 {object} service.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get the currently authenticated user's information
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.User
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout the current user (client should discard tokens)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT implementation, logout is handled client-side
	// by removing the token. For enhanced security, you could implement
	// token blacklisting here.
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
