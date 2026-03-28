package dto

import (
	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	TeamID      *uuid.UUID `json:"team_id"`
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description"`
	Timezone    string     `json:"timezone"`
}

type UpdateScheduleRequest struct {
	TeamID      *uuid.UUID `json:"team_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Timezone    *string    `json:"timezone"`
}

type CreateRotationRequest struct {
	Name           string  `json:"name" binding:"required"`
	RotationType   string  `json:"rotation_type" binding:"required"`
	RotationLength int     `json:"rotation_length" binding:"required"`
	StartDate      string  `json:"start_date" binding:"required"`
	StartTime      string  `json:"start_time"`
	EndTime        *string `json:"end_time"`
	HandoffDay     *int    `json:"handoff_day"`
	HandoffTime    string  `json:"handoff_time"`
}

type UpdateRotationRequest struct {
	Name           *string `json:"name"`
	RotationType   *string `json:"rotation_type"`
	RotationLength *int    `json:"rotation_length"`
	StartDate      *string `json:"start_date"`
	StartTime      *string `json:"start_time"`
	EndTime        *string `json:"end_time"`
	HandoffDay     *int    `json:"handoff_day"`
	HandoffTime    *string `json:"handoff_time"`
}

type AddParticipantRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Position int       `json:"position" binding:"min=0"`
}

type ReorderParticipantsRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
}

type CreateOverrideRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Note      *string   `json:"note"`
}

type UpdateOverrideRequest struct {
	UserID    *uuid.UUID `json:"user_id"`
	StartTime *string    `json:"start_time"`
	EndTime   *string    `json:"end_time"`
	Note      *string    `json:"note"`
}
