package domain

import "errors"

var (
	// General errors
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")

	// Alert errors
	ErrInvalidPriority = errors.New("invalid alert priority")
	ErrInvalidStatus   = errors.New("invalid alert status")

	// Schedule errors
	ErrInvalidRotationType = errors.New("invalid rotation type")
	ErrInvalidTimezone     = errors.New("invalid timezone")
	ErrOverlapOverride     = errors.New("override overlaps with existing override")

	// Escalation errors
	ErrInvalidEscalationTarget = errors.New("invalid escalation target type")
)
