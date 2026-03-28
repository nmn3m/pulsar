package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/delivery/rest/middleware"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/usecase"
)

type IncidentHandler struct {
	incidentUsecase *usecase.IncidentUsecase
}

func NewIncidentHandler(incidentUsecase *usecase.IncidentUsecase) *IncidentHandler {
	return &IncidentHandler{
		incidentUsecase: incidentUsecase,
	}
}

// Create godoc
// @Summary      Create a new incident
// @Description  Creates a new incident with the provided details. The incident is created in the investigating status.
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body usecase.CreateIncidentRequest true "Incident creation request"
// @Success      201 {object} domain.Incident
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents [post]
func (h *IncidentHandler) Create(c *gin.Context) {
	var req usecase.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)
	userID, _ := middleware.GetUserID(c)

	incident, err := h.incidentUsecase.CreateIncident(c.Request.Context(), orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// Get retrieves a basic incident by ID (internal use - GetWithDetails is exposed via API)
func (h *IncidentHandler) Get(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	incident, err := h.incidentUsecase.GetIncident(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// GetWithDetails godoc
// @Summary      Get an incident with full details
// @Description  Retrieves an incident with all related data including responders, alerts, and timeline events
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Success      200 {object} domain.IncidentWithDetails
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id} [get]
func (h *IncidentHandler) GetWithDetails(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	incident, err := h.incidentUsecase.GetIncidentWithDetails(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// Update godoc
// @Summary      Update an incident
// @Description  Updates an existing incident with the provided fields. Only fields included in the request body will be updated.
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        request body usecase.UpdateIncidentRequest true "Incident update request"
// @Success      200 {object} domain.Incident
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id} [patch]
func (h *IncidentHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	var req usecase.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)
	userID, _ := middleware.GetUserID(c)

	incident, err := h.incidentUsecase.UpdateIncident(c.Request.Context(), id, orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// Delete godoc
// @Summary      Delete an incident
// @Description  Permanently deletes an incident and all associated data
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id} [delete]
func (h *IncidentHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	if err := h.incidentUsecase.DeleteIncident(c.Request.Context(), id, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "incident deleted successfully"})
}

// List godoc
// @Summary      List incidents
// @Description  Retrieves a paginated list of incidents with optional filtering by status, severity, team, and search term
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status query []string false "Filter by status (investigating, identified, monitoring, resolved)"
// @Param        severity query []string false "Filter by severity (critical, high, medium, low)"
// @Param        assigned_to_team_id query string false "Filter by assigned team ID" format(uuid)
// @Param        search query string false "Search term for title and description"
// @Param        page query int false "Page number (default: 1)"
// @Param        page_size query int false "Page size (default: 20, max: 100)"
// @Success      200 {object} usecase.ListIncidentsResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents [get]
func (h *IncidentHandler) List(c *gin.Context) {
	var req usecase.ListIncidentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	response, err := h.incidentUsecase.ListIncidents(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AddResponder godoc
// @Summary      Add a responder to an incident
// @Description  Assigns a user as a responder to an incident with a specific role
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        request body usecase.AddResponderRequest true "Add responder request"
// @Success      201 {object} domain.IncidentResponder
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/responders [post]
func (h *IncidentHandler) AddResponder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	var req usecase.AddResponderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	responder, err := h.incidentUsecase.AddResponder(c.Request.Context(), id, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responder)
}

// RemoveResponder godoc
// @Summary      Remove a responder from an incident
// @Description  Removes a user from the list of responders for an incident
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        responderId path string true "Responder ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/responders/{responderId} [delete]
func (h *IncidentHandler) RemoveResponder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	responderIDParam := c.Param("responderId")
	responderID, err := uuid.Parse(responderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid responder ID"})
		return
	}

	userID, _ := middleware.GetUserID(c)

	orgID, _ := middleware.GetOrganizationID(c)

	if err := h.incidentUsecase.RemoveResponder(c.Request.Context(), id, orgID, responderID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "responder removed successfully"})
}

// UpdateResponderRole godoc
// @Summary      Update a responder's role
// @Description  Updates the role of a responder assigned to an incident
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        responderId path string true "Responder ID" format(uuid)
// @Param        request body usecase.UpdateResponderRoleRequest true "Update responder role request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/responders/{responderId} [patch]
func (h *IncidentHandler) UpdateResponderRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	responderIDParam := c.Param("responderId")
	responderID, err := uuid.Parse(responderIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid responder ID"})
		return
	}

	var req usecase.UpdateResponderRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	if err := h.incidentUsecase.UpdateResponderRole(c.Request.Context(), id, orgID, responderID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "responder role updated successfully"})
}

