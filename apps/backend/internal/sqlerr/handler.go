package sqlerr

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/reche13/habitum/internal/errs"
)

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return handlePostgresError(pgErr)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return errs.NewNotFoundError("Resource not found")
	}

	return fmt.Errorf("database error: %w", err)
}

func handlePostgresError(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	case "23505": // unique_violation
		return errs.NewBadRequestError(fmt.Sprintf("Duplicate entry: %s", pgErr.ConstraintName))
	case "23503": // foreign_key_violation
		return errs.NewBadRequestError("Referenced resource does not exist")
	case "23502": // not_null_violation
		return errs.NewBadRequestError(fmt.Sprintf("Required field missing: %s", pgErr.ColumnName))
	case "23514": // check_violation
		return errs.NewBadRequestError("Data validation failed")
	case "42P01": // undefined_table
		return errs.NewInternalServerError("Database configuration error")
	case "42703": // undefined_column
		return errs.NewInternalServerError("Database schema error")
	default:
		return fmt.Errorf("database error: %s", pgErr.Message)
	}
}

