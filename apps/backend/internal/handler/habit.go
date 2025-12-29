package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/middleware"
	"github.com/reche13/habitum/internal/model"
	"github.com/reche13/habitum/internal/model/habit"
	"github.com/reche13/habitum/internal/model/habitlog"
	"github.com/reche13/habitum/internal/service"
	"github.com/rs/zerolog"
)

type HabitHandler struct {
	logger         zerolog.Logger
	habitService   *service.HabitService
	habitLogService *service.HabitLogService
}

func NewHabitHandler(
	habitService *service.HabitService,
	habitLogService *service.HabitLogService,
) *HabitHandler {
	return &HabitHandler{
		habitService:   habitService,
		habitLogService: habitLogService,
	}
}

func (h *HabitHandler) CreateHabit(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	var payload habit.CreateHabitPayload

	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&payload); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	createdHabit, err := h.habitService.CreateHabit(c.Request().Context(), userID, &payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, model.SuccessResponse(createdHabit))
}

func (h *HabitHandler) GetHabits(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	// Parse query parameters
	filters := &habit.ListFilters{}

	// Category filter
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		category := habit.Category(categoryParam)
		filters.Category = &category
	}

	// Search filter
	if searchParam := c.QueryParam("search"); searchParam != "" {
		filters.Search = &searchParam
	}

	// Sort parameter
	if sortParam := c.QueryParam("sort"); sortParam != "" {
		// Validate sort option
		validSorts := map[string]bool{
			"name":       true,
			"date":       true,
			"streak":     true,
			"completion": true,
		}
		if validSorts[sortParam] {
			filters.Sort = &sortParam
		}
	}

	// Order parameter (asc/desc)
	if orderParam := c.QueryParam("order"); orderParam != "" {
		orderLower := strings.ToLower(orderParam)
		if orderLower == "asc" || orderLower == "desc" {
			filters.Order = &orderLower
		}
	}

	// Pagination parameters
	if pageParam := c.QueryParam("page"); pageParam != "" {
		if page, err := strconv.Atoi(pageParam); err == nil && page > 0 {
			filters.Page = &page
		}
	}

	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil && limit > 0 {
			filters.Limit = &limit
		}
	}

	habits, total, err := h.habitService.GetHabits(c.Request().Context(), userID, filters)
	if err != nil {
		return err
	}

	// Calculate pagination metadata
	page := 1
	limit := 50
	if filters.Page != nil {
		page = *filters.Page
	}
	if filters.Limit != nil {
		limit = *filters.Limit
	}
	totalPages := (total + limit - 1) / limit // Ceiling division

	meta := &model.Meta{
		RequestID:  middleware.GetRequestID(c),
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, model.SuccessResponseWithMeta(habits, meta))
}

func (h *HabitHandler) GetHabit(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	habit, err := h.habitService.GetHabit(c.Request().Context(), habitID, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(habit))
}

func (h *HabitHandler) UpdateHabit(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	var payload habit.UpdateHabitPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("Invalid request payload")
	}

	if fieldErrors := middleware.ValidateStruct(&payload); fieldErrors != nil {
		return errs.NewValidationError(fieldErrors)
	}

	updatedHabit, err := h.habitService.UpdateHabit(c.Request().Context(), habitID, userID, &payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(updatedHabit))
}

func (h *HabitHandler) DeleteHabit(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	if err := h.habitService.DeleteHabit(c.Request().Context(), habitID, userID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HabitHandler) MarkComplete(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	// Parse optional date parameter (defaults to today)
	var logDate time.Time
	if dateParam := c.QueryParam("date"); dateParam != "" {
		parsedDate, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			return errs.NewBadRequestError("Invalid date format. Use YYYY-MM-DD")
		}
		logDate = parsedDate
	} else {
		logDate = time.Now().UTC()
	}

	// Verify habit exists and belongs to user
	_, err = h.habitService.GetHabit(c.Request().Context(), habitID, userID)
	if err != nil {
		return err
	}

	// Mark as complete
	payload := &habitlog.HabitLogPayload{
		HabitID:   habitID,
		LogDate:   logDate,
		Completed: true,
	}

	log, err := h.habitLogService.MarkComplete(c.Request().Context(), userID, habitID, payload)
	if err != nil {
		return err
	}

	// Get updated habit with computed fields
	updatedHabit, err := h.habitService.GetHabit(c.Request().Context(), habitID, userID)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"id":          log.ID,
		"habitId":     log.HabitID,
		"userId":      log.UserID,
		"completedAt": log.CreatedAt,
		"date":        log.LogDate.Format("2006-01-02"),
		"habit":       updatedHabit,
	}

	return c.JSON(http.StatusCreated, model.SuccessResponse(response))
}

