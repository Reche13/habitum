package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/lib"
	"github.com/reche13/habitum/internal/model/analytics"
	"github.com/reche13/habitum/internal/repository"
)

type AnalyticsService struct {
	*BaseService
	habitRepo    *repository.HabitRepository
	habitLogRepo *repository.HabitLogRepository
}

func NewAnalyticsService(
	habitRepo *repository.HabitRepository,
	habitLogRepo *repository.HabitLogRepository,
) *AnalyticsService {
	return &AnalyticsService{
		BaseService: &BaseService{
			resourceName: "analytics",
		},
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
	}
}

func (s *AnalyticsService) GetCompletionTrend(ctx context.Context, userID uuid.UUID, period string) (*analytics.CompletionTrendResponse, error) {
	// Calculate date range based on period
	now := lib.NormalizeDate(time.Now().UTC())
	var startDate time.Time

	switch period {
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "90d":
		startDate = now.AddDate(0, 0, -90)
	case "all":
		// Get the earliest habit creation date
		habits, _, err := s.habitRepo.List(ctx, userID, nil)
		if err != nil {
			return nil, s.wrapError(err)
		}
		if len(habits) == 0 {
			return &analytics.CompletionTrendResponse{Data: []analytics.CompletionTrendDataPoint{}}, nil
		}
		startDate = habits[0].CreatedAt
		for _, h := range habits {
			if h.CreatedAt.Before(startDate) {
				startDate = h.CreatedAt
			}
		}
		startDate = lib.NormalizeDate(startDate)
	default:
		// Default to 30 days
		startDate = now.AddDate(0, 0, -30)
	}

	// Get all habits (including archived - we'll check archived date per day)
	allHabits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Get all logs in the date range
	logs, err := s.habitLogRepo.GetByDateRange(ctx, userID, startDate, now)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Group logs by date and count completions
	completionsByDate := make(map[string]int) // date -> count of completions
	
	for _, log := range logs {
		if log.Completed {
			dateKey := log.LogDate.Format("2006-01-02")
			completionsByDate[dateKey]++
		}
	}

	// Build response data points for each day in range
	dataPoints := make([]analytics.CompletionTrendDataPoint, 0)
	currentDate := startDate
	
	for !currentDate.After(now) {
		dateKey := currentDate.Format("2006-01-02")
		completions := completionsByDate[dateKey]
		
		// Count active habits on this date (habits created on or before this date and not archived on this date)
		totalHabits := 0
		for _, h := range allHabits {
			habitCreatedDate := lib.NormalizeDate(h.CreatedAt)
			// Habit must be created on or before this date
			if !currentDate.Before(habitCreatedDate) {
				// Check if habit was archived on or before this date
				if h.ArchivedAt == nil {
					totalHabits++
				} else {
					archivedDate := lib.NormalizeDate(*h.ArchivedAt)
					if currentDate.Before(archivedDate) {
						totalHabits++
					}
				}
			}
		}

		// Calculate completion rate
		completionRate := 0.0
		if totalHabits > 0 {
			completionRate = (float64(completions) / float64(totalHabits)) * 100
		}

		dataPoints = append(dataPoints, analytics.CompletionTrendDataPoint{
			Date:           dateKey,
			Completions:    completions,
			TotalHabits:    totalHabits,
			CompletionRate: completionRate,
		})

		// Move to next day
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &analytics.CompletionTrendResponse{
		Data: dataPoints,
	}, nil
}

func (s *AnalyticsService) GetCategoryBreakdown(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (s *AnalyticsService) GetDayOfWeekAnalysis(ctx context.Context, userID uuid.UUID, period string) (interface{}, error) {
	_ = period // Placeholder
	return nil, nil
}

func (s *AnalyticsService) GetMetrics(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (s *AnalyticsService) GetTopHabits(ctx context.Context, userID uuid.UUID, limit int, sortBy string) (interface{}, error) {
	_ = limit
	_ = sortBy
	return nil, nil
}

func (s *AnalyticsService) GetStreakLeaderboard(ctx context.Context, userID uuid.UUID, limit int) (interface{}, error) {
	_ = limit
	return nil, nil
}

func (s *AnalyticsService) GetInsights(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	return nil, nil
}

