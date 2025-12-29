package habit

type ListFilters struct {
	Category *Category `json:"category,omitempty"`
	Search   *string   `json:"search,omitempty"`
	Sort     *string   `json:"sort,omitempty"` // "name", "date", "streak", "completion"
	Order    *string   `json:"order,omitempty"` // "asc" or "desc"
	Page     *int      `json:"page,omitempty"`
	Limit    *int      `json:"limit,omitempty"`
}

type SortOption string

const (
	SortByName        SortOption = "name"
	SortByDate        SortOption = "date"
	SortByStreak      SortOption = "streak"
	SortByCompletion   SortOption = "completion"
)

