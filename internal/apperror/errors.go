package apperror

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorCode string

const (
	ErrCodeTeamExists   ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists     ErrorCode = "PR_EXISTS"
	ErrCodePRMerged     ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned  ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate  ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeInternal     ErrorCode = "INTERNAL"
)

type AppError struct {
	Code       ErrorCode `json:"-"`
	Message    string    `json:"-"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewAppError(code ErrorCode, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

func (e *AppError) ToErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    e.Code,
			Message: e.Message,
		},
	}
}

func NewTeamExistsError(teamName string) *AppError {
	return &AppError{
		Code:       ErrCodeTeamExists,
		Message:    "team '" + teamName + "' already exists",
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewPRExistsError(prID string) *AppError {
	return &AppError{
		Code:       ErrCodePRExists,
		Message:    "pull request '" + prID + "' already exists",
		HTTPStatus: http.StatusConflict,
	}
}

func NewPRMergedError() *AppError {
	return &AppError{
		Code:       ErrCodePRMerged,
		Message:    "cannot reassign on merged PR",
		HTTPStatus: http.StatusConflict,
	}
}

func NewNotAssignedError() *AppError {
	return &AppError{
		Code:       ErrCodeNotAssigned,
		Message:    "reviewer is not assigned to this PR",
		HTTPStatus: http.StatusConflict,
	}
}

func NewNoCandidateError() *AppError {
	return &AppError{
		Code:       ErrCodeNoCandidate,
		Message:    "no active replacement candidate in team",
		HTTPStatus: http.StatusConflict,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    resource + " not found",
		HTTPStatus: http.StatusNotFound,
	}
}

func NewInvalidInputError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    "internal server error",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, err *AppError) {
	render.Status(r, err.HTTPStatus)
	render.JSON(w, r, err.ToErrorResponse())
}
