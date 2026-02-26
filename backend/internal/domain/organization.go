package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	Plan      string
	Settings  map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
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
