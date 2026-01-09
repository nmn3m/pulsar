package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type RoutingHandler struct {
	routingService *service.RoutingService
}

func NewRoutingHandler(routingService *service.RoutingService) *RoutingHandler {
	return &RoutingHandler{
		routingService: routingService,
	}
}

// List godoc
// @Summary      List routing rules
// @Description  Retrieves a paginated list of alert routing rules for the authenticated user's organization
// @Tags         Routing Rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query    int  false  "Page number"      default(1)
// @Param        page_size  query    int  false  "Page size"        default(50)
// @Success      200  {object}  map[string][]domain.AlertRoutingRule  "List of routing rules"
// @Failure      401  {object}  map[string]string                     "Unauthorized"
// @Failure      500  {object}  map[string]string                     "Internal server error"
// @Router       /routing-rules [get]
func (h *RoutingHandler) List(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	rules, err := h.routingService.ListRules(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

// Create godoc
// @Summary      Create routing rule
// @Description  Creates a new alert routing rule for the authenticated user's organization
// @Tags         Routing Rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      domain.CreateRoutingRuleRequest  true  "Routing rule creation request"
// @Success      201      {object}  domain.AlertRoutingRule          "Created routing rule"
// @Failure      400      {object}  map[string]string                "Bad request"
// @Failure      401      {object}  map[string]string                "Unauthorized"
// @Router       /routing-rules [post]
func (h *RoutingHandler) Create(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateRoutingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.routingService.CreateRule(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// Get godoc
// @Summary      Get routing rule
// @Description  Retrieves a specific routing rule by ID
// @Tags         Routing Rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Routing rule ID"  format(uuid)
// @Success      200  {object}  domain.AlertRoutingRule  "Routing rule"
// @Failure      400  {object}  map[string]string        "Invalid rule ID"
// @Failure      404  {object}  map[string]string        "Rule not found"
// @Router       /routing-rules/{id} [get]
func (h *RoutingHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	rule, err := h.routingService.GetRule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// Update godoc
// @Summary      Update routing rule
// @Description  Updates an existing routing rule by ID
// @Tags         Routing Rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                           true  "Routing rule ID"  format(uuid)
// @Param        request  body      domain.UpdateRoutingRuleRequest  true  "Routing rule update request"
// @Success      200      {object}  domain.AlertRoutingRule          "Updated routing rule"
// @Failure      400      {object}  map[string]string                "Invalid request or rule ID"
// @Router       /routing-rules/{id} [patch]
func (h *RoutingHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var req domain.UpdateRoutingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule, err := h.routingService.UpdateRule(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// Delete godoc
// @Summary      Delete routing rule
// @Description  Deletes a routing rule by ID
// @Tags         Routing Rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Routing rule ID"  format(uuid)
// @Success      200  {object}  map[string]string  "Rule deleted successfully"
// @Failure      400  {object}  map[string]string  "Invalid rule ID"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /routing-rules/{id} [delete]
func (h *RoutingHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	if err := h.routingService.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted"})
}

// Reorder godoc
// @Summary      Reorder routing rules
// @Description  Reorders routing rules by setting their priorities based on the provided order
// @Tags         Routing Rules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      domain.ReorderRoutingRulesRequest  true  "Routing rules reorder request"
// @Success      200      {object}  map[string]string                  "Rules reordered successfully"
// @Failure      400      {object}  map[string]string                  "Bad request"
// @Failure      401      {object}  map[string]string                  "Unauthorized"
// @Failure      500      {object}  map[string]string                  "Internal server error"
// @Router       /routing-rules/reorder [put]
func (h *RoutingHandler) Reorder(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.ReorderRoutingRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.routingService.ReorderRules(c.Request.Context(), orgID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rules reordered"})
}
