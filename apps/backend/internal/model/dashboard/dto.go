package dashboard

import "time"

// DashboardResponse represents the dashboard home data
type DashboardResponse struct {
	Today           TodayStats        `json:"today"`
	HabitsToComplete []HabitSummary   `json:"habitsToComplete"`
	HabitsCompleted  []HabitSummary   `json:"habitsCompleted"`
	ActiveStreaks    []StreakSummary  `json:"activeStreaks"`
	QuickStats       QuickStats       `json:"quickStats"`
	Achievements     []AchievementSummary `json:"recentAchievements,omitempty"`
}

// TodayStats represents today's completion statistics
type TodayStats struct {
	Date          string  `json:"date"`          // Format: "yyyy-MM-dd"
	CompletionRate float64 `json:"completionRate"` // Percentage (0-100)
	CompletedCount int     `json:"completedCount"`
	TotalCount     int     `json:"totalCount"`
}

// HabitSummary represents a habit for dashboard display
type HabitSummary struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Icon            *string    `json:"icon,omitempty"`
	IconID          string     `json:"iconId,omitempty"`
	Color           *string    `json:"color,omitempty"`
	Frequency       string     `json:"frequency"`
	Category        string     `json:"category"`
	CurrentStreak   int        `json:"currentStreak"`
	CompletedToday  bool       `json:"completedToday"`
	CompletedTodayAt *time.Time `json:"completedTodayAt,omitempty"`
}

// StreakSummary represents a habit with active streak
type StreakSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Icon          *string `json:"icon,omitempty"`
	Color         *string `json:"color,omitempty"`
	CurrentStreak int    `json:"currentStreak"`
	LongestStreak int    `json:"longestStreak"`
}

// QuickStats represents quick statistics
type QuickStats struct {
	TodayRate      float64 `json:"todayRate"`      // Percentage (0-100)
	ThisWeek       int     `json:"thisWeek"`       // Total completions this week
	LongestStreak  int     `json:"longestStreak"`  // Longest streak across all habits
	TotalHabits    int     `json:"totalHabits"`
}

// AchievementSummary represents an achievement (placeholder for future)
type AchievementSummary struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	UnlockedAt  *time.Time `json:"unlockedAt,omitempty"`
	Progress    *int       `json:"progress,omitempty"`
	Target      *int       `json:"target,omitempty"`
}

