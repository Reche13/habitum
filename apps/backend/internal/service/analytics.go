package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/lib"
	"github.com/reche13/habitum/internal/model/analytics"
	"github.com/reche13/habitum/internal/model/habit"
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

func (s *AnalyticsService) GetCategoryBreakdown(ctx context.Context, userID uuid.UUID) (*analytics.CategoryBreakdownResponse, error) {
	// Get all active habits
	habits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter out archived habits
	activeHabits := make([]habit.Habit, 0)
	for _, h := range habits {
		if h.ArchivedAt == nil {
			activeHabits = append(activeHabits, h)
		}
	}

	// Group habits by category and calculate stats
	categoryStats := make(map[string]struct {
		habitCount int
		completionRates []float64
		totalCompletions int
	})

	for _, h := range activeHabits {
		// Calculate completion rate for this habit
		completionRate, err := CalculateCompletionRate(
			ctx,
			s.habitLogRepo,
			userID,
			h.ID,
			h.CreatedAt,
			h.Frequency,
		)
		if err != nil {
			completionRate = 0
		}

		// Get total completions for this habit
		startDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := lib.NormalizeDate(time.Now().UTC())
		logs, err := s.habitLogRepo.GetByHabit(ctx, userID, h.ID, startDate, endDate)
		totalCompletions := 0
		if err == nil {
			for _, log := range logs {
				if log.Completed {
					totalCompletions++
				}
			}
		}

		categoryKey := string(h.Category)
		stats := categoryStats[categoryKey]
		stats.habitCount++
		stats.completionRates = append(stats.completionRates, completionRate)
		stats.totalCompletions += totalCompletions
		categoryStats[categoryKey] = stats
	}

	// Build response
	dataPoints := make([]analytics.CategoryBreakdownDataPoint, 0)
	for category, stats := range categoryStats {
		// Calculate average completion rate
		avgCompletionRate := 0.0
		if len(stats.completionRates) > 0 {
			sum := 0.0
			for _, rate := range stats.completionRates {
				sum += rate
			}
			avgCompletionRate = sum / float64(len(stats.completionRates))
		}

		// Get category label (capitalize first letter)
		label := category
		if len(label) > 0 {
			label = strings.ToUpper(string(label[0])) + strings.ToLower(label[1:])
		}

		dataPoints = append(dataPoints, analytics.CategoryBreakdownDataPoint{
			Category: category,
			Label: label,
			HabitCount: stats.habitCount,
			AvgCompletionRate: avgCompletionRate,
			TotalCompletions: stats.totalCompletions,
		})
	}

	return &analytics.CategoryBreakdownResponse{
		Data: dataPoints,
	}, nil
}

