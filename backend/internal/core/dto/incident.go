package dto

import (
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type CreateIncidentRequest struct {
	Title            string     `json:"title" binding:"required"`
	Description      *string    `json:"description"`
	Severity         string     `json:"severity" binding:"required"`
	Priority         string     `json:"priority" binding:"required"`
	AssignedToTeamID *uuid.UUID `json:"assigned_to_team_id"`
}

type UpdateIncidentRequest struct {
	Title            *string    `json:"title"`
	Description      *string    `json:"description"`
	Severity         *string    `json:"severity"`
	Status           *string    `json:"status"`
	Priority         *string    `json:"priority"`
	AssignedToTeamID *uuid.UUID `json:"assigned_to_team_id"`
}

type AddResponderRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role" binding:"required"`
}

type UpdateResponderRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type AddNoteRequest struct {
	Note string `json:"note" binding:"required"`
}

type LinkAlertRequest struct {
	AlertID uuid.UUID `json:"alert_id" binding:"required"`
}

type ListIncidentsRequest struct {
	Status           []string   `form:"status"`
	Severity         []string   `form:"severity"`
	AssignedToTeamID *uuid.UUID `form:"assigned_to_team_id"`
	Search           *string    `form:"search"`
	Page             int        `form:"page"`
	PageSize         int        `form:"page_size"`
}

type ListIncidentsResponse struct {
	Incidents []*domain.Incident `json:"incidents"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}
