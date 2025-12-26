package habitlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/reche13/habitum/internal/model"
)

type HabitLog struct {
	model.Base

	UserID uuid.UUID `json:"user_id" db:"user_id"`
	HabitID uuid.UUID `json:"habit_id" db:"habit_id"`
	LogDate time.Time `json:"log_date" db:"log_date"`
	Completed bool `json:"completed" db:"completed"`
}