// ListResponders godoc
// @Summary      List responders for an incident
// @Description  Retrieves all responders assigned to an incident with their user details
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Success      200 {array} domain.ResponderWithUser
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/responders [get]
func (h *IncidentHandler) ListResponders(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	responders, err := h.incidentUsecase.ListResponders(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we return an empty array instead of null
	if responders == nil {
		responders = []*domain.ResponderWithUser{}
	}

	c.JSON(http.StatusOK, responders)
}

// AddNote godoc
// @Summary      Add a note to an incident
// @Description  Adds a note to the incident timeline as a timeline event
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        request body usecase.AddNoteRequest true "Add note request"
// @Success      201 {object} domain.IncidentTimelineEvent
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/notes [post]
func (h *IncidentHandler) AddNote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	var req usecase.AddNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)
	userID, _ := middleware.GetUserID(c)

	event, err := h.incidentUsecase.AddNote(c.Request.Context(), id, orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetTimeline godoc
// @Summary      Get incident timeline
// @Description  Retrieves all timeline events for an incident including status changes, notes, and responder actions
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Success      200 {array} domain.TimelineEventWithUser
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/timeline [get]
func (h *IncidentHandler) GetTimeline(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	timeline, err := h.incidentUsecase.GetTimeline(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we return an empty array instead of null
	if timeline == nil {
		timeline = []*domain.TimelineEventWithUser{}
	}

	c.JSON(http.StatusOK, timeline)
}

// LinkAlert godoc
// @Summary      Link an alert to an incident
// @Description  Associates an existing alert with an incident for tracking and correlation
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        request body usecase.LinkAlertRequest true "Link alert request"
// @Success      201 {object} domain.IncidentAlert
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/alerts [post]
func (h *IncidentHandler) LinkAlert(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	var req usecase.LinkAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)
	userID, _ := middleware.GetUserID(c)

	link, err := h.incidentUsecase.LinkAlert(c.Request.Context(), id, orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, link)
}

// UnlinkAlert godoc
// @Summary      Unlink an alert from an incident
// @Description  Removes the association between an alert and an incident
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Param        alertId path string true "Alert ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/alerts/{alertId} [delete]
func (h *IncidentHandler) UnlinkAlert(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	alertIDParam := c.Param("alertId")
	alertID, err := uuid.Parse(alertIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)
	userID, _ := middleware.GetUserID(c)

	if err := h.incidentUsecase.UnlinkAlert(c.Request.Context(), id, orgID, alertID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert unlinked successfully"})
}

// ListAlerts godoc
// @Summary      List alerts linked to an incident
// @Description  Retrieves all alerts that are associated with an incident
// @Tags         Incidents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Incident ID" format(uuid)
// @Success      200 {array} domain.IncidentAlertWithDetails
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /incidents/{id}/alerts [get]
func (h *IncidentHandler) ListAlerts(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid incident ID"})
		return
	}

	orgID, _ := middleware.GetOrganizationID(c)

	alerts, err := h.incidentUsecase.ListAlerts(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we return an empty array instead of null
	if alerts == nil {
		alerts = []*domain.IncidentAlertWithDetails{}
	}

	c.JSON(http.StatusOK, alerts)
}
