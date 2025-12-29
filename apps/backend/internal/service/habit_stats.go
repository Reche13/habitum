package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/lib"
	"github.com/reche13/habitum/internal/model/habit"
	"github.com/reche13/habitum/internal/repository"
)

// HabitStats contains computed statistics for a habit
type HabitStats struct {
	CurrentStreak      int
	LongestStreak      int
	CompletionRate     float64
	CompletedToday     bool
	CompletedTodayAt   *time.Time
	CompletedThisWeek  int
	CompletionHistory  []string
}

// CalculateCurrentStreak calculates the current consecutive streak for a habit
func CalculateCurrentStreak(
	ctx context.Context,
	habitLogRepo *repository.HabitLogRepository,
	userID uuid.UUID,
	habitID uuid.UUID,
	frequency habit.Frequency,
) (int, error) {
	// Get all completed logs for this habit
	startDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := lib.NormalizeDate(time.Now().UTC())

	logs, err := habitLogRepo.GetByHabit(ctx, userID, habitID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	// Extract completed dates
	completedDates := make([]time.Time, 0)
	for _, log := range logs {
		if log.Completed {
			completedDates = append(completedDates, log.LogDate)
		}
	}

	if len(completedDates) == 0 {
		return 0, nil
	}

	// Sort dates descending (most recent first)
	sort.Slice(completedDates, func(i, j int) bool {
		return completedDates[i].After(completedDates[j])
	})

	// Calculate streak based on frequency
	if frequency == habit.Daily {
		return calculateDailyStreak(completedDates), nil
	} else {
		return calculateWeeklyStreak(completedDates), nil
	}
}

// calculateDailyStreak calculates streak for daily habits
func calculateDailyStreak(completedDates []time.Time) int {
	if len(completedDates) == 0 {
		return 0
	}

	today := lib.NormalizeDate(time.Now().UTC())
	streak := 0

	// Check if today is completed
	expectedDate := today
	if !isDateInSlice(completedDates, expectedDate) {
		// If today is not completed, start from yesterday
		expectedDate = today.AddDate(0, 0, -1)
	}

	// Count consecutive days backwards
	for i := 0; i < len(completedDates); i++ {
		if completedDates[i].Equal(expectedDate) {
			streak++
			expectedDate = expectedDate.AddDate(0, 0, -1)
		} else if completedDates[i].Before(expectedDate) {
			// Gap found, streak broken
			break
		}
	}

	return streak
}

// calculateWeeklyStreak calculates streak for weekly habits
func calculateWeeklyStreak(completedDates []time.Time) int {
	if len(completedDates) == 0 {
		return 0
	}

	// Group by week (ISO week)
	weekMap := make(map[string]bool)
	for _, date := range completedDates {
		year, week := date.ISOWeek()
		weekKey := formatWeekKey(year, week)
		weekMap[weekKey] = true
	}

	// Get all weeks and sort
	weeks := make([]string, 0, len(weekMap))
	for week := range weekMap {
		weeks = append(weeks, week)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(weeks)))

	if len(weeks) == 0 {
		return 0
	}

	// Calculate consecutive weeks
	streak := 0
	currentYear, currentWeek := time.Now().UTC().ISOWeek()
	currentWeekKey := formatWeekKey(currentYear, currentWeek)

	// Check if current week is completed
	expectedWeek := currentWeekKey
	if !contains(weeks, expectedWeek) {
		// If current week is not completed, go to previous week
		prevWeek := time.Now().UTC().AddDate(0, 0, -7)
		prevYear, prevWeekNum := prevWeek.ISOWeek()
		expectedWeek = formatWeekKey(prevYear, prevWeekNum)
	}

	// Count consecutive weeks backwards
	streak = 0
	if len(weeks) > 0 {
		// Start from current/previous week and count backwards
		checkDate := time.Now().UTC()
		if !contains(weeks, currentWeekKey) {
			checkDate = checkDate.AddDate(0, 0, -7)
		}
		
		for i := 0; i < 1000; i++ { // Max 1000 weeks (~19 years)
			year, week := checkDate.ISOWeek()
			weekKey := formatWeekKey(year, week)
			if contains(weeks, weekKey) {
				streak++
				checkDate = checkDate.AddDate(0, 0, -7)
			} else {
				break
			}
		}
	}

	return streak
}

