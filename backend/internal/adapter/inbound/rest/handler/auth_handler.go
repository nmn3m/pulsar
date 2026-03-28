package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/nmn3m/pulsar/backend/internal/adapter/inbound/rest/middleware"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
	"github.com/nmn3m/pulsar/backend/internal/core/port/inbound"
	"github.com/nmn3m/pulsar/backend/internal/core/port/outbound"
)

// VerifyEmailRequest is the request body for email verification
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// ResendOTPRequest is the request body for resending OTP
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthHandler struct {
	authService              inbound.AuthService
	emailVerificationService inbound.EmailVerificationService
	blacklist                outbound.TokenRevoker
}

func NewAuthHandler(authService inbound.AuthService, emailVerificationService inbound.EmailVerificationService, blacklist outbound.TokenRevoker) *AuthHandler {
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
// @Param        request body dto.RegisterRequest true "Registration request"
// @Success      201 {object} dto.AuthResponse
// @Failure      400 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
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
// @Param        request body dto.LoginRequest true "Login request"
// @Success      200 {object} dto.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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
// @Success      200 {object} dto.AuthResponse
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
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		tokenString := parts[1]
		// Parse without validation to extract expiry for blacklist TTL
		token, _ := jwt.ParseWithClaims(tokenString, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		expiry := time.Now().Add(24 * time.Hour) // fallback expiry
		if token != nil {
			if claims, ok := token.Claims.(*middleware.Claims); ok && claims.ExpiresAt != nil {
				expiry = claims.ExpiresAt.Time
			}
		}
		h.blacklist.Revoke(tokenString, expiry)
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// VerifyEmail godoc
// @Summary      Verify email with OTP
// @Description  Verify user email address using OTP code sent via email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body VerifyEmailRequest true "Verify email request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	if h.emailVerificationService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "email verification not configured"})
		return
	}

	var req VerifyEmailRequest
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
// @Param        request body ResendOTPRequest true "Resend OTP request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	if h.emailVerificationService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "email verification not configured"})
		return
	}

	var req ResendOTPRequest
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
