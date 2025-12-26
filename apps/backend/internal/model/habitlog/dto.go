package habitlog

import (
	"time"

	"github.com/google/uuid"
)

type HabitLogPayload struct {
	HabitID uuid.UUID `json:"habit_id" db:"habit_id"`
	LogDate time.Time `json:"log_date" db:"log_date"`
	Completed bool `json:"completed" db:"completed"`
}