// CalculateLongestStreak calculates the longest streak ever achieved
func CalculateLongestStreak(
	ctx context.Context,
	habitLogRepo *repository.HabitLogRepository,
	userID uuid.UUID,
	habitID uuid.UUID,
	frequency habit.Frequency,
) (int, error) {
	// Get all completed logs
	startDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := lib.NormalizeDate(time.Now().UTC())

	logs, err := habitLogRepo.GetByHabit(ctx, userID, habitID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	// Extract completed dates
	completedDates := make([]time.Time, 0)
	for _, log := range logs {
		if log.Completed {
			completedDates = append(completedDates, log.LogDate)
		}
	}

	if len(completedDates) == 0 {
		return 0, nil
	}

	// Sort dates ascending
	sort.Slice(completedDates, func(i, j int) bool {
		return completedDates[i].Before(completedDates[j])
	})

	if frequency == habit.Daily {
		return findLongestDailyStreak(completedDates), nil
	} else {
		return findLongestWeeklyStreak(completedDates), nil
	}
}

// findLongestDailyStreak finds the longest consecutive daily streak
func findLongestDailyStreak(completedDates []time.Time) int {
	if len(completedDates) == 0 {
		return 0
	}

	longest := 1
	current := 1

	for i := 1; i < len(completedDates); i++ {
		daysDiff := int(completedDates[i].Sub(completedDates[i-1]).Hours() / 24)
		if daysDiff == 1 {
			// Consecutive day
			current++
			if current > longest {
				longest = current
			}
		} else {
			// Gap found, reset
			current = 1
		}
	}

	return longest
}

// findLongestWeeklyStreak finds the longest consecutive weekly streak
func findLongestWeeklyStreak(completedDates []time.Time) int {
	if len(completedDates) == 0 {
		return 0
	}

	// Group by week
	weekMap := make(map[string]bool)
	for _, date := range completedDates {
		year, week := date.ISOWeek()
		weekKey := formatWeekKey(year, week)
		weekMap[weekKey] = true
	}

	// Get all weeks and sort
	weeks := make([]string, 0, len(weekMap))
	for week := range weekMap {
		weeks = append(weeks, week)
	}
	sort.Strings(weeks)

	if len(weeks) == 0 {
		return 0
	}

	longest := 1
	current := 1

	// For weekly streaks, use a simpler approach: convert weeks to dates and check consecutive
	// Get dates for each week (use first day of week)
	weekDates := make([]time.Time, 0, len(weeks))
	for _, weekStr := range weeks {
		year, weekNum := parseWeekKey(weekStr)
		if year == 0 {
			continue
		}
		// Find first day of ISO week for that year/week
		date := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC) // Jan 4 is always in week 1
		dateYear, dateWeek := date.ISOWeek()
		for dateYear != year || dateWeek != weekNum {
			if dateWeek < weekNum {
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
		weekDates = append(weekDates, date)
	}
	
	// Sort by date
	sort.Slice(weekDates, func(i, j int) bool {
		return weekDates[i].Before(weekDates[j])
	})

	// Check consecutive weeks
	for i := 1; i < len(weekDates); i++ {
		daysDiff := int(weekDates[i].Sub(weekDates[i-1]).Hours() / (24 * 7))
		if daysDiff == 1 {
			current++
			if current > longest {
				longest = current
			}
		} else {
			current = 1
		}
	}

	return longest
}

// CalculateCompletionRate calculates the completion rate for a habit
func CalculateCompletionRate(
	ctx context.Context,
	habitLogRepo *repository.HabitLogRepository,
	userID uuid.UUID,
	habitID uuid.UUID,
	habitCreatedAt time.Time,
	frequency habit.Frequency,
) (float64, error) {
	now := lib.NormalizeDate(time.Now().UTC())
	createdAt := lib.NormalizeDate(habitCreatedAt)

	// Get all completed logs
	startDate := createdAt
	endDate := now

	logs, err := habitLogRepo.GetByHabit(ctx, userID, habitID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	// Count completed logs
	completedCount := 0
	for _, log := range logs {
		if log.Completed {
			completedCount++
		}
	}

	if frequency == habit.Daily {
		// For daily habits: completed days / total days since creation
		totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
		if totalDays == 0 {
			return 0, nil
		}
		return (float64(completedCount) / float64(totalDays)) * 100, nil
	} else {
		// For weekly habits: completed weeks / total weeks since creation
		// Calculate weeks between dates
		weeksDiff := int(endDate.Sub(startDate).Hours()/(24*7)) + 1
		if weeksDiff == 0 {
			return 0, nil
		}
		return (float64(completedCount) / float64(weeksDiff)) * 100, nil
	}
}

// GetTodayStatus checks if habit is completed today
func GetTodayStatus(
	ctx context.Context,
	habitLogRepo *repository.HabitLogRepository,
	userID uuid.UUID,
	habitID uuid.UUID,
) (bool, *time.Time, error) {
	today := lib.NormalizeDate(time.Now().UTC())
	logs, err := habitLogRepo.GetByDate(ctx, userID, today)
	if err != nil {
		return false, nil, err
	}

	for _, log := range logs {
		if log.HabitID == habitID && log.Completed {
			return true, &log.CreatedAt, nil
		}
	}

	return false, nil, nil
}

// GetCompletedThisWeek counts completions this week
func GetCompletedThisWeek(
	ctx context.Context,
	habitLogRepo *repository.HabitLogRepository,
	userID uuid.UUID,
	habitID uuid.UUID,
) (int, error) {
	now := lib.NormalizeDate(time.Now().UTC())
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekEnd := weekStart.AddDate(0, 0, 6)

	logs, err := habitLogRepo.GetByHabit(ctx, userID, habitID, weekStart, weekEnd)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, log := range logs {
		if log.Completed {
			count++
		}
	}

	return count, nil
}

// GetCompletionHistoryDates gets array of completion dates
func GetCompletionHistoryDates(
	ctx context.Context,
	habitLogRepo *repository.HabitLogRepository,
	userID uuid.UUID,
	habitID uuid.UUID,
	limit int,
) ([]string, error) {
	startDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := lib.NormalizeDate(time.Now().UTC())

	logs, err := habitLogRepo.GetByHabit(ctx, userID, habitID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	dates := make([]string, 0)
	for _, log := range logs {
		if log.Completed {
			dates = append(dates, log.LogDate.Format("2006-01-02"))
		}
	}

	// Sort descending (most recent first) and limit
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	if limit > 0 && len(dates) > limit {
		dates = dates[:limit]
	}

	return dates, nil
}

// Helper functions

func isDateInSlice(dates []time.Time, target time.Time) bool {
	for _, d := range dates {
		if d.Equal(target) {
			return true
		}
	}
	return false
}

func formatWeekKey(year, week int) string {
	// Format as "YYYY-WW" for sorting (zero-padded week)
	return fmt.Sprintf("%04d-W%02d", year, week)
}

func parseWeekKey(key string) (int, int) {
	// Parse "YYYY-WW" format
	var year, week int
	_, err := fmt.Sscanf(key, "%04d-W%02d", &year, &week)
	if err != nil {
		return 0, 0
	}
	return year, week
}

func getLastWeekOfYear(year int) int {
	// Get last week of year
	lastDay := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	_, week := lastDay.ISOWeek()
	return week
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
