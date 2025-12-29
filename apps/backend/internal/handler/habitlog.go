package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/reche13/habitum/internal/model/habitlog"
	"github.com/reche13/habitum/internal/service"
	"github.com/rs/zerolog"
)

type HabitLogHandler struct {
	logger           zerolog.Logger
	habitLogService  *service.HabitLogService
}

func NewHabitLogHandler(
	habitLogService *service.HabitLogService,
) *HabitLogHandler {
	return &HabitLogHandler{
		habitLogService: habitLogService,
	}
}

func hardcodedUserID() uuid.UUID {
	return uuid.MustParse("11111111-1111-1111-1111-111111111111")
}

func (h *HabitLogHandler) Create(c echo.Context) error {
	var payload habitlog.HabitLogPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.habitLogService.SetCompletion(
		c.Request().Context(),
		hardcodedUserID(),
		&payload,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *HabitLogHandler) GetByDate(c echo.Context) error {
	date, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.habitLogService.GetByDate(
		c.Request().Context(),
		hardcodedUserID(),
		date,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *HabitLogHandler) GetByDateRange(c echo.Context) error {
	from, err := time.Parse("2006-01-02", c.QueryParam("from"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	to, err := time.Parse("2006-01-02", c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.habitLogService.GetByDateRange(
		c.Request().Context(),
		hardcodedUserID(),
		from,
		to,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *HabitLogHandler) GetByHabit(c echo.Context) error {
	habitID, err := uuid.Parse(c.Param("habitId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	from, err := time.Parse("2006-01-02", c.QueryParam("from"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	to, err := time.Parse("2006-01-02", c.QueryParam("to"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := h.habitLogService.GetByHabit(
		c.Request().Context(),
		hardcodedUserID(),
		habitID,
		from,
		to,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}
