package habit

import "time"

// HabitResponse includes the habit with computed fields
type HabitResponse struct {
	Habit
	CompletionRate     float64    `json:"completionRate"`
	CompletedToday     bool       `json:"completedToday"`
	CompletedTodayAt   *time.Time `json:"completedTodayAt,omitempty"`
	CompletedThisWeek  int        `json:"completedThisWeek"`
	CompletionHistory  []string   `json:"completionHistory,omitempty"`
}

