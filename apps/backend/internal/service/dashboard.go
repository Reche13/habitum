package service

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/lib"
	"github.com/reche13/habitum/internal/model/dashboard"
	"github.com/reche13/habitum/internal/model/habit"
	"github.com/reche13/habitum/internal/repository"
)

type DashboardService struct {
	*BaseService
	habitRepo    *repository.HabitRepository
	habitLogRepo *repository.HabitLogRepository
}

func NewDashboardService(
	habitRepo *repository.HabitRepository,
	habitLogRepo *repository.HabitLogRepository,
) *DashboardService {
	return &DashboardService{
		BaseService: &BaseService{
			resourceName: "dashboard",
		},
		habitRepo:    habitRepo,
		habitLogRepo: habitLogRepo,
	}
}

func (s *DashboardService) GetHome(ctx context.Context, userID uuid.UUID) (*dashboard.DashboardResponse, error) {
	// Get all active habits
	allHabits, _, err := s.habitRepo.List(ctx, userID, nil)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Filter out archived habits
	activeHabits := make([]habit.Habit, 0)
	for _, h := range allHabits {
		if h.ArchivedAt == nil {
			activeHabits = append(activeHabits, h)
		}
	}

	if len(activeHabits) == 0 {
		today := lib.NormalizeDate(time.Now().UTC())
		return &dashboard.DashboardResponse{
			Today: dashboard.TodayStats{
				Date:          today.Format("2006-01-02"),
				CompletionRate: 0,
				CompletedCount: 0,
				TotalCount:     0,
			},
			HabitsToComplete: []dashboard.HabitSummary{},
			HabitsCompleted:  []dashboard.HabitSummary{},
			ActiveStreaks:    []dashboard.StreakSummary{},
			QuickStats: dashboard.QuickStats{
				TodayRate:     0,
				ThisWeek:      0,
				LongestStreak: 0,
				TotalHabits:   0,
			},
			Achievements: []dashboard.AchievementSummary{},
		}, nil
	}

	today := lib.NormalizeDate(time.Now().UTC())
	todayStr := today.Format("2006-01-02")

	// Get today's logs
	todayLogs, err := s.habitLogRepo.GetByDate(ctx, userID, today)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Create map of completed habits today
	completedTodayMap := make(map[uuid.UUID]bool)
	completedTodayAtMap := make(map[uuid.UUID]time.Time)
	for _, log := range todayLogs {
		if log.Completed {
			completedTodayMap[log.HabitID] = true
			completedTodayAtMap[log.HabitID] = log.CreatedAt
		}
	}

	// Get this week's logs for quick stats
	now := lib.NormalizeDate(time.Now().UTC())
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekLogs, err := s.habitLogRepo.GetByDateRange(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, s.wrapError(err)
	}

	// Count completions this week
	completionsThisWeek := 0
	for _, log := range weekLogs {
		if log.Completed {
			completionsThisWeek++
		}
	}

	// Build habit summaries
	habitsToComplete := make([]dashboard.HabitSummary, 0)
	habitsCompleted := make([]dashboard.HabitSummary, 0)
	activeStreaks := make([]dashboard.StreakSummary, 0)

	totalCompleted := 0
	longestStreak := 0

	for _, h := range activeHabits {
		completedToday := completedTodayMap[h.ID]
		var completedTodayAt *time.Time
		if completedToday {
			completedAt := completedTodayAtMap[h.ID]
			completedTodayAt = &completedAt
			totalCompleted++
		}

		// Get current streak (use stored value for now, could calculate if needed)
		currentStreak := h.CurrentStreak
		if h.LongestStreak > longestStreak {
			longestStreak = h.LongestStreak
		}

		habitSummary := dashboard.HabitSummary{
			ID:             h.ID.String(),
			Name:           h.Name,
			Description:    h.Description,
			Icon:           h.Icon,
			IconID:         lib.GetStringValue(h.Icon),
			Color:          h.Color,
			Frequency:      string(h.Frequency),
			Category:       string(h.Category),
			CurrentStreak:  currentStreak,
			CompletedToday: completedToday,
			CompletedTodayAt: completedTodayAt,
		}

		if completedToday {
			habitsCompleted = append(habitsCompleted, habitSummary)
		} else {
			habitsToComplete = append(habitsToComplete, habitSummary)
		}

		// Add to active streaks if streak > 0
		if currentStreak > 0 {
			activeStreaks = append(activeStreaks, dashboard.StreakSummary{
				ID:            h.ID.String(),
				Name:          h.Name,
				Icon:          h.Icon,
				Color:         h.Color,
				CurrentStreak: currentStreak,
				LongestStreak: h.LongestStreak,
			})
		}
	}

	// Sort active streaks by current streak (descending)
	sort.Slice(activeStreaks, func(i, j int) bool {
		return activeStreaks[i].CurrentStreak > activeStreaks[j].CurrentStreak
	})

	// Limit to top 5
	if len(activeStreaks) > 5 {
		activeStreaks = activeStreaks[:5]
	}

	// Calculate today's completion rate
	totalCount := len(activeHabits)
	completionRate := 0.0
	if totalCount > 0 {
		completionRate = (float64(totalCompleted) / float64(totalCount)) * 100
	}

	// Build response
	return &dashboard.DashboardResponse{
		Today: dashboard.TodayStats{
			Date:          todayStr,
			CompletionRate: completionRate,
			CompletedCount: totalCompleted,
			TotalCount:     totalCount,
		},
		HabitsToComplete: habitsToComplete,
		HabitsCompleted:  habitsCompleted,
		ActiveStreaks:    activeStreaks,
		QuickStats: dashboard.QuickStats{
			TodayRate:     completionRate,
			ThisWeek:      completionsThisWeek,
			LongestStreak: longestStreak,
			TotalHabits:   totalCount,
		},
		Achievements: []dashboard.AchievementSummary{}, // Empty for now, will be populated when achievements are implemented
	}, nil
}


