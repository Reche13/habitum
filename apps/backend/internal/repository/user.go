package repository

import (
	"context"
	"fmt"
	"time"

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

// GetByEmail finds a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM 
			users
		WHERE
			email = @email
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
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

// GetByOAuthProvider finds a user by OAuth provider and provider ID
func (r *UserRepository) GetByOAuthProvider(ctx context.Context, provider, providerID string) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM 
			users
		WHERE
			oauth_provider = @provider
			AND oauth_provider_id = @provider_id
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"provider":    provider,
		"provider_id": providerID,
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

// CreateWithPassword creates a user with a password hash
func (r *UserRepository) CreateWithPassword(
	ctx context.Context,
	name, email, passwordHash string,
) (*user.User, error) {
	stmt := `
		INSERT INTO 
			users (
				name,
				email,
				password_hash
			)
		VALUES 
			(
				@name,
				@email,
				@password_hash
			)
		RETURNING
			*
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"name":          name,
		"email":         email,
		"password_hash": passwordHash,
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

// UpdateEmailVerification updates email verification status and token
func (r *UserRepository) UpdateEmailVerification(
	ctx context.Context,
	userID uuid.UUID,
	verified bool,
	token *string,
	expiresAt *time.Time,
) error {
	stmt := `
		UPDATE users
		SET 
			email_verified = @verified,
			email_verification_token = @token,
			email_verification_expires_at = @expires_at,
			updated_at = NOW()
		WHERE id = @id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id":         userID,
		"verified":   verified,
		"token":      token,
		"expires_at": expiresAt,
	})
	return err
}

// UpdatePassword updates user password hash
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	stmt := `
		UPDATE users
		SET 
			password_hash = @password_hash,
			updated_at = NOW()
		WHERE id = @id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id":           userID,
		"password_hash": passwordHash,
	})
	return err
}

// UpdatePasswordResetToken updates password reset token
func (r *UserRepository) UpdatePasswordResetToken(
	ctx context.Context,
	userID uuid.UUID,
	token *string,
	expiresAt *time.Time,
) error {
	stmt := `
		UPDATE users
		SET 
			password_reset_token = @token,
			password_reset_expires_at = @expires_at,
			updated_at = NOW()
		WHERE id = @id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id":         userID,
		"token":      token,
		"expires_at": expiresAt,
	})
	return err
}

// GetByVerificationToken finds a user by email verification token
func (r *UserRepository) GetByVerificationToken(ctx context.Context, token string) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM 
			users
		WHERE
			email_verification_token = @token
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"token": token,
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

// GetByPasswordResetToken finds a user by password reset token
func (r *UserRepository) GetByPasswordResetToken(ctx context.Context, token string) (*user.User, error) {
	stmt := `
		SELECT
			*
		FROM 
			users
		WHERE
			password_reset_token = @token
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"token": token,
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

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	stmt := `
		UPDATE users
		SET 
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = @id
	`

	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id": userID,
	})
	return err
}

// CreateOAuthUser creates a user from OAuth provider
func (r *UserRepository) CreateOAuthUser(
	ctx context.Context,
	name, email, provider, providerID string,
) (*user.User, error) {
	stmt := `
		INSERT INTO 
			users (
				name,
				email,
				oauth_provider,
				oauth_provider_id,
				email_verified
			)
		VALUES 
			(
				@name,
				@email,
				@provider,
				@provider_id,
				true
			)
		RETURNING
			*
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"name":         name,
		"email":        email,
		"provider":     provider,
		"provider_id":  providerID,
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
