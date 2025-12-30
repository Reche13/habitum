package calendar

// CompletionsResponse represents completions for a date range
type CompletionsResponse struct {
	Completions []CompletionDay `json:"completions"`
	Statistics  PeriodStats     `json:"statistics"`
}

// CompletionDay represents completions for a single day
type CompletionDay struct {
	Date            string        `json:"date"`            // Format: "yyyy-MM-dd"
	Habits          []HabitInfo   `json:"habits"`          // Completed habits on this day
	CompletionRate  float64       `json:"completionRate"`  // Percentage (0-100)
	TotalHabits     int           `json:"totalHabits"`
	CompletedHabits int           `json:"completedHabits"`
}

// HabitInfo represents basic habit info for calendar
type HabitInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

// PeriodStats represents statistics for a period
type PeriodStats struct {
	TotalCompletions    int     `json:"totalCompletions"`
	DaysWithCompletions int     `json:"daysWithCompletions"`
	CompletionRate      float64 `json:"completionRate"` // Percentage (0-100)
	TotalDays           int     `json:"totalDays"`
}

// MonthResponse represents month view data
type MonthResponse struct {
	Year       int           `json:"year"`
	Month      int           `json:"month"`
	Days       []DayData     `json:"days"`
	Statistics PeriodStats   `json:"statistics"`
}

// WeekResponse represents week view data
type WeekResponse struct {
	Year       int           `json:"year"`
	Week       int           `json:"week"`
	Days       []DayData     `json:"days"`
	Statistics PeriodStats   `json:"statistics"`
}

// YearResponse represents year view (heatmap) data
type YearResponse struct {
	Year       int           `json:"year"`
	Heatmap    []HeatmapDay  `json:"heatmap"`
	Statistics PeriodStats   `json:"statistics"`
}

// DayData represents a single day in month/week view
type DayData struct {
	Date           string   `json:"date"`           // Format: "yyyy-MM-dd"
	DayOfWeek      int      `json:"dayOfWeek,omitempty"` // 0=Monday, 6=Sunday
	Completions    []string `json:"completions"`    // Array of habit IDs
	CompletionRate float64  `json:"completionRate"` // Percentage (0-100)
}

// HeatmapDay represents a day in year heatmap
type HeatmapDay struct {
	Date           string  `json:"date"`           // Format: "yyyy-MM-dd"
	CompletionRate float64 `json:"completionRate"` // Percentage (0-100)
	Intensity      int     `json:"intensity"`      // 0-4 (for heatmap visualization)
}

