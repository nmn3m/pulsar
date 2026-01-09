package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type EscalationHandler struct {
	escalationService *service.EscalationService
}

func NewEscalationHandler(escalationService *service.EscalationService) *EscalationHandler {
	return &EscalationHandler{
		escalationService: escalationService,
	}
}

// Policy handlers

// List godoc
// @Summary      List escalation policies
// @Description  Retrieves a paginated list of escalation policies for the authenticated user's organization
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query    int  false  "Page number"      default(1)
// @Param        page_size  query    int  false  "Page size"        default(20)
// @Success      200  {object}  map[string][]domain.EscalationPolicy  "List of escalation policies"
// @Failure      401  {object}  map[string]string                     "Unauthorized"
// @Failure      500  {object}  map[string]string                     "Internal server error"
// @Router       /escalation-policies [get]
func (h *EscalationHandler) List(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	policies, err := h.escalationService.ListPolicies(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"policies": policies})
}

// Create godoc
// @Summary      Create escalation policy
// @Description  Creates a new escalation policy for the authenticated user's organization
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      service.CreateEscalationPolicyRequest  true  "Escalation policy creation request"
// @Success      201      {object}  domain.EscalationPolicy                "Created escalation policy"
// @Failure      400      {object}  map[string]string                      "Bad request"
// @Failure      401      {object}  map[string]string                      "Unauthorized"
// @Router       /escalation-policies [post]
func (h *EscalationHandler) Create(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req service.CreateEscalationPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy, err := h.escalationService.CreatePolicy(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, policy)
}

