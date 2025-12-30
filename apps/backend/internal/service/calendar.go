package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/lib"
	"github.com/reche13/habitum/internal/model/calendar"
	"github.com/reche13/habitum/internal/model/habitlog"
	"github.com/reche13/habitum/internal/repository"
)

type CalendarService struct {
	*BaseService
	habitRepo    *repository.HabitRepository
	habitLogRepo *repository.HabitLogRepository
}

func NewCalendarService(
	habitRepo *repository.HabitRepository,
	habitLogRepo *repository.HabitLogRepository,
) *CalendarService {
	return &CalendarService{
		BaseService: &BaseService{
			resourceName: "calendar",
		},
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
	}
}

func (s *CalendarService) GetCompletions(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate time.Time,
	habitIDs []uuid.UUID,
) (*calendar.CompletionsResponse, error) {
	normalizedStart := lib.NormalizeDate(startDate)
	normalizedEnd := lib.NormalizeDate(endDate)

	// Get all logs in date range
	allLogs, err := s.habitLogRepo.GetByDateRange(ctx, userID, normalizedStart, normalizedEnd)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter by habit IDs if provided, and only completed logs
	logs := make([]habitlog.HabitLog, 0)
	if len(habitIDs) > 0 {
		habitIDMap := make(map[uuid.UUID]bool)
		for _, id := range habitIDs {
			habitIDMap[id] = true
		}
		for _, log := range allLogs {
			if habitIDMap[log.HabitID] && log.Completed {
				logs = append(logs, log)
			}
		}
	} else {
		// Only completed logs
		for _, log := range allLogs {
			if log.Completed {
				logs = append(logs, log)
			}
		}
	}

	// Get habits for habit info
	allHabits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	habitMap := make(map[uuid.UUID]struct {
		ID    uuid.UUID
		Name  string
		Color *string
		Icon  *string
	})
	for _, h := range allHabits {
		if h.ArchivedAt == nil {
			habitMap[h.ID] = struct {
				ID    uuid.UUID
				Name  string
				Color *string
				Icon  *string
			}{h.ID, h.Name, h.Color, h.Icon}
		}
	}

	// Filter habits by habitIDs if provided
	if len(habitIDs) > 0 {
		filteredHabits := make(map[uuid.UUID]struct {
			ID    uuid.UUID
			Name  string
			Color *string
			Icon  *string
		})
		for _, id := range habitIDs {
			if h, ok := habitMap[id]; ok {
				filteredHabits[id] = h
			}
		}
		habitMap = filteredHabits
	}

	// Group logs by date
	completionsByDate := make(map[string][]uuid.UUID) // date -> habit IDs
	for _, log := range logs {
		if log.Completed {
			dateKey := log.LogDate.Format("2006-01-02")
			completionsByDate[dateKey] = append(completionsByDate[dateKey], log.HabitID)
		}
	}

	// Build completion days
	completionDays := make([]calendar.CompletionDay, 0)
	currentDate := normalizedStart
	totalCompletions := 0
	daysWithCompletions := 0

	for !currentDate.After(normalizedEnd) {
		dateKey := currentDate.Format("2006-01-02")
		completedHabitIDs := completionsByDate[dateKey]

		// Count active habits on this date
		totalHabits := 0
		completedHabits := make([]calendar.HabitInfo, 0)

		for habitID, habitInfo := range habitMap {
			// Check if habit was active on this date (simplified - assume all habits are active)
			totalHabits++
			for _, completedID := range completedHabitIDs {
				if completedID == habitID {
					completedHabits = append(completedHabits, calendar.HabitInfo{
						ID:    habitInfo.ID.String(),
						Name:  habitInfo.Name,
						Color: lib.GetStringValue(habitInfo.Color),
						Icon:  lib.GetStringValue(habitInfo.Icon),
					})
					break
				}
			}
		}

		completionRate := 0.0
		if totalHabits > 0 {
			completionRate = (float64(len(completedHabits)) / float64(totalHabits)) * 100
		}

		if len(completedHabits) > 0 {
			daysWithCompletions++
			totalCompletions += len(completedHabits)
		}

		completionDays = append(completionDays, calendar.CompletionDay{
			Date:            dateKey,
			Habits:          completedHabits,
			CompletionRate:  completionRate,
			TotalHabits:     totalHabits,
			CompletedHabits: len(completedHabits),
		})

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Calculate statistics
	totalDays := int(normalizedEnd.Sub(normalizedStart).Hours()/24) + 1
	overallCompletionRate := 0.0
	if totalDays > 0 && len(habitMap) > 0 {
		expectedCompletions := totalDays * len(habitMap)
		if expectedCompletions > 0 {
			overallCompletionRate = (float64(totalCompletions) / float64(expectedCompletions)) * 100
		}
	}

	return &calendar.CompletionsResponse{
		Completions: completionDays,
		Statistics: calendar.PeriodStats{
			TotalCompletions:    totalCompletions,
			DaysWithCompletions: daysWithCompletions,
			CompletionRate:      overallCompletionRate,
			TotalDays:           totalDays,
		},
	}, nil
}

func (s *CalendarService) GetMonth(
	ctx context.Context,
	userID uuid.UUID,
	year int,
	month int,
	habitIDs []uuid.UUID,
) (*calendar.MonthResponse, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of month

	completions, err := s.GetCompletions(ctx, userID, startDate, endDate, habitIDs)
	if err != nil {
		return nil, err
	}

	// Convert to day data format
	days := make([]calendar.DayData, 0)
	for _, completion := range completions.Completions {
		habitIDStrings := make([]string, 0)
		for _, habit := range completion.Habits {
			habitIDStrings = append(habitIDStrings, habit.ID)
		}

		days = append(days, calendar.DayData{
			Date:           completion.Date,
			Completions:    habitIDStrings,
			CompletionRate: completion.CompletionRate,
		})
	}

	return &calendar.MonthResponse{
		Year:       year,
		Month:      month,
		Days:       days,
		Statistics: completions.Statistics,
	}, nil
}

func (s *CalendarService) GetWeek(
	ctx context.Context,
	userID uuid.UUID,
	year int,
	week int,
	habitIDs []uuid.UUID,
) (*calendar.WeekResponse, error) {
	// Calculate start date of week (Monday)
	date := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC) // Jan 4 is always in week 1
	dateYear, dateWeek := date.ISOWeek()
	for dateYear != year || dateWeek != week {
		if dateWeek < week {
			date = date.AddDate(0, 0, 7)
		} else {
			date = date.AddDate(0, 0, -7)
		}
		dateYear, dateWeek = date.ISOWeek()
	}

	// Get Monday of that week
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}

	startDate := lib.NormalizeDate(date)
	endDate := startDate.AddDate(0, 0, 6) // Sunday

	completions, err := s.GetCompletions(ctx, userID, startDate, endDate, habitIDs)
	if err != nil {
		return nil, err
	}

	// Convert to day data format with day of week
	days := make([]calendar.DayData, 0)
	for _, completion := range completions.Completions {
		parsedDate, _ := time.Parse("2006-01-02", completion.Date)
		weekday := int(parsedDate.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday becomes 7
		}
		dayOfWeek := weekday - 1 // Monday = 0, Sunday = 6

		habitIDStrings := make([]string, 0)
		for _, habit := range completion.Habits {
			habitIDStrings = append(habitIDStrings, habit.ID)
		}

		days = append(days, calendar.DayData{
			Date:           completion.Date,
			DayOfWeek:      dayOfWeek,
			Completions:    habitIDStrings,
			CompletionRate: completion.CompletionRate,
		})
	}

	return &calendar.WeekResponse{
		Year:       year,
		Week:       week,
		Days:       days,
		Statistics: completions.Statistics,
	}, nil
}

func (s *CalendarService) GetYear(
	ctx context.Context,
	userID uuid.UUID,
	year int,
	habitIDs []uuid.UUID,
) (*calendar.YearResponse, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	completions, err := s.GetCompletions(ctx, userID, startDate, endDate, habitIDs)
	if err != nil {
		return nil, err
	}

	// Convert to heatmap format
	heatmap := make([]calendar.HeatmapDay, 0)
	for _, completion := range completions.Completions {
		// Calculate intensity (0-4) based on completion rate
		intensity := 0
		if completion.CompletionRate > 0 {
			intensity = 1
		}
		if completion.CompletionRate >= 25 {
			intensity = 2
		}
		if completion.CompletionRate >= 50 {
			intensity = 3
		}
		if completion.CompletionRate >= 75 {
			intensity = 4
		}

		heatmap = append(heatmap, calendar.HeatmapDay{
			Date:           completion.Date,
			CompletionRate: completion.CompletionRate,
			Intensity:      intensity,
		})
	}

	return &calendar.YearResponse{
		Year:       year,
		Heatmap:    heatmap,
		Statistics: completions.Statistics,
	}, nil
}


