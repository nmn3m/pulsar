package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	Name      string                 `json:"name" db:"name"`
	Slug      string                 `json:"slug" db:"slug"`
	Plan      string                 `json:"plan" db:"plan"`
	Settings  map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// PlanType represents the organization's subscription plan
type PlanType string

const (
	PlanFree       PlanType = "free"
	PlanPro        PlanType = "pro"
	PlanEnterprise PlanType = "enterprise"
)

func (p PlanType) String() string {
	return string(p)
}

func (p PlanType) IsValid() bool {
	switch p {
	case PlanFree, PlanPro, PlanEnterprise:
		return true
	}
	return false
}
