package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		return nil, err
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM 
			users
		WHERE
			id = @id
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]user.User, error) {
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
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}

	if users == nil {
		return []user.User{}, nil
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, payload *user.UpdateUserPayload) (*user.User, error) {
	// Build dynamic update query
	updates := []string{}
	args := pgx.NamedArgs{"id": id}

	if payload.Name != nil {
		updates = append(updates, "name = @name")
		args["name"] = *payload.Name
	}

	if payload.Email != nil {
		updates = append(updates, "email = @email")
		args["email"] = *payload.Email
	}

	if len(updates) == 0 {
		// No updates, just return the user
		return r.GetByID(ctx, id)
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

	stmt := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = @id
		RETURNING *
	`, updateClause)

	rows, err := r.db.Query(ctx, stmt, args)
	if err != nil {
		return nil, err
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt := `
		DELETE FROM users
		WHERE id = @id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return err
	}

	return nil
}
