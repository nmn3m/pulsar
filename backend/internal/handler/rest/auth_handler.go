package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type AuthHandler struct {
	authService              *service.AuthService
	emailVerificationService *service.EmailVerificationService
}

func NewAuthHandler(authService *service.AuthService, emailVerificationService *service.EmailVerificationService) *AuthHandler {
	return &AuthHandler{
		authService:              authService,
		emailVerificationService: emailVerificationService,
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

// VerifyEmail godoc
// @Summary      Verify email with OTP
// @Description  Verify user email address using OTP code sent via email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body domain.VerifyEmailRequest true "Verify email request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	if h.emailVerificationService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "email verification not configured"})
		return
	}

	var req domain.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailVerificationService.VerifyOTP(c.Request.Context(), req.Email, req.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

// ResendOTP godoc
// @Summary      Resend OTP verification code
// @Description  Resend OTP verification code to user's email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body domain.ResendOTPRequest true "Resend OTP request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	if h.emailVerificationService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "email verification not configured"})
		return
	}

	var req domain.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailVerificationService.ResendOTP(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification code sent"})
}
