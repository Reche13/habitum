package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/middleware"
	"github.com/reche13/habitum/internal/model"
	"github.com/reche13/habitum/internal/model/habit"
	"github.com/reche13/habitum/internal/service"
	"github.com/rs/zerolog"
)

type HabitHandler struct {
	logger      zerolog.Logger
	habitService *service.HabitService
}

func NewHabitHandler(
	habitService *service.HabitService,
) *HabitHandler {
	return &HabitHandler{
		habitService: habitService,
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