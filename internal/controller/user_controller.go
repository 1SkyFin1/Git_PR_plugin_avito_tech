package controller

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/constants"
	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UserController struct {
	userService        *service.UserService
	pullRequestService *service.PullRequestService
}

func NewUserController(userService *service.UserService, pullRequestService *service.PullRequestService) *UserController {
	return &UserController{
		userService:        userService,
		pullRequestService: pullRequestService,
	}
}

// SetIsActive godoc
// @Summary Установить флаг активности пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.SetUserIsActiveRequest true "Запрос на изменение активности пользователя"
// @Success 200 {object} dto.UserResponse "Обновлённый пользователь"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} apperror.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/setIsActive [post]
func (c *UserController) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetUserIsActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appErr := apperror.NewInvalidInputError("Invalid request: " + err.Error())
		apperror.RenderError(w, r, appErr)
		return
	}

	if req.UserID == uuid.Nil {
		appErr := apperror.NewInvalidInputError("Invalid request: user_id is required and must be a valid UUID")
		apperror.RenderError(w, r, appErr)
		return
	}

	response, err := c.userService.SetIsActive(r.Context(), &req)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "User")
		apperror.RenderError(w, r, appErr)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// GetPullRequestsByReviewerID godoc
// @Summary Получить список PR, где пользователь является ревьювером
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя (UUID)"
// @Success 200 {object} dto.GetPullRequestsByReviewerIDResponse "Список PR"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/getReview [get]
func (c *UserController) GetPullRequestsByReviewerID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get(constants.UserIDQueryParameter)

	if userIDStr == "" {
		appErr := apperror.NewInvalidInputError("Invalid request: user_id is required")
		apperror.RenderError(w, r, appErr)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		appErr := apperror.NewInvalidInputError("Invalid request: user_id is invalid UUID")
		apperror.RenderError(w, r, appErr)
		return
	}

	response, err := c.pullRequestService.GetPullRequestsByReviewerID(r.Context(), userID)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "PullRequest")
		apperror.RenderError(w, r, appErr)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
