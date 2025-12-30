package analytics

// CompletionTrendResponse represents daily completion data
type CompletionTrendResponse struct {
	Data []CompletionTrendDataPoint `json:"data"`
}

// CompletionTrendDataPoint represents completion data for a single day
type CompletionTrendDataPoint struct {
	Date           string  `json:"date"`            // Format: "yyyy-MM-dd"
	Completions    int     `json:"completions"`     // Number of completions on this day
	TotalHabits    int     `json:"totalHabits"`     // Total active habits on this day
	CompletionRate float64 `json:"completionRate"`  // Percentage (0-100)
}

// CategoryBreakdownResponse represents completion stats by category
type CategoryBreakdownResponse struct {
	Data []CategoryBreakdownDataPoint `json:"data"`
}

// CategoryBreakdownDataPoint represents stats for a single category
type CategoryBreakdownDataPoint struct {
	Category       string  `json:"category"`
	Label          string  `json:"label"`
	Icon           string  `json:"icon,omitempty"`
	HabitCount     int     `json:"habitCount"`
	AvgCompletionRate float64 `json:"avgCompletionRate"` // Percentage (0-100)
	TotalCompletions int     `json:"totalCompletions"`
}

// DayOfWeekResponse represents completion stats by day of week
type DayOfWeekResponse struct {
	Data []DayOfWeekDataPoint `json:"data"`
}

// DayOfWeekDataPoint represents stats for a single day of week
type DayOfWeekDataPoint struct {
	Day            string  `json:"day"`            // "Monday", "Tuesday", etc.
	DayIndex       int     `json:"dayIndex"`       // 0=Monday, 6=Sunday
	Completions    int     `json:"completions"`
	TotalHabits    int     `json:"totalHabits"`
	CompletionRate float64 `json:"completionRate"` // Percentage (0-100)
}

// MetricsResponse represents overall analytics metrics
type MetricsResponse struct {
	AvgCompletionRate float64 `json:"avgCompletionRate"` // Average across all habits
	AvgStreak         float64 `json:"avgStreak"`         // Average current streak
	TotalCompletions  int     `json:"totalCompletions"`   // Total all-time completions
	ConsistencyScore  float64 `json:"consistencyScore"`  // Percentage (0-100)
}

// TopHabitsResponse represents top performing habits
type TopHabitsResponse struct {
	Data []TopHabitDataPoint `json:"data"`
}

// TopHabitDataPoint represents a single top habit
type TopHabitDataPoint struct {
	HabitID        string  `json:"habitId"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	CompletionRate float64 `json:"completionRate"` // Percentage (0-100)
	CurrentStreak  int     `json:"currentStreak"`
	LongestStreak  int     `json:"longestStreak"`
}

// StreakLeaderboardResponse represents habits sorted by streak
type StreakLeaderboardResponse struct {
	Data []StreakLeaderboardDataPoint `json:"data"`
}

// StreakLeaderboardDataPoint represents a single habit in leaderboard
type StreakLeaderboardDataPoint struct {
	HabitID       string `json:"habitId"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	CurrentStreak int    `json:"currentStreak"`
	LongestStreak int    `json:"longestStreak"`
}

// InsightsResponse represents personalized insights
type InsightsResponse struct {
	Data []Insight `json:"data"`
}

// Insight represents a single insight
type Insight struct {
	Type        string `json:"type"`        // "positive", "suggestion", "achievement"
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`    // "high", "medium", "low"
}

