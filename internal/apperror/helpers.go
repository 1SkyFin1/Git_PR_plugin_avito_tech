package apperror

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
	pgNotNullViolation    = "23502"
)

func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func HandleDBError(err error, resource string) *AppError {
	if err == nil {
		return nil
	}

	if IsNotFoundError(err) {
		return NewNotFoundError(resource)
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case pgUniqueViolation:
			if pqErr.Constraint == "teams_team_name_key" {
				return NewTeamExistsError("team")
			}
			if pqErr.Constraint == "pull_requests_pkey" {
				return NewPRExistsError("pull request")
			}
			return NewInvalidInputError("resource already exists")
		case pgForeignKeyViolation:
			return NewNotFoundError("related resource")
		case pgNotNullViolation:
			return NewInvalidInputError("required field is missing: " + pqErr.Column)
		case pgCheckViolation:
			return NewInvalidInputError("constraint violation: " + pqErr.Message)
		}
	}

	return NewInternalError(err)
}

func WrapServiceError(err error, defaultResource string) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return HandleDBError(err, defaultResource)
}
