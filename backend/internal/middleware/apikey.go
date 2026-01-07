package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
)

// APIKeyValidator interface for validating API keys
type APIKeyValidator interface {
	ValidateAPIKey(ctx context.Context, rawKey string) (*domain.APIKey, error)
}

// APIKeyMiddleware handles API key authentication
type APIKeyMiddleware struct {
	validator APIKeyValidator
}

// NewAPIKeyMiddleware creates a new API key middleware
func NewAPIKeyMiddleware(validator APIKeyValidator) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		validator: validator,
	}
}

// RequireAPIKey middleware that requires a valid API key
func (m *APIKeyMiddleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := extractAPIKey(c)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			c.Abort()
			return
		}

		key, err := m.validator.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		// Set context values
		c.Set("user_id", key.UserID)
		c.Set("organization_id", key.OrganizationID)
		c.Set("api_key", key)
		c.Set("auth_type", "api_key")

		c.Next()
	}
}

// RequireAPIKeyWithScope middleware that requires a valid API key with specific scope
func (m *APIKeyMiddleware) RequireAPIKeyWithScope(scope domain.APIKeyScope) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := extractAPIKey(c)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			c.Abort()
			return
		}

		key, err := m.validator.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		// Check scope
		if !key.HasScope(scope) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "required_scope": string(scope)})
			c.Abort()
			return
		}

		// Set context values
		c.Set("user_id", key.UserID)
		c.Set("organization_id", key.OrganizationID)
		c.Set("api_key", key)
		c.Set("auth_type", "api_key")

		c.Next()
	}
}

// extractAPIKey extracts the API key from the request
// Supports: X-API-Key header, Authorization: ApiKey <key>, and query parameter ?api_key=<key>
func extractAPIKey(c *gin.Context) string {
	// Try X-API-Key header first
	if key := c.GetHeader("X-API-Key"); key != "" {
		return key
	}

	// Try Authorization header with ApiKey scheme
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "apikey" {
			return parts[1]
		}
	}

	// Try query parameter (not recommended for security, but supported)
	if key := c.Query("api_key"); key != "" {
		return key
	}

	return ""
}

// GetAPIKey gets the API key from context if authenticated via API key
func GetAPIKey(c *gin.Context) (*domain.APIKey, bool) {
	key, exists := c.Get("api_key")
	if !exists {
		return nil, false
	}
	apiKey, ok := key.(*domain.APIKey)
	return apiKey, ok
}

// IsAPIKeyAuth checks if the request was authenticated via API key
func IsAPIKeyAuth(c *gin.Context) bool {
	authType, exists := c.Get("auth_type")
	if !exists {
		return false
	}
	return authType == "api_key"
}

// CombinedAuthMiddleware creates a middleware that accepts both JWT and API key authentication
type CombinedAuthMiddleware struct {
	jwtAuth    *AuthMiddleware
	apiKeyAuth *APIKeyMiddleware
}

// NewCombinedAuthMiddleware creates a new combined auth middleware
func NewCombinedAuthMiddleware(jwtAuth *AuthMiddleware, apiKeyAuth *APIKeyMiddleware) *CombinedAuthMiddleware {
	return &CombinedAuthMiddleware{
		jwtAuth:    jwtAuth,
		apiKeyAuth: apiKeyAuth,
	}
}

// RequireAuth middleware that accepts either JWT or API key authentication
func (m *CombinedAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key first
		apiKey := extractAPIKey(c)
		if apiKey != "" {
			key, err := m.apiKeyAuth.validator.ValidateAPIKey(c.Request.Context(), apiKey)
			if err == nil {
				// Valid API key
				c.Set("user_id", key.UserID)
				c.Set("organization_id", key.OrganizationID)
				c.Set("api_key", key)
				c.Set("auth_type", "api_key")
				c.Next()
				return
			}
		}

		// Fall back to JWT authentication
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"})
			c.Abort()
			return
		}

		// Use JWT auth middleware logic
		m.jwtAuth.RequireAuth()(c)
	}
}

// RequireAuthWithScope middleware that accepts either JWT or API key with specific scope
func (m *CombinedAuthMiddleware) RequireAuthWithScope(scope domain.APIKeyScope) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key first
		apiKey := extractAPIKey(c)
		if apiKey != "" {
			key, err := m.apiKeyAuth.validator.ValidateAPIKey(c.Request.Context(), apiKey)
			if err == nil {
				// Check scope for API key
				if !key.HasScope(scope) {
					c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "required_scope": string(scope)})
					c.Abort()
					return
				}
				// Valid API key with scope
				c.Set("user_id", key.UserID)
				c.Set("organization_id", key.OrganizationID)
				c.Set("api_key", key)
				c.Set("auth_type", "api_key")
				c.Next()
				return
			}
		}

		// Fall back to JWT authentication (JWT users have full access)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"})
			c.Abort()
			return
		}

		// Use JWT auth middleware logic
		m.jwtAuth.RequireAuth()(c)
	}
}

// Helper to check if user has required scope (works for both JWT and API key)
func CheckScope(c *gin.Context, scope domain.APIKeyScope) bool {
	// JWT authenticated users have full access
	if !IsAPIKeyAuth(c) {
		return true
	}

	// For API key auth, check scope
	apiKey, ok := GetAPIKey(c)
	if !ok {
		return false
	}

	return apiKey.HasScope(scope)
}

// GetAuthUserID returns the user ID regardless of auth method
func GetAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	return GetUserID(c)
}

// GetAuthOrganizationID returns the organization ID regardless of auth method
func GetAuthOrganizationID(c *gin.Context) (uuid.UUID, bool) {
	return GetOrganizationID(c)
}
