package repository

import (
	"context"
	"fmt"
	"strings"

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

func (r *HabitRepository) List(ctx context.Context, userID uuid.UUID, filters *habit.ListFilters) ([]habit.Habit, int, error) {
	// Build WHERE clause
	whereConditions := []string{"user_id = @user_id"}
	args := pgx.NamedArgs{"user_id": userID}

	// Add category filter
	if filters != nil && filters.Category != nil {
		whereConditions = append(whereConditions, "category = @category")
		args["category"] = *filters.Category
	}

	// Add search filter (search in name and description)
	if filters != nil && filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + strings.ToLower(*filters.Search) + "%"
		whereConditions = append(whereConditions, "(LOWER(name) LIKE @search OR LOWER(description) LIKE @search)")
		args["search"] = searchTerm
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Build ORDER BY clause (whitelist approach for security)
	orderBy := "created_at DESC" // default
	orderDir := "DESC"            // default direction
	
	if filters != nil && filters.Order != nil {
		orderDir = strings.ToUpper(*filters.Order)
		if orderDir != "ASC" && orderDir != "DESC" {
			orderDir = "DESC"
		}
	}

	if filters != nil && filters.Sort != nil {
		sortValue := strings.ToLower(*filters.Sort)
		switch sortValue {
		case "name":
			orderBy = fmt.Sprintf("name %s", orderDir)
		case "date":
			orderBy = fmt.Sprintf("created_at %s", orderDir)
		case "streak":
			// Note: This will need to be updated when streaks are added to the table
			orderBy = fmt.Sprintf("created_at %s", orderDir)
		case "completion":
			// Note: This will need to be updated when completion rates are calculated
			orderBy = fmt.Sprintf("created_at %s", orderDir)
		default:
			orderBy = fmt.Sprintf("created_at %s", orderDir)
		}
	} else {
		// If no sort specified, use default with order direction
		orderBy = fmt.Sprintf("created_at %s", orderDir)
	}

	// Build base query
	baseStmt := fmt.Sprintf(`
		SELECT
			*
		FROM 
			habits
		WHERE
			%s
		ORDER BY 
			%s
	`, whereClause, orderBy)

	// Get total count for pagination
	countStmt := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM habits
		WHERE %s
	`, whereClause)

	var total int
	err := r.db.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := 1
	limit := 50 // default limit
	if filters != nil {
		if filters.Page != nil && *filters.Page > 0 {
			page = *filters.Page
		}
		if filters.Limit != nil && *filters.Limit > 0 {
			limit = *filters.Limit
			// Cap limit at 100 to prevent abuse
			if limit > 100 {
				limit = 100
			}
		}
	}

	offset := (page - 1) * limit
	args["limit"] = limit
	args["offset"] = offset
	stmt := baseStmt + " LIMIT @limit OFFSET @offset"

	rows, err := r.db.Query(ctx, stmt, args)
	if err != nil {
		return nil, 0, err
	}

	habits, err := pgx.CollectRows(rows, pgx.RowToStructByName[habit.Habit])
	if err != nil {
		return nil, 0, err
	}

	if habits == nil {
		return []habit.Habit{}, total, nil
	}

	return habits, total, nil
}

func (r *HabitRepository) GetByID(ctx context.Context, habitID uuid.UUID, userID uuid.UUID) (*habit.Habit, error) {
	stmt := `
		SELECT
			*
		FROM 
			habits
		WHERE
			id = @habit_id
			AND user_id = @user_id
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"habit_id": habitID,
		"user_id":  userID,
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

func (r *HabitRepository) Update(ctx context.Context, habitID uuid.UUID, userID uuid.UUID, payload *habit.UpdateHabitPayload) (*habit.Habit, error) {
	// Build dynamic update query
	updates := []string{}
	args := pgx.NamedArgs{
		"habit_id": habitID,
		"user_id":  userID,
	}

	if payload.Name != nil {
		updates = append(updates, "name = @name")
		args["name"] = *payload.Name
	}

	if payload.Description != nil {
		updates = append(updates, "description = @description")
		args["description"] = *payload.Description
	}

	if payload.Icon != nil {
		updates = append(updates, "icon = @icon")
		args["icon"] = *payload.Icon
	}

	if payload.Color != nil {
		updates = append(updates, "color = @color")
		args["color"] = *payload.Color
	}

	if payload.Category != nil {
		updates = append(updates, "category = @category")
		args["category"] = *payload.Category
	}

	if payload.Frequency != nil {
		updates = append(updates, "frequency = @frequency")
		args["frequency"] = *payload.Frequency
	}

	if payload.TimesPerWeek != nil {
		updates = append(updates, "times_per_week = @times_per_week")
		args["times_per_week"] = *payload.TimesPerWeek
	}

	if len(updates) == 0 {
		// No updates, just return the habit
		return r.GetByID(ctx, habitID, userID)
	}

	updates = append(updates, "updated_at = NOW()")

	// Build update clause
	updateClause := ""
	for i, update := range updates {
		if i > 0 {
			updateClause += ", "
		}
		updateClause += update
	}

	stmt := `
		UPDATE habits
		SET ` + updateClause + `
		WHERE id = @habit_id
			AND user_id = @user_id
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, args)
	if err != nil {
		return nil, err
	}

	h, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[habit.Habit])
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (r *HabitRepository) Delete(ctx context.Context, habitID uuid.UUID, userID uuid.UUID) error {
	stmt := `
		DELETE FROM habits
		WHERE id = @habit_id
			AND user_id = @user_id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"habit_id": habitID,
		"user_id":  userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *HabitRepository) UpdateStreaks(
	ctx context.Context,
	habitID uuid.UUID,
	userID uuid.UUID,
	currentStreak int,
	longestStreak int,
) error {
	stmt := `
		UPDATE habits
		SET current_streak = @current_streak,
			longest_streak = @longest_streak,
			updated_at = NOW()
		WHERE id = @habit_id
			AND user_id = @user_id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"habit_id":      habitID,
		"user_id":       userID,
		"current_streak": currentStreak,
		"longest_streak": longestStreak,
	})
	if err != nil {
		return err
	}

	return nil
}