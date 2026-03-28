package dto

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type SendNotificationRequest struct {
	ChannelID uuid.UUID  `json:"channel_id" binding:"required"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	AlertID   *uuid.UUID `json:"alert_id,omitempty"`
	Recipient string     `json:"recipient" binding:"required"`
	Subject   *string    `json:"subject,omitempty"`
	Message   string     `json:"message" binding:"required"`
}

type CreateNotificationChannelRequest struct {
	Name        string             `json:"name" binding:"required"`
	ChannelType domain.ChannelType `json:"channel_type" binding:"required"`
	IsEnabled   bool               `json:"is_enabled"`
	Config      json.RawMessage    `json:"config" binding:"required"`
}

type UpdateNotificationChannelRequest struct {
	Name        *string             `json:"name,omitempty"`
	ChannelType *domain.ChannelType `json:"channel_type,omitempty"`
	IsEnabled   *bool               `json:"is_enabled,omitempty"`
	Config      json.RawMessage     `json:"config,omitempty"`
}

type CreateUserNotificationPreferenceRequest struct {
	ChannelID    uuid.UUID `json:"channel_id" binding:"required"`
	IsEnabled    bool      `json:"is_enabled"`
	DNDEnabled   bool      `json:"dnd_enabled"`
	DNDStartTime *string   `json:"dnd_start_time,omitempty"`
	DNDEndTime   *string   `json:"dnd_end_time,omitempty"`
	MinPriority  *string   `json:"min_priority,omitempty"`
}

type UpdateUserNotificationPreferenceRequest struct {
	IsEnabled    *bool   `json:"is_enabled,omitempty"`
	DNDEnabled   *bool   `json:"dnd_enabled,omitempty"`
	DNDStartTime *string `json:"dnd_start_time,omitempty"`
	DNDEndTime   *string `json:"dnd_end_time,omitempty"`
	MinPriority  *string `json:"min_priority,omitempty"`
}
