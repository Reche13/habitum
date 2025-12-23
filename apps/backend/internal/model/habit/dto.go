package habit

type CreateHabitPayload struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description"`
	Icon *string `json:"icon"`
	Color *string `json:"color"`
	Category Category `json:"category" validate:"required"`
	Frequency Frequency `json:"frequency" validate:"required"`
	TimesPerWeek *int `json:"times_per_week,omitempty"`
}