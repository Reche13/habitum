package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reche13/habitum/internal/model"
	"github.com/reche13/habitum/internal/model/user"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(
	ctx context.Context,
	payload *user.CreateUserPayload,
) (*user.User, error) {
	stmt := `
		INSERT INTO 
			users (
				name,
				email
			)
		VALUES 
			(
				@name,
				@email
			)
		RETURNING
			*
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"name":  payload.Name,
		"email": payload.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create user query for name=%s email=%s: %w", payload.Name, payload.Email, err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for name=%s email=%s: %w", payload.Name, payload.Email, err)
	}

	return &user, nil
}


func (r *UserRepository) List(ctx context.Context) (*model.PaginatedResponse[user.User], error) {
	stmt := `
		SELECT
			*
		FROM 
			users
		ORDER BY 
			created_at 
		DESC
	`

	rows, err := r.db.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get users query: %w", err)
	}
	
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[user.User]{
				Data:       []user.User{},
				Page:       0,
				Limit:      0,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:users: %w", err)
	}

	return &model.PaginatedResponse[user.User]{
		Data:       users,
		Page:       0,
		Limit:      0,
		Total:      len(users),
		TotalPages: 0,
	}, nil
}
