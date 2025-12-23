package handler

import (
	"net/http"

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

	habits, err := h.habitService.GetHabits(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	meta := &model.Meta{
		RequestID: middleware.GetRequestID(c),
		Total:     len(habits),
	}

	return c.JSON(http.StatusOK, model.SuccessResponseWithMeta(habits, meta))
}