func (s *AnalyticsService) GetDayOfWeekAnalysis(ctx context.Context, userID uuid.UUID, period string) (*analytics.DayOfWeekResponse, error) {
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
			return &analytics.DayOfWeekResponse{Data: []analytics.DayOfWeekDataPoint{}}, nil
		}
		startDate = habits[0].CreatedAt
		for _, h := range habits {
			if h.CreatedAt.Before(startDate) {
				startDate = h.CreatedAt
			}
		}
		startDate = lib.NormalizeDate(startDate)
	default:
		// Default to all time if period is empty or invalid
		if period == "" {
			habits, _, err := s.habitRepo.List(ctx, userID, nil)
			if err != nil {
				return nil, s.wrapError(err)
			}
			if len(habits) == 0 {
				return &analytics.DayOfWeekResponse{Data: []analytics.DayOfWeekDataPoint{}}, nil
			}
			startDate = habits[0].CreatedAt
			for _, h := range habits {
				if h.CreatedAt.Before(startDate) {
					startDate = h.CreatedAt
				}
			}
			startDate = lib.NormalizeDate(startDate)
		} else {
			// Default to 30 days for invalid period
			startDate = now.AddDate(0, 0, -30)
		}
	}

	// Get all habits
	allHabits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Get all logs in the date range
	logs, err := s.habitLogRepo.GetByDateRange(ctx, userID, startDate, now)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Group completions by day of week (0=Monday, 6=Sunday)
	completionsByDay := make(map[int]int) // day index -> count
	habitsByDay := make(map[int]int)      // day index -> total habits active on that day

	// Count completions by day of week
	for _, log := range logs {
		if log.Completed {
			// Convert to weekday (Monday=0, Sunday=6)
			weekday := int(log.LogDate.Weekday())
			if weekday == 0 {
				weekday = 7 // Sunday becomes 7, then we'll convert to 6
			}
			dayIndex := weekday - 1 // Monday=0, Sunday=6
			completionsByDay[dayIndex]++
		}
	}

	// Count total habits active per day of week
	// For each day in the range, count active habits
	currentDate := startDate
	for !currentDate.After(now) {
		weekday := int(currentDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		dayIndex := weekday - 1

		// Count active habits on this date
		activeHabits := 0
		for _, h := range allHabits {
			habitCreatedDate := lib.NormalizeDate(h.CreatedAt)
			if !currentDate.Before(habitCreatedDate) {
				if h.ArchivedAt == nil {
					activeHabits++
				} else {
					archivedDate := lib.NormalizeDate(*h.ArchivedAt)
					if currentDate.Before(archivedDate) {
						activeHabits++
					}
				}
			}
		}
		habitsByDay[dayIndex] += activeHabits

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Build response data points for each day of week
	dayNames := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	dataPoints := make([]analytics.DayOfWeekDataPoint, 7)

	for dayIndex := 0; dayIndex < 7; dayIndex++ {
		completions := completionsByDay[dayIndex]
		totalHabits := habitsByDay[dayIndex]
		
		// Calculate completion rate
		completionRate := 0.0
		if totalHabits > 0 {
			completionRate = (float64(completions) / float64(totalHabits)) * 100
		}

		dataPoints[dayIndex] = analytics.DayOfWeekDataPoint{
			Day:            dayNames[dayIndex],
			DayIndex:       dayIndex,
			Completions:    completions,
			TotalHabits:    totalHabits,
			CompletionRate: completionRate,
		}
	}

	return &analytics.DayOfWeekResponse{
		Data: dataPoints,
	}, nil
}

func (s *AnalyticsService) GetMetrics(ctx context.Context, userID uuid.UUID) (*analytics.MetricsResponse, error) {
	// Get all active habits
	habits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter out archived habits
	activeHabits := make([]habit.Habit, 0)
	for _, h := range habits {
		if h.ArchivedAt == nil {
			activeHabits = append(activeHabits, h)
		}
	}

	if len(activeHabits) == 0 {
		return &analytics.MetricsResponse{
			AvgCompletionRate: 0,
			AvgStreak:         0,
			TotalCompletions:  0,
			ConsistencyScore:  0,
		}, nil
	}

	// Calculate metrics
	totalCompletionRate := 0.0
	totalStreak := 0
	totalCompletions := 0

	for _, h := range activeHabits {
		// Get completion rate
		completionRate, err := CalculateCompletionRate(
			ctx,
			s.habitLogRepo,
			userID,
			h.ID,
			h.CreatedAt,
			h.Frequency,
		)
		if err != nil {
			completionRate = 0
		}
		totalCompletionRate += completionRate

		// Get current streak
		totalStreak += h.CurrentStreak

		// Get total completions for this habit
		startDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := lib.NormalizeDate(time.Now().UTC())
		logs, err := s.habitLogRepo.GetByHabit(ctx, userID, h.ID, startDate, endDate)
		if err == nil {
			for _, log := range logs {
				if log.Completed {
					totalCompletions++
				}
			}
		}
	}

	// Calculate averages
	avgCompletionRate := totalCompletionRate / float64(len(activeHabits))
	avgStreak := float64(totalStreak) / float64(len(activeHabits))

	// Calculate consistency score (based on completion rate - higher is better)
	// Simple approach: use average completion rate as consistency score
	// Could be enhanced with variance calculation, but keeping it simple
	consistencyScore := avgCompletionRate

	return &analytics.MetricsResponse{
		AvgCompletionRate: avgCompletionRate,
		AvgStreak:        avgStreak,
		TotalCompletions: totalCompletions,
		ConsistencyScore: consistencyScore,
	}, nil
}

func (s *AnalyticsService) GetTopHabits(ctx context.Context, userID uuid.UUID, limit int, sortBy string) (*analytics.TopHabitsResponse, error) {
	// Get all active habits
	habits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter out archived habits
	activeHabits := make([]habit.Habit, 0)
	for _, h := range habits {
		if h.ArchivedAt == nil {
			activeHabits = append(activeHabits, h)
		}
	}

	if len(activeHabits) == 0 {
		return &analytics.TopHabitsResponse{Data: []analytics.TopHabitDataPoint{}}, nil
	}

	// Calculate completion rate and prepare data for sorting
	type habitWithStats struct {
		habit          habit.Habit
		completionRate float64
	}

	habitsWithStats := make([]habitWithStats, 0, len(activeHabits))
	for _, h := range activeHabits {
		completionRate, err := CalculateCompletionRate(
			ctx,
			s.habitLogRepo,
			userID,
			h.ID,
			h.CreatedAt,
			h.Frequency,
		)
		if err != nil {
			completionRate = 0
		}

		habitsWithStats = append(habitsWithStats, habitWithStats{
			habit:          h,
			completionRate: completionRate,
		})
	}

	// Sort by sortBy parameter
	switch sortBy {
	case "completion":
		// Sort by completion rate (descending)
		for i := 0; i < len(habitsWithStats)-1; i++ {
			for j := i + 1; j < len(habitsWithStats); j++ {
				if habitsWithStats[i].completionRate < habitsWithStats[j].completionRate {
					habitsWithStats[i], habitsWithStats[j] = habitsWithStats[j], habitsWithStats[i]
				}
			}
		}
	case "streak":
		// Sort by current streak (descending)
		for i := 0; i < len(habitsWithStats)-1; i++ {
			for j := i + 1; j < len(habitsWithStats); j++ {
				if habitsWithStats[i].habit.CurrentStreak < habitsWithStats[j].habit.CurrentStreak {
					habitsWithStats[i], habitsWithStats[j] = habitsWithStats[j], habitsWithStats[i]
				}
			}
		}
	default:
		// Default to completion rate
		for i := 0; i < len(habitsWithStats)-1; i++ {
			for j := i + 1; j < len(habitsWithStats); j++ {
				if habitsWithStats[i].completionRate < habitsWithStats[j].completionRate {
					habitsWithStats[i], habitsWithStats[j] = habitsWithStats[j], habitsWithStats[i]
				}
			}
		}
	}

	// Apply limit
	if limit > 0 && limit < len(habitsWithStats) {
		habitsWithStats = habitsWithStats[:limit]
	}

	// Build response
	dataPoints := make([]analytics.TopHabitDataPoint, len(habitsWithStats))
	for i, hws := range habitsWithStats {
		dataPoints[i] = analytics.TopHabitDataPoint{
			HabitID:        hws.habit.ID.String(),
			Name:           hws.habit.Name,
			Category:       string(hws.habit.Category),
			CompletionRate: hws.completionRate,
			CurrentStreak:  hws.habit.CurrentStreak,
			LongestStreak:  hws.habit.LongestStreak,
		}
	}

	return &analytics.TopHabitsResponse{
		Data: dataPoints,
	}, nil
}

func (s *AnalyticsService) GetStreakLeaderboard(ctx context.Context, userID uuid.UUID, limit int) (*analytics.StreakLeaderboardResponse, error) {
	// Get all active habits
	habits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter out archived habits
	activeHabits := make([]habit.Habit, 0)
	for _, h := range habits {
		if h.ArchivedAt == nil {
			activeHabits = append(activeHabits, h)
		}
	}

	if len(activeHabits) == 0 {
		return &analytics.StreakLeaderboardResponse{Data: []analytics.StreakLeaderboardDataPoint{}}, nil
	}

	// Sort by current streak (descending)
	for i := 0; i < len(activeHabits)-1; i++ {
		for j := i + 1; j < len(activeHabits); j++ {
			if activeHabits[i].CurrentStreak < activeHabits[j].CurrentStreak {
				activeHabits[i], activeHabits[j] = activeHabits[j], activeHabits[i]
			}
		}
	}

	// Apply limit
	if limit > 0 && limit < len(activeHabits) {
		activeHabits = activeHabits[:limit]
	}

	// Build response
	dataPoints := make([]analytics.StreakLeaderboardDataPoint, len(activeHabits))
	for i, h := range activeHabits {
		dataPoints[i] = analytics.StreakLeaderboardDataPoint{
			HabitID:       h.ID.String(),
			Name:          h.Name,
			Category:      string(h.Category),
			CurrentStreak: h.CurrentStreak,
			LongestStreak: h.LongestStreak,
		}
	}

	return &analytics.StreakLeaderboardResponse{
		Data: dataPoints,
	}, nil
}

func (s *AnalyticsService) GetInsights(ctx context.Context, userID uuid.UUID) (*analytics.InsightsResponse, error) {
	insights := make([]analytics.Insight, 0)

	// Get all active habits
	habits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter out archived habits
	activeHabits := make([]habit.Habit, 0)
	for _, h := range habits {
		if h.ArchivedAt == nil {
			activeHabits = append(activeHabits, h)
		}
	}

	if len(activeHabits) == 0 {
		return &analytics.InsightsResponse{Data: insights}, nil
	}

	// 1. Positive feedback - good streaks
	bestStreak := 0
	bestStreakHabit := ""
	for _, h := range activeHabits {
		if h.CurrentStreak > bestStreak {
			bestStreak = h.CurrentStreak
			bestStreakHabit = h.Name
		}
	}
	if bestStreak >= 7 {
		insights = append(insights, analytics.Insight{
			Type:        "positive",
			Title:       "Great Streak!",
			Description: bestStreakHabit + " has a " + fmt.Sprintf("%d", bestStreak) + "-day streak. Keep it up!",
			Priority:    "high",
		})
	}

	// 2. Suggestions - declining habits (low completion rate or broken streaks)
	now := lib.NormalizeDate(time.Now().UTC())
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	
	for _, h := range activeHabits {
		// Check recent completion rate
		logs, err := s.habitLogRepo.GetByHabit(ctx, userID, h.ID, thirtyDaysAgo, now)
		if err != nil {
			continue
		}

		recentCompletions := 0
		for _, log := range logs {
			if log.Completed {
				recentCompletions++
			}
		}

		// Calculate expected completions (rough estimate)
		expectedCompletions := 30 // for daily habits
		if h.Frequency == habit.Weekly {
			expectedCompletions = 4 // roughly 4 weeks
		}

		completionRate := 0.0
		if expectedCompletions > 0 {
			completionRate = (float64(recentCompletions) / float64(expectedCompletions)) * 100
		}

		// Suggest if completion rate is low or streak is broken
		if completionRate < 50 && h.CurrentStreak == 0 {
			insights = append(insights, analytics.Insight{
				Type:        "suggestion",
				Title:       "Get Back on Track",
				Description: h.Name + " has been inactive. Try to complete it today!",
				Priority:    "medium",
			})
		}
	}

	// 3. Best day analysis (from day of week data)
	dayOfWeekData, err := s.GetDayOfWeekAnalysis(ctx, userID, "30d")
	if err == nil && len(dayOfWeekData.Data) > 0 {
		bestDay := dayOfWeekData.Data[0]
		for _, day := range dayOfWeekData.Data {
			if day.CompletionRate > bestDay.CompletionRate {
				bestDay = day
			}
		}
		if bestDay.CompletionRate > 70 {
			insights = append(insights, analytics.Insight{
				Type:        "positive",
				Title:       "Best Day",
				Description: bestDay.Day + " is your most productive day with " + fmt.Sprintf("%.0f", bestDay.CompletionRate) + "% completion rate!",
				Priority:    "low",
			})
		}
	}

	// 4. Overall achievement progress
	metrics, err := s.GetMetrics(ctx, userID)
	if err == nil {
		if metrics.AvgCompletionRate >= 80 {
			insights = append(insights, analytics.Insight{
				Type:        "achievement",
				Title:       "Excellent Consistency",
				Description: "You're maintaining an 80%+ average completion rate across all habits!",
				Priority:    "high",
			})
		} else if metrics.AvgCompletionRate < 50 {
			insights = append(insights, analytics.Insight{
				Type:        "suggestion",
				Title:       "Room for Improvement",
				Description: "Your average completion rate is below 50%. Focus on consistency!",
				Priority:    "medium",
			})
		}
	}

	return &analytics.InsightsResponse{
		Data: insights,
	}, nil
}

