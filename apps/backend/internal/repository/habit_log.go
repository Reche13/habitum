package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reche13/habitum/internal/model/habitlog"
)

type HabitLogRepository struct {
	db *pgxpool.Pool
}

func NewHabitLogRepository(db *pgxpool.Pool) * HabitLogRepository {
	return &HabitLogRepository{db:db}
}

func (r *HabitLogRepository) Create(
	ctx context.Context,
	userID uuid.UUID,
	payload *habitlog.HabitLogPayload,
) (*habitlog.HabitLog, error) {
	stmt := `
		INSERT INTO 
		habit_logs (user_id, habit_id, log_date, completed) 
		VALUES (@user_id, @habit_id, @log_date, @completed) 
		ON CONFLICT (habit_id, log_date) 
		DO UPDATE SET
		completed = EXCLUDED.completed, updated_at = NOW()
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
		"habit_id": payload.HabitID,
		"log_date": payload.LogDate,
		"completed": payload.Completed,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hl, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[habitlog.HabitLog])
	if err != nil {
		return nil, err
	}

	return &hl, nil
}


func (r *HabitLogRepository) GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]habitlog.HabitLog, error) {
	stmt := `
		SELECT * 
		FROM habit_logs 
		WHERE user_id = @user_id 
		AND log_date = @log_date 
		ORDER BY created_at
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
		"log_date": date,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hls, err := pgx.CollectRows(rows, pgx.RowToStructByName[habitlog.HabitLog])
	if err != nil {
		return nil, err
	}

	return hls, nil
}


func (r *HabitLogRepository) GetByDateRange(
	ctx context.Context,
	userID uuid.UUID,
	from time.Time,
	to time.Time,
) ([]habitlog.HabitLog, error) {
	stmt := `
		SELECT *
		FROM habit_logs
		WHERE user_id = @user_id
		AND log_date BETWEEN @from AND @to
		ORDER BY log_date
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
		"from": from,
		"to": to,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[habitlog.HabitLog])
}


func (r *HabitLogRepository) GetByHabit(
	ctx context.Context,
	userID uuid.UUID,
	habitID uuid.UUID,
	from time.Time,
	to time.Time,
) ([]habitlog.HabitLog, error) {
	stmt := `
		SELECT *
		FROM habit_logs
		WHERE user_id = @user_id
		AND habit_id = @habit_id
		AND log_date BETWEEN @from AND @to
		ORDER BY log_date
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":userID,
		"habit_id":habitID,
		"from":from,
		"to":to,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[habitlog.HabitLog])
}
