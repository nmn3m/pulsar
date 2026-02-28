package domain

import (
	"time"

	"github.com/google/uuid"
)

// ChannelType represents the type of notification channel
type ChannelType string

const (
	ChannelTypeEmail   ChannelType = "email"
	ChannelTypeSlack   ChannelType = "slack"
	ChannelTypeTeams   ChannelType = "teams"
	ChannelTypeWebhook ChannelType = "webhook"
	ChannelTypeSMS     ChannelType = "sms"
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
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	ChannelType    ChannelType
	IsEnabled      bool
	Config         []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// UserNotificationPreference represents a user's notification preferences for a specific channel
type UserNotificationPreference struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	ChannelID    uuid.UUID
	IsEnabled    bool
	DNDEnabled   bool
	DNDStartTime *time.Time
	DNDEndTime   *time.Time
	MinPriority  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NotificationLog represents a record of a sent (or attempted) notification
type NotificationLog struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ChannelID      uuid.UUID
	UserID         *uuid.UUID
	AlertID        *uuid.UUID
	Recipient      string
	Subject        *string
	Message        string
	Status         NotificationStatus
	ErrorMessage   *string
	SentAt         *time.Time
	CreatedAt      time.Time
}
