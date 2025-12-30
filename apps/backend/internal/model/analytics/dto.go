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

