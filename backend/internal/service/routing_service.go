package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type RoutingService struct {
	routingRepo repository.RoutingRuleRepository
}

func NewRoutingService(routingRepo repository.RoutingRuleRepository) *RoutingService {
	return &RoutingService{
		routingRepo: routingRepo,
	}
}

// CreateRule creates a new routing rule
func (s *RoutingService) CreateRule(ctx context.Context, orgID uuid.UUID, req *domain.CreateRoutingRuleRequest) (*domain.AlertRoutingRule, error) {
	// Validate conditions JSON
	var conditions domain.RoutingConditions
	if err := json.Unmarshal(req.Conditions, &conditions); err != nil {
		return nil, fmt.Errorf("invalid conditions format: %w", err)
	}

	// Validate actions JSON
	var actions domain.RoutingActions
	if err := json.Unmarshal(req.Actions, &actions); err != nil {
		return nil, fmt.Errorf("invalid actions format: %w", err)
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	rule := &domain.AlertRoutingRule{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		Priority:       req.Priority,
		Conditions:     req.Conditions,
		Actions:        req.Actions,
		Enabled:        enabled,
	}

	if err := s.routingRepo.Create(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// GetRule retrieves a routing rule by ID
func (s *RoutingService) GetRule(ctx context.Context, id uuid.UUID) (*domain.AlertRoutingRule, error) {
	return s.routingRepo.GetByID(ctx, id)
}

// UpdateRule updates an existing routing rule
func (s *RoutingService) UpdateRule(ctx context.Context, id uuid.UUID, req *domain.UpdateRoutingRuleRequest) (*domain.AlertRoutingRule, error) {
	rule, err := s.routingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Description != nil {
		rule.Description = req.Description
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.Conditions != nil {
		// Validate conditions
		var conditions domain.RoutingConditions
		if err := json.Unmarshal(req.Conditions, &conditions); err != nil {
			return nil, fmt.Errorf("invalid conditions format: %w", err)
		}
		rule.Conditions = req.Conditions
	}
	if req.Actions != nil {
		// Validate actions
		var actions domain.RoutingActions
		if err := json.Unmarshal(req.Actions, &actions); err != nil {
			return nil, fmt.Errorf("invalid actions format: %w", err)
		}
		rule.Actions = req.Actions
	}
	if req.Enabled != nil {
		rule.Enabled = *req.Enabled
	}

	if err := s.routingRepo.Update(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// DeleteRule deletes a routing rule
func (s *RoutingService) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return s.routingRepo.Delete(ctx, id)
}

// ListRules lists routing rules for an organization
func (s *RoutingService) ListRules(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.AlertRoutingRule, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize
	return s.routingRepo.List(ctx, orgID, pageSize, offset)
}

// ReorderRules reorders routing rules by priority
func (s *RoutingService) ReorderRules(ctx context.Context, orgID uuid.UUID, req *domain.ReorderRoutingRulesRequest) error {
	return s.routingRepo.Reorder(ctx, orgID, req.RuleIDs)
}

// ApplyRouting evaluates routing rules and returns the actions to apply to an alert
func (s *RoutingService) ApplyRouting(ctx context.Context, orgID uuid.UUID, alert *domain.Alert) (*domain.RoutingActions, error) {
	rules, err := s.routingRepo.ListEnabled(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list routing rules: %w", err)
	}

	// Rules are already sorted by priority (ascending)
	for _, rule := range rules {
		conditions, err := rule.ParseConditions()
		if err != nil {
			// Skip rule with invalid conditions
			continue
		}

		if s.evaluateConditions(alert, conditions) {
			actions, err := rule.ParseActions()
			if err != nil {
				// Skip rule with invalid actions
				continue
			}
			return actions, nil
		}
	}

	// No matching rule found
	return nil, nil
}

// evaluateConditions evaluates if the alert matches the routing conditions
func (s *RoutingService) evaluateConditions(alert *domain.Alert, conditions *domain.RoutingConditions) bool {
	if len(conditions.Conditions) == 0 {
		return true // No conditions means match all
	}

	matchAll := conditions.Match == "all"
	matchedCount := 0

	for _, condition := range conditions.Conditions {
		matched := s.evaluateCondition(alert, &condition)

		if matched {
			matchedCount++
			if !matchAll {
				// For "any" match, return true on first match
				return true
			}
		} else if matchAll {
			// For "all" match, return false on first non-match
			return false
		}
	}

	// For "all" match, all conditions must be matched
	return matchAll && matchedCount == len(conditions.Conditions)
}

// evaluateCondition evaluates a single condition against an alert
func (s *RoutingService) evaluateCondition(alert *domain.Alert, condition *domain.RoutingCondition) bool {
	var fieldValue string

	switch condition.Field {
	case "source":
		fieldValue = alert.Source
	case "priority":
		fieldValue = string(alert.Priority)
	case "message":
		fieldValue = alert.Message
	case "tags":
		// For tags, we check if any tag matches
		return s.evaluateTagsCondition(alert.Tags, condition)
	default:
		// Check custom fields
		if val, ok := alert.CustomFields[condition.Field]; ok {
			fieldValue = fmt.Sprintf("%v", val)
		} else {
			return false
		}
	}

	return s.evaluateOperator(fieldValue, condition.Operator, condition.Value)
}

// evaluateTagsCondition evaluates conditions on tags
func (s *RoutingService) evaluateTagsCondition(tags []string, condition *domain.RoutingCondition) bool {
	switch condition.Operator {
	case "contains":
		for _, tag := range tags {
			if tag == condition.Value {
				return true
			}
		}
		return false
	case "not_contains":
		for _, tag := range tags {
			if tag == condition.Value {
				return false
			}
		}
		return true
	default:
		// For other operators, join tags and compare
		tagsStr := strings.Join(tags, ",")
		return s.evaluateOperator(tagsStr, condition.Operator, condition.Value)
	}
}

// evaluateOperator evaluates an operator against field and value
func (s *RoutingService) evaluateOperator(fieldValue, operator, conditionValue string) bool {
	switch operator {
	case "equals":
		return fieldValue == conditionValue
	case "not_equals":
		return fieldValue != conditionValue
	case "contains":
		return strings.Contains(fieldValue, conditionValue)
	case "not_contains":
		return !strings.Contains(fieldValue, conditionValue)
	case "regex":
		matched, err := regexp.MatchString(conditionValue, fieldValue)
		return err == nil && matched
	case "gte":
		fieldNum, err1 := strconv.Atoi(fieldValue)
		condNum, err2 := strconv.Atoi(conditionValue)
		if err1 == nil && err2 == nil {
			return fieldNum >= condNum
		}
		return fieldValue >= conditionValue
	case "lte":
		fieldNum, err1 := strconv.Atoi(fieldValue)
		condNum, err2 := strconv.Atoi(conditionValue)
		if err1 == nil && err2 == nil {
			return fieldNum <= condNum
		}
		return fieldValue <= conditionValue
	case "starts_with":
		return strings.HasPrefix(fieldValue, conditionValue)
	case "ends_with":
		return strings.HasSuffix(fieldValue, conditionValue)
	default:
		return false
	}
}
