package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// Schedule handlers

// List godoc
// @Summary      List schedules
// @Description  Retrieves a paginated list of schedules for the current organization
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "Page number"      default(1)
// @Param        page_size  query     int  false  "Page size"        default(20)
// @Success      200        {object}  map[string][]domain.Schedule
// @Failure      401        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /schedules [get]
func (h *ScheduleHandler) List(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	schedules, err := h.scheduleService.ListSchedules(c.Request.Context(), orgID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

// Create godoc
// @Summary      Create schedule
// @Description  Creates a new schedule for the current organization
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      service.CreateScheduleRequest  true  "Schedule creation request"
// @Success      201      {object}  domain.Schedule
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /schedules [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.scheduleService.CreateSchedule(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

// Get godoc
// @Summary      Get schedule
// @Description  Retrieves a schedule by ID including its rotations
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Schedule ID"  format(uuid)
// @Success      200  {object}  domain.ScheduleWithRotations
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /schedules/{id} [get]
func (h *ScheduleHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	schedule, err := h.scheduleService.GetScheduleWithRotations(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// Update godoc
// @Summary      Update schedule
// @Description  Updates an existing schedule by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                          true  "Schedule ID"  format(uuid)
// @Param        request  body      service.UpdateScheduleRequest   true  "Schedule update request"
// @Success      200      {object}  domain.Schedule
// @Failure      400      {object}  map[string]string
// @Router       /schedules/{id} [patch]
func (h *ScheduleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := h.scheduleService.UpdateSchedule(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// Delete godoc
// @Summary      Delete schedule
// @Description  Deletes a schedule by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Schedule ID"  format(uuid)
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /schedules/{id} [delete]
func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	if err := h.scheduleService.DeleteSchedule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}

// Rotation handlers

// ListRotations godoc
// @Summary      List rotations
// @Description  Retrieves all rotations for a schedule
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Schedule ID"  format(uuid)
// @Success      200  {object}  map[string][]domain.ScheduleRotation
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /schedules/{id}/rotations [get]
func (h *ScheduleHandler) ListRotations(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	rotations, err := h.scheduleService.ListRotations(c.Request.Context(), scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rotations": rotations})
}

// CreateRotation godoc
// @Summary      Create rotation
// @Description  Creates a new rotation for a schedule
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                          true  "Schedule ID"  format(uuid)
// @Param        request  body      service.CreateRotationRequest   true  "Rotation creation request"
// @Success      201      {object}  domain.ScheduleRotation
// @Failure      400      {object}  map[string]string
// @Router       /schedules/{id}/rotations [post]
func (h *ScheduleHandler) CreateRotation(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var req service.CreateRotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rotation, err := h.scheduleService.CreateRotation(c.Request.Context(), scheduleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rotation)
}

// GetRotation godoc
// @Summary      Get rotation
// @Description  Retrieves a rotation by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string  true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string  true  "Rotation ID"   format(uuid)
// @Success      200         {object}  domain.ScheduleRotation
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId} [get]
func (h *ScheduleHandler) GetRotation(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	rotation, err := h.scheduleService.GetRotation(c.Request.Context(), rotationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rotation)
}

// UpdateRotation godoc
// @Summary      Update rotation
// @Description  Updates an existing rotation by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string                          true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string                          true  "Rotation ID"   format(uuid)
// @Param        request     body      service.UpdateRotationRequest   true  "Rotation update request"
// @Success      200         {object}  domain.ScheduleRotation
// @Failure      400         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId} [patch]
func (h *ScheduleHandler) UpdateRotation(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	var req service.UpdateRotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rotation, err := h.scheduleService.UpdateRotation(c.Request.Context(), rotationID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rotation)
}

// DeleteRotation godoc
// @Summary      Delete rotation
// @Description  Deletes a rotation by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string  true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string  true  "Rotation ID"   format(uuid)
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId} [delete]
func (h *ScheduleHandler) DeleteRotation(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	if err := h.scheduleService.DeleteRotation(c.Request.Context(), rotationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rotation deleted"})
}

// Participant handlers

// ListParticipants godoc
// @Summary      List participants
// @Description  Retrieves all participants for a rotation
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string  true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string  true  "Rotation ID"   format(uuid)
// @Success      200         {object}  map[string][]domain.ParticipantWithUser
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId}/participants [get]
func (h *ScheduleHandler) ListParticipants(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	participants, err := h.scheduleService.ListParticipants(c.Request.Context(), rotationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

// AddParticipant godoc
// @Summary      Add participant
// @Description  Adds a participant to a rotation
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string                          true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string                          true  "Rotation ID"   format(uuid)
// @Param        request     body      service.AddParticipantRequest   true  "Add participant request"
// @Success      201         {object}  domain.ScheduleRotationParticipant
// @Failure      400         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId}/participants [post]
func (h *ScheduleHandler) AddParticipant(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	var req service.AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participant, err := h.scheduleService.AddParticipant(c.Request.Context(), rotationID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, participant)
}

// RemoveParticipant godoc
// @Summary      Remove participant
// @Description  Removes a participant from a rotation
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string  true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string  true  "Rotation ID"   format(uuid)
// @Param        userId      path      string  true  "User ID"       format(uuid)
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId}/participants/{userId} [delete]
func (h *ScheduleHandler) RemoveParticipant(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.scheduleService.RemoveParticipant(c.Request.Context(), rotationID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "participant removed"})
}

// ReorderParticipants godoc
// @Summary      Reorder participants
// @Description  Reorders participants in a rotation
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string                              true  "Schedule ID"   format(uuid)
// @Param        rotationId  path      string                              true  "Rotation ID"   format(uuid)
// @Param        request     body      service.ReorderParticipantsRequest  true  "Reorder participants request"
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Router       /schedules/{id}/rotations/{rotationId}/participants/reorder [put]
func (h *ScheduleHandler) ReorderParticipants(c *gin.Context) {
	rotationID, err := uuid.Parse(c.Param("rotationId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rotation id"})
		return
	}

	var req service.ReorderParticipantsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.scheduleService.ReorderParticipants(c.Request.Context(), rotationID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "participants reordered"})
}

// Override handlers

// ListOverrides godoc
// @Summary      List overrides
// @Description  Retrieves all overrides for a schedule within a time range
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      string  true   "Schedule ID"                        format(uuid)
// @Param        start  query     string  false  "Start time (RFC3339 format)"        format(date-time)
// @Param        end    query     string  false  "End time (RFC3339 format)"          format(date-time)
// @Success      200    {object}  map[string][]domain.ScheduleOverride
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /schedules/{id}/overrides [get]
func (h *ScheduleHandler) ListOverrides(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	// Parse time range from query params
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
			return
		}
	} else {
		start = time.Now()
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
			return
		}
	} else {
		end = start.AddDate(0, 1, 0) // Default to 1 month
	}

	overrides, err := h.scheduleService.ListOverrides(c.Request.Context(), scheduleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"overrides": overrides})
}

