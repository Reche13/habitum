package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/reche13/habitum/internal/errs"
	"github.com/reche13/habitum/internal/service"
)

type CalendarHandler struct {
	calendarService *service.CalendarService
}

func NewCalendarHandler(calendarService *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{
		calendarService: calendarService,
	}
}

func (h *CalendarHandler) GetCompletions(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	// Parse query params
	startDateStr := c.QueryParam("startDate")
	endDateStr := c.QueryParam("endDate")
	habitIDsStr := c.QueryParam("habitIds")

	if startDateStr == "" || endDateStr == "" {
		return errs.NewBadRequestError("startDate and endDate are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return errs.NewBadRequestError("Invalid startDate format. Use yyyy-MM-dd")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return errs.NewBadRequestError("Invalid endDate format. Use yyyy-MM-dd")
	}

	// Parse habit IDs if provided
	var habitIDs []uuid.UUID
	if habitIDsStr != "" {
		ids := strings.Split(habitIDsStr, ",")
		habitIDs = make([]uuid.UUID, 0, len(ids))
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				return errs.NewBadRequestError("Invalid habit ID: " + idStr)
			}
			habitIDs = append(habitIDs, id)
		}
	}

	completions, err := h.calendarService.GetCompletions(
		c.Request().Context(),
		userID,
		startDate,
		endDate,
		habitIDs,
	)
	if err != nil {
		return err
	}

	return c.JSON(200, completions)
}

func (h *CalendarHandler) GetMonth(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	// Parse query params
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	habitIDsStr := c.QueryParam("habitIds")

	if yearStr == "" || monthStr == "" {
		return errs.NewBadRequestError("year and month are required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		return errs.NewBadRequestError("Invalid year")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return errs.NewBadRequestError("Invalid month. Must be 1-12")
	}

	// Parse habit IDs if provided
	var habitIDs []uuid.UUID
	if habitIDsStr != "" {
		ids := strings.Split(habitIDsStr, ",")
		habitIDs = make([]uuid.UUID, 0, len(ids))
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				return errs.NewBadRequestError("Invalid habit ID: " + idStr)
			}
			habitIDs = append(habitIDs, id)
		}
	}

	monthData, err := h.calendarService.GetMonth(
		c.Request().Context(),
		userID,
		year,
		month,
		habitIDs,
	)
	if err != nil {
		return err
	}

	return c.JSON(200, monthData)
}

func (h *CalendarHandler) GetWeek(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	// Parse query params
	yearStr := c.QueryParam("year")
	weekStr := c.QueryParam("week")
	habitIDsStr := c.QueryParam("habitIds")

	if yearStr == "" || weekStr == "" {
		return errs.NewBadRequestError("year and week are required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		return errs.NewBadRequestError("Invalid year")
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 53 {
		return errs.NewBadRequestError("Invalid week. Must be 1-53")
	}

	// Parse habit IDs if provided
	var habitIDs []uuid.UUID
	if habitIDsStr != "" {
		ids := strings.Split(habitIDsStr, ",")
		habitIDs = make([]uuid.UUID, 0, len(ids))
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				return errs.NewBadRequestError("Invalid habit ID: " + idStr)
			}
			habitIDs = append(habitIDs, id)
		}
	}

	weekData, err := h.calendarService.GetWeek(
		c.Request().Context(),
		userID,
		year,
		week,
		habitIDs,
	)
	if err != nil {
		return err
	}

	return c.JSON(200, weekData)
}

func (h *CalendarHandler) GetYear(c echo.Context) error {
	userID := uuid.MustParse("04b151e6-7631-4548-9384-1e11bbaa84e8")

	// Parse query params
	yearStr := c.QueryParam("year")
	habitIDsStr := c.QueryParam("habitIds")

	if yearStr == "" {
		return errs.NewBadRequestError("year is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		return errs.NewBadRequestError("Invalid year")
	}

	// Parse habit IDs if provided
	var habitIDs []uuid.UUID
	if habitIDsStr != "" {
		ids := strings.Split(habitIDsStr, ",")
		habitIDs = make([]uuid.UUID, 0, len(ids))
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				return errs.NewBadRequestError("Invalid habit ID: " + idStr)
			}
			habitIDs = append(habitIDs, id)
		}
	}

	yearData, err := h.calendarService.GetYear(
		c.Request().Context(),
		userID,
		year,
		habitIDs,
	)
	if err != nil {
		return err
	}

	return c.JSON(200, yearData)
}

