package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ChannelType represents the type of notification channel
type ChannelType string

const (
	ChannelTypeEmail    ChannelType = "email"
	ChannelTypeSlack    ChannelType = "slack"
	ChannelTypeTeams    ChannelType = "teams"
	ChannelTypeWebhook  ChannelType = "webhook"
	ChannelTypePush     ChannelType = "push"      // Firebase Cloud Messaging
	ChannelTypePushAPNS ChannelType = "push_apns" // Apple Push Notification Service (future)
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

// NotificationChannel represents a notification delivery channel
type NotificationChannel struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	OrganizationID uuid.UUID       `json:"organization_id" db:"organization_id"`
	Name           string          `json:"name" db:"name"`
	ChannelType    ChannelType     `json:"channel_type" db:"channel_type"`
	IsEnabled      bool            `json:"is_enabled" db:"is_enabled"`
	Config         json.RawMessage `json:"config" db:"config" swaggertype:"object"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// UserNotificationPreference represents a user's notification preferences for a specific channel
type UserNotificationPreference struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	ChannelID    uuid.UUID  `json:"channel_id" db:"channel_id"`
	IsEnabled    bool       `json:"is_enabled" db:"is_enabled"`
	DNDEnabled   bool       `json:"dnd_enabled" db:"dnd_enabled"`
	DNDStartTime *time.Time `json:"dnd_start_time,omitempty" db:"dnd_start_time"`
	DNDEndTime   *time.Time `json:"dnd_end_time,omitempty" db:"dnd_end_time"`
	MinPriority  *string    `json:"min_priority,omitempty" db:"min_priority"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// NotificationLog represents a record of a sent (or attempted) notification
type NotificationLog struct {
	ID             uuid.UUID          `json:"id" db:"id"`
	OrganizationID uuid.UUID          `json:"organization_id" db:"organization_id"`
	ChannelID      uuid.UUID          `json:"channel_id" db:"channel_id"`
	UserID         *uuid.UUID         `json:"user_id,omitempty" db:"user_id"`
	AlertID        *uuid.UUID         `json:"alert_id,omitempty" db:"alert_id"`
	Recipient      string             `json:"recipient" db:"recipient"`
	Subject        *string            `json:"subject,omitempty" db:"subject"`
	Message        string             `json:"message" db:"message"`
	Status         NotificationStatus `json:"status" db:"status"`
	ErrorMessage   *string            `json:"error_message,omitempty" db:"error_message"`
	SentAt         *time.Time         `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt      time.Time          `json:"created_at" db:"created_at"`
}

// NotificationProvider is the interface that all notification providers must implement
type NotificationProvider interface {
	// Send sends a notification to the specified recipient
	Send(recipient string, subject string, message string) error

	// ValidateConfig validates the provider configuration
	ValidateConfig(config json.RawMessage) error
}

// CreateNotificationChannelRequest represents a request to create a notification channel
type CreateNotificationChannelRequest struct {
	Name        string          `json:"name" binding:"required"`
	ChannelType ChannelType     `json:"channel_type" binding:"required"`
	IsEnabled   bool            `json:"is_enabled"`
	Config      json.RawMessage `json:"config" binding:"required" swaggertype:"object"`
}

// UpdateNotificationChannelRequest represents a request to update a notification channel
type UpdateNotificationChannelRequest struct {
	Name        *string         `json:"name,omitempty"`
	ChannelType *ChannelType    `json:"channel_type,omitempty"`
	IsEnabled   *bool           `json:"is_enabled,omitempty"`
	Config      json.RawMessage `json:"config,omitempty" swaggertype:"object"`
}

// CreateUserNotificationPreferenceRequest represents a request to create user notification preferences
type CreateUserNotificationPreferenceRequest struct {
	ChannelID    uuid.UUID `json:"channel_id" binding:"required"`
	IsEnabled    bool      `json:"is_enabled"`
	DNDEnabled   bool      `json:"dnd_enabled"`
	DNDStartTime *string   `json:"dnd_start_time,omitempty"`
	DNDEndTime   *string   `json:"dnd_end_time,omitempty"`
	MinPriority  *string   `json:"min_priority,omitempty"`
}

// UpdateUserNotificationPreferenceRequest represents a request to update user notification preferences
type UpdateUserNotificationPreferenceRequest struct {
	IsEnabled    *bool   `json:"is_enabled,omitempty"`
	DNDEnabled   *bool   `json:"dnd_enabled,omitempty"`
	DNDStartTime *string `json:"dnd_start_time,omitempty"`
	DNDEndTime   *string `json:"dnd_end_time,omitempty"`
	MinPriority  *string `json:"min_priority,omitempty"`
}

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	ChannelID uuid.UUID  `json:"channel_id" binding:"required"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	AlertID   *uuid.UUID `json:"alert_id,omitempty"`
	Recipient string     `json:"recipient" binding:"required"`
	Subject   *string    `json:"subject,omitempty"`
	Message   string     `json:"message" binding:"required"`
}