// Get godoc
// @Summary      Get escalation policy
// @Description  Retrieves a specific escalation policy by ID with its rules and targets
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Escalation policy ID"  format(uuid)
// @Success      200  {object}  domain.EscalationPolicyWithRules  "Escalation policy with rules"
// @Failure      400  {object}  map[string]string                 "Invalid policy ID"
// @Failure      404  {object}  map[string]string                 "Policy not found"
// @Router       /escalation-policies/{id} [get]
func (h *EscalationHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	policy, err := h.escalationService.GetPolicyWithRules(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// Update godoc
// @Summary      Update escalation policy
// @Description  Updates an existing escalation policy by ID
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                                 true  "Escalation policy ID"  format(uuid)
// @Param        request  body      service.UpdateEscalationPolicyRequest  true  "Escalation policy update request"
// @Success      200      {object}  domain.EscalationPolicy                "Updated escalation policy"
// @Failure      400      {object}  map[string]string                      "Invalid request or policy ID"
// @Router       /escalation-policies/{id} [patch]
func (h *EscalationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	var req service.UpdateEscalationPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy, err := h.escalationService.UpdatePolicy(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// Delete godoc
// @Summary      Delete escalation policy
// @Description  Deletes an escalation policy by ID
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Escalation policy ID"  format(uuid)
// @Success      200  {object}  map[string]string  "Policy deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid policy ID"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /escalation-policies/{id} [delete]
func (h *EscalationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	if err := h.escalationService.DeletePolicy(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "policy deleted"})
}

// Rule handlers

// ListRules godoc
// @Summary      List escalation rules
// @Description  Retrieves all escalation rules for a specific policy
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Escalation policy ID"  format(uuid)
// @Success      200  {object}  map[string][]domain.EscalationRule  "List of escalation rules"
// @Failure      400  {object}  map[string]string                   "Invalid policy ID"
// @Failure      500  {object}  map[string]string                   "Internal server error"
// @Router       /escalation-policies/{id}/rules [get]
func (h *EscalationHandler) ListRules(c *gin.Context) {
	policyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	rules, err := h.escalationService.ListRules(c.Request.Context(), policyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

// CreateRule godoc
// @Summary      Create escalation rule
// @Description  Creates a new escalation rule for a specific policy
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                               true  "Escalation policy ID"  format(uuid)
// @Param        request  body      service.CreateEscalationRuleRequest  true  "Escalation rule creation request"
// @Success      201      {object}  domain.EscalationRule                "Created escalation rule"
// @Failure      400      {object}  map[string]string                    "Invalid request or policy ID"
// @Router       /escalation-policies/{id}/rules [post]
func (h *EscalationHandler) CreateRule(c *gin.Context) {
	policyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}

	var req service.CreateEscalationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.escalationService.CreateRule(c.Request.Context(), policyID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// GetRule godoc
// @Summary      Get escalation rule
// @Description  Retrieves a specific escalation rule by ID
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string  true  "Escalation policy ID"  format(uuid)
// @Param        ruleId  path      string  true  "Escalation rule ID"    format(uuid)
// @Success      200     {object}  domain.EscalationRule  "Escalation rule"
// @Failure      400     {object}  map[string]string      "Invalid rule ID"
// @Failure      404     {object}  map[string]string      "Rule not found"
// @Router       /escalation-policies/{id}/rules/{ruleId} [get]
func (h *EscalationHandler) GetRule(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	rule, err := h.escalationService.GetRule(c.Request.Context(), ruleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// UpdateRule godoc
// @Summary      Update escalation rule
// @Description  Updates an existing escalation rule by ID
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                               true  "Escalation policy ID"  format(uuid)
// @Param        ruleId   path      string                               true  "Escalation rule ID"    format(uuid)
// @Param        request  body      service.UpdateEscalationRuleRequest  true  "Escalation rule update request"
// @Success      200      {object}  domain.EscalationRule                "Updated escalation rule"
// @Failure      400      {object}  map[string]string                    "Invalid request or rule ID"
// @Router       /escalation-policies/{id}/rules/{ruleId} [patch]
func (h *EscalationHandler) UpdateRule(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var req service.UpdateEscalationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.escalationService.UpdateRule(c.Request.Context(), ruleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// DeleteRule godoc
// @Summary      Delete escalation rule
// @Description  Deletes an escalation rule by ID
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string  true  "Escalation policy ID"  format(uuid)
// @Param        ruleId  path      string  true  "Escalation rule ID"    format(uuid)
// @Success      200     {object}  map[string]string  "Rule deleted successfully"
// @Failure      400     {object}  map[string]string  "Invalid rule ID"
// @Failure      500     {object}  map[string]string  "Internal server error"
// @Router       /escalation-policies/{id}/rules/{ruleId} [delete]
func (h *EscalationHandler) DeleteRule(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	if err := h.escalationService.DeleteRule(c.Request.Context(), ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}

// Target handlers

// ListTargets godoc
// @Summary      List escalation targets
// @Description  Retrieves all escalation targets for a specific rule
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string  true  "Escalation policy ID"  format(uuid)
// @Param        ruleId  path      string  true  "Escalation rule ID"    format(uuid)
// @Success      200     {object}  map[string][]domain.EscalationTarget  "List of escalation targets"
// @Failure      400     {object}  map[string]string                     "Invalid rule ID"
// @Failure      500     {object}  map[string]string                     "Internal server error"
// @Router       /escalation-policies/{id}/rules/{ruleId}/targets [get]
func (h *EscalationHandler) ListTargets(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	targets, err := h.escalationService.ListTargets(c.Request.Context(), ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"targets": targets})
}

// AddTarget godoc
// @Summary      Add escalation target
// @Description  Adds a new target (user, team, or schedule) to an escalation rule
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                              true  "Escalation policy ID"  format(uuid)
// @Param        ruleId   path      string                              true  "Escalation rule ID"    format(uuid)
// @Param        request  body      service.AddEscalationTargetRequest  true  "Escalation target request"
// @Success      201      {object}  domain.EscalationTarget             "Created escalation target"
// @Failure      400      {object}  map[string]string                   "Invalid request or rule ID"
// @Router       /escalation-policies/{id}/rules/{ruleId}/targets [post]
func (h *EscalationHandler) AddTarget(c *gin.Context) {
	ruleID, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var req service.AddEscalationTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := h.escalationService.AddTarget(c.Request.Context(), ruleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// RemoveTarget godoc
// @Summary      Remove escalation target
// @Description  Removes a target from an escalation rule
// @Tags         Escalation Policies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      string  true  "Escalation policy ID"  format(uuid)
// @Param        ruleId    path      string  true  "Escalation rule ID"    format(uuid)
// @Param        targetId  path      string  true  "Escalation target ID"  format(uuid)
// @Success      200       {object}  map[string]string  "Target removed successfully"
// @Failure      400       {object}  map[string]string  "Invalid target ID"
// @Failure      500       {object}  map[string]string  "Internal server error"
// @Router       /escalation-policies/{id}/rules/{ruleId}/targets/{targetId} [delete]
func (h *EscalationHandler) RemoveTarget(c *gin.Context) {
	targetID, err := uuid.Parse(c.Param("targetId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target id"})
		return
	}

	if err := h.escalationService.RemoveTarget(c.Request.Context(), targetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "target removed"})
}