// CreateOverride godoc
// @Summary      Create override
// @Description  Creates a new schedule override
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                          true  "Schedule ID"  format(uuid)
// @Param        request  body      service.CreateOverrideRequest   true  "Override creation request"
// @Success      201      {object}  domain.ScheduleOverride
// @Failure      400      {object}  map[string]string
// @Router       /schedules/{id}/overrides [post]
func (h *ScheduleHandler) CreateOverride(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var req service.CreateOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	override, err := h.scheduleService.CreateOverride(c.Request.Context(), scheduleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, override)
}

// GetOverride godoc
// @Summary      Get override
// @Description  Retrieves an override by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string  true  "Schedule ID"   format(uuid)
// @Param        overrideId  path      string  true  "Override ID"   format(uuid)
// @Success      200         {object}  domain.ScheduleOverride
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Router       /schedules/{id}/overrides/{overrideId} [get]
func (h *ScheduleHandler) GetOverride(c *gin.Context) {
	overrideID, err := uuid.Parse(c.Param("overrideId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid override id"})
		return
	}

	override, err := h.scheduleService.GetOverride(c.Request.Context(), overrideID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, override)
}

// UpdateOverride godoc
// @Summary      Update override
// @Description  Updates an existing override by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string                          true  "Schedule ID"   format(uuid)
// @Param        overrideId  path      string                          true  "Override ID"   format(uuid)
// @Param        request     body      service.UpdateOverrideRequest   true  "Override update request"
// @Success      200         {object}  domain.ScheduleOverride
// @Failure      400         {object}  map[string]string
// @Router       /schedules/{id}/overrides/{overrideId} [patch]
func (h *ScheduleHandler) UpdateOverride(c *gin.Context) {
	overrideID, err := uuid.Parse(c.Param("overrideId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid override id"})
		return
	}

	var req service.UpdateOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	override, err := h.scheduleService.UpdateOverride(c.Request.Context(), overrideID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, override)
}

// DeleteOverride godoc
// @Summary      Delete override
// @Description  Deletes an override by ID
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id          path      string  true  "Schedule ID"   format(uuid)
// @Param        overrideId  path      string  true  "Override ID"   format(uuid)
// @Success      200         {object}  map[string]string
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /schedules/{id}/overrides/{overrideId} [delete]
func (h *ScheduleHandler) DeleteOverride(c *gin.Context) {
	overrideID, err := uuid.Parse(c.Param("overrideId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid override id"})
		return
	}

	if err := h.scheduleService.DeleteOverride(c.Request.Context(), overrideID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "override deleted"})
}

// On-call handler

// GetOnCall godoc
// @Summary      Get on-call user
// @Description  Retrieves the user currently on-call for a schedule at a specific time
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true   "Schedule ID"                        format(uuid)
// @Param        at   query     string  false  "Time to check (RFC3339 format)"     format(date-time)
// @Success      200  {object}  domain.OnCallUser
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /schedules/{id}/oncall [get]
func (h *ScheduleHandler) GetOnCall(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	// Parse optional time parameter
	atStr := c.Query("at")
	var at time.Time
	if atStr != "" {
		at, err = time.Parse(time.RFC3339, atStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
			return
		}
	} else {
		at = time.Now()
	}

	onCallUser, err := h.scheduleService.GetOnCallUser(c.Request.Context(), scheduleID, at)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, onCallUser)
}
