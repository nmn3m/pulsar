package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pulsar/backend/internal/middleware"
	"github.com/pulsar/backend/internal/service"
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