func (h *HabitHandler) UnmarkComplete(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	// Parse optional date parameter (defaults to today)
	var logDate time.Time
	if dateParam := c.QueryParam("date"); dateParam != "" {
		parsedDate, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			return errs.NewBadRequestError("Invalid date format. Use YYYY-MM-DD")
		}
		logDate = parsedDate
	} else {
		logDate = time.Now().UTC()
	}

	// Verify habit exists and belongs to user
	_, err = h.habitService.GetHabit(c.Request().Context(), habitID, userID)
	if err != nil {
		return err
	}

	// Unmark completion
	err = h.habitLogService.UnmarkComplete(c.Request().Context(), userID, habitID, logDate)
	if err != nil {
		return err
	}

	// Get updated habit with computed fields
	updatedHabit, err := h.habitService.GetHabit(c.Request().Context(), habitID, userID)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"message": "Completion removed",
		"habit":   updatedHabit,
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(response))
}

func (h *HabitHandler) GetCompletions(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	// Parse optional date range parameters
	var startDate, endDate time.Time
	if startParam := c.QueryParam("startDate"); startParam != "" {
		startDate, err = time.Parse("2006-01-02", startParam)
		if err != nil {
			return errs.NewBadRequestError("Invalid startDate format. Use YYYY-MM-DD")
		}
	} else {
		// Default to 365 days ago
		startDate = time.Now().UTC().AddDate(0, 0, -365)
	}

	if endParam := c.QueryParam("endDate"); endParam != "" {
		endDate, err = time.Parse("2006-01-02", endParam)
		if err != nil {
			return errs.NewBadRequestError("Invalid endDate format. Use YYYY-MM-DD")
		}
	} else {
		endDate = time.Now().UTC()
	}

	// Parse limit
	limit := 365 // default
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 1000 {
				limit = 1000 // cap at 1000
			}
		}
	}

	completions, total, err := h.habitLogService.GetCompletions(c.Request().Context(), userID, habitID, startDate, endDate, limit)
	if err != nil {
		return err
	}

	// Format response
	completionRecords := make([]map[string]interface{}, len(completions))
	for i, comp := range completions {
		completionRecords[i] = map[string]interface{}{
			"id":          comp.ID,
			"habitId":     comp.HabitID,
			"date":        comp.LogDate.Format("2006-01-02"),
			"completedAt": comp.CreatedAt,
		}
	}

	response := map[string]interface{}{
		"completions": completionRecords,
		"total":       total,
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(response))
}

func (h *HabitHandler) GetCompletionHistory(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	idParam := c.Param("id")
	habitID, err := uuid.Parse(idParam)
	if err != nil {
		return errs.NewBadRequestError("Invalid habit ID format")
	}

	// Parse query parameters
	allTime := c.QueryParam("allTime") == "true"
	var year *int
	if yearParam := c.QueryParam("year"); yearParam != "" {
		if parsedYear, err := strconv.Atoi(yearParam); err == nil && parsedYear > 0 {
			year = &parsedYear
		}
	}

	dates, totalDays, completedDays, err := h.habitLogService.GetCompletionHistory(c.Request().Context(), userID, habitID, year, allTime)
	if err != nil {
		return err
	}

	// Format dates as strings
	dateStrings := make([]string, len(dates))
	for i, date := range dates {
		dateStrings[i] = date.Format("2006-01-02")
	}

	response := map[string]interface{}{
		"dates":         dateStrings,
		"totalDays":     totalDays,
		"completedDays": completedDays,
	}

	return c.JSON(http.StatusOK, model.SuccessResponse(response))
}