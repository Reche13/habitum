package habit

import (
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/model"
)

type Frequency string

const (
	Daily Frequency = "daily"
	Weekly Frequency = "weekly"
)

type Category string

const (
	Health Category = "health"
	Productivity Category = "productivity"
	Learning Category = "learning"
	Work Category = "work"
	Personal Category = "personal"
	Mindfulness Category = "mindfulness"
	Social Category = "social"
	Creative Category = "creative"
	Finance Category = "finance"
	Other Category = "other"
)

type Habit struct {
	model.Base

	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	Icon *string `json:"icon" db:"icon"`
	Color *string `json:"color" db:"color"`
	Category Category `json:"category" db:"category"`
	Frequency Frequency `json:"frequency" db:"frequency"`
	TimesPerWeek *int `json:"times_per_week,omitempty" db:"times_per_week"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" db:"archived_at"`
}