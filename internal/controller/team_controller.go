package controller

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/constants"
	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

type TeamController struct {
	teamService *service.TeamService
}

func NewTeamController(teamService *service.TeamService) *TeamController {
	return &TeamController{teamService: teamService}
}

// CreateTeamWithMembers godoc
// @Summary Создать команду с участниками
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body dto.TeamInfo true "Информация о команде и участниках"
// @Success 201 {object} dto.TeamResponse "Созданная команда"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 409 {object} apperror.ErrorResponse "Команда уже существует"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /team/add [post]
func (c *TeamController) CreateTeamWithMembers(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamInfo

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appErr := apperror.NewInvalidInputError("Invalid request: " + err.Error())
		apperror.RenderError(w, r, appErr)
		return
	}

	if req.TeamName == "" {
		appErr := apperror.NewInvalidInputError("Invalid request: team_name is required")
		apperror.RenderError(w, r, appErr)
		return
	}

	response, err := c.teamService.CreateTeamWithMembers(r.Context(), &req)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "Team")
		apperror.RenderError(w, r, appErr)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// GetTeamByName godoc
// @Summary Получить команду по названию
// @Tags Teams
// @Accept json
// @Produce json
// @Param team_name query string true "Название команды"
// @Success 200 {object} dto.TeamInfo "Информация о команде"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} apperror.ErrorResponse "Команда не найдена"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /team/get [get]
func (c *TeamController) GetTeamByName(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get(constants.TeamNameQueryParameter)

	if teamName == "" {
		appErr := apperror.NewInvalidInputError("Invalid request: team_name is required")
		apperror.RenderError(w, r, appErr)
		return
	}

	response, err := c.teamService.GetTeamByName(r.Context(), teamName)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "Team")
		apperror.RenderError(w, r, appErr)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
