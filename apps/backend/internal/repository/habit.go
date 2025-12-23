package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reche13/habitum/internal/model/habit"
)


type HabitRepository struct {
	db *pgxpool.Pool
}

func NewHabitRepository(db *pgxpool.Pool) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) Create(ctx context.Context, userID uuid.UUID, payload *habit.CreateHabitPayload) (*habit.Habit, error) {
	stmt := `
		INSERT INTO habits (
			user_id, name, description, icon, color,
			category, frequency, times_per_week
		)
		VALUES (
			@user_id, @name, @description, @icon, @color,
			@category, @frequency, @times_per_week
		)
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":      userID,
		"name":         payload.Name,
		"description":  payload.Description,
		"icon":         payload.Icon,
		"color":        payload.Color,
		"category":     payload.Category,
		"frequency":    payload.Frequency,
		"times_per_week": payload.TimesPerWeek,
	})

	if err != nil {
		return nil, err
	}

	h, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[habit.Habit])
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (r *HabitRepository) List(ctx context.Context, userID uuid.UUID) ([]habit.Habit, error) {
	stmt := `
		SELECT
			*
		FROM 
			habits
		WHERE
			user_id = @user_id
		ORDER BY 
			created_at 
		DESC
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"user_id":userID})
	if err != nil {
		return nil, err
	}

	habits, err := pgx.CollectRows(rows, pgx.RowToStructByName[habit.Habit])
	if err != nil {
		return nil, err
	}

	if habits == nil {
		return []habit.Habit{}, nil
	}

	return habits, nil
}