package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type CreateAlertRequest struct {
	Source       string                 `json:"source" binding:"required"`
	SourceID     *string                `json:"source_id"`
	Priority     string                 `json:"priority" binding:"required"`
	Message      string                 `json:"message" binding:"required"`
	Description  *string                `json:"description"`
	Tags         []string               `json:"tags"`
	CustomFields map[string]interface{} `json:"custom_fields"`
	DedupKey     *string                `json:"dedup_key"` // Optional deduplication key
}

type UpdateAlertRequest struct {
	Priority     *string                `json:"priority"`
	Message      *string                `json:"message"`
	Description  *string                `json:"description"`
	Tags         []string               `json:"tags"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type AcknowledgeAlertRequest struct {
	UserID uuid.UUID
}

type CloseAlertRequest struct {
	UserID uuid.UUID
	Reason string `json:"reason"`
}

type SnoozeAlertRequest struct {
	Until time.Time `json:"until" binding:"required"`
}

type AssignAlertRequest struct {
	UserID *uuid.UUID `json:"user_id"`
	TeamID *uuid.UUID `json:"team_id"`
}

type ListAlertsRequest struct {
	Status         []string   `form:"status"`
	Priority       []string   `form:"priority"`
	AssignedToUser *uuid.UUID `form:"assigned_to_user"`
	AssignedToTeam *uuid.UUID `form:"assigned_to_team"`
	Source         *string    `form:"source"`
	Search         *string    `form:"search"`
	Page           int        `form:"page"`
	PageSize       int        `form:"page_size"`
}

type ListAlertsResponse struct {
	Alerts   []*domain.Alert `json:"alerts"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
