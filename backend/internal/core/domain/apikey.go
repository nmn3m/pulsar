package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Name           string
	KeyPrefix      string // First 8 chars for identification
	KeyHash        string // SHA-256 hash of the full key
	Scopes         []string
	LastUsedAt     *time.Time
	ExpiresAt      *time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// APIKeyScope defines available API key scopes
type APIKeyScope string

const (
	ScopeAlertsRead         APIKeyScope = "alerts:read"
	ScopeAlertsWrite        APIKeyScope = "alerts:write"
	ScopeIncidentsRead      APIKeyScope = "incidents:read"
	ScopeIncidentsWrite     APIKeyScope = "incidents:write"
	ScopeTeamsRead          APIKeyScope = "teams:read"
	ScopeTeamsWrite         APIKeyScope = "teams:write"
	ScopeSchedulesRead      APIKeyScope = "schedules:read"
	ScopeSchedulesWrite     APIKeyScope = "schedules:write"
	ScopeWebhooksRead       APIKeyScope = "webhooks:read"
	ScopeWebhooksWrite      APIKeyScope = "webhooks:write"
	ScopeNotificationsRead  APIKeyScope = "notifications:read"
	ScopeNotificationsWrite APIKeyScope = "notifications:write"
	ScopeUsersRead          APIKeyScope = "users:read"
	ScopeAll                APIKeyScope = "*"
)

// ValidScopes returns all valid API key scopes
func ValidScopes() []APIKeyScope {
	return []APIKeyScope{
		ScopeAlertsRead, ScopeAlertsWrite,
		ScopeIncidentsRead, ScopeIncidentsWrite,
		ScopeTeamsRead, ScopeTeamsWrite,
		ScopeSchedulesRead, ScopeSchedulesWrite,
		ScopeWebhooksRead, ScopeWebhooksWrite,
		ScopeNotificationsRead, ScopeNotificationsWrite,
		ScopeUsersRead,
		ScopeAll,
	}
}

// IsValidScope checks if a scope string is valid
func IsValidScope(scope string) bool {
	for _, s := range ValidScopes() {
		if string(s) == scope {
			return true
		}
	}
	return false
}

// HasScope checks if the API key has a specific scope
func (k *APIKey) HasScope(scope APIKeyScope) bool {
	for _, s := range k.Scopes {
		if s == string(ScopeAll) || s == string(scope) {
			return true
		}
		// Check for wildcard read/write permissions
		if scope == ScopeAlertsRead && s == string(ScopeAlertsWrite) {
			return true
		}
		if scope == ScopeIncidentsRead && s == string(ScopeIncidentsWrite) {
			return true
		}
		if scope == ScopeTeamsRead && s == string(ScopeTeamsWrite) {
			return true
		}
		if scope == ScopeSchedulesRead && s == string(ScopeSchedulesWrite) {
			return true
		}
		if scope == ScopeWebhooksRead && s == string(ScopeWebhooksWrite) {
			return true
		}
		if scope == ScopeNotificationsRead && s == string(ScopeNotificationsWrite) {
			return true
		}
	}
	return false
}

// GenerateAPIKey generates a new API key and returns the raw key and its hash
// The raw key format: pls_<32 random hex chars>
func GenerateAPIKey() (rawKey, keyPrefix, keyHash string, err error) {
	// Generate 32 random bytes
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", "", err
	}

	// Create the key with prefix
	rawKey = "pls_" + hex.EncodeToString(randomBytes)
	keyPrefix = rawKey[:12] // "pls_" + first 8 hex chars

	// Hash the key for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash = hex.EncodeToString(hash[:])

	return rawKey, keyPrefix, keyHash, nil
}

// HashAPIKey hashes an API key for comparison
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
