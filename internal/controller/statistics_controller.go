package controller

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type StatisticsController struct {
	statisticsService *service.StatisticsService
}

func NewStatisticsController(statisticsService *service.StatisticsService) *StatisticsController {
	return &StatisticsController{
		statisticsService: statisticsService,
	}
}

// GetStatistics godoc
// @Summary Получить статистику системы
// @Description Возвращает общую статистику: количество пользователей, PR, назначений, топ ревьюверов и PR с наибольшим количеством ревьюверов
// @Tags Statistics
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetStatisticsResponse "Статистика системы"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /statistics [get]
func (c *StatisticsController) GetStatistics(w http.ResponseWriter, r *http.Request) {
	response, err := c.statisticsService.GetStatistics(r.Context())
	if err != nil {
		appErr := apperror.WrapServiceError(err, "Statistics")
		apperror.RenderError(w, r, appErr)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// DeactivateUsers godoc
// @Summary Массовая деактивация пользователей
// @Description Деактивирует указанных пользователей и удаляет все их назначения на ревью PR
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.DeactivateUsersRequest true "Список ID пользователей для деактивации"
// @Success 200 {object} dto.DeactivateUsersResponse "Результат деактивации"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/deactivate [post]
func (c *StatisticsController) DeactivateUsers(w http.ResponseWriter, r *http.Request) {
	var req dto.DeactivateUsersRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appErr := apperror.NewInvalidInputError("Invalid request: " + err.Error())
		apperror.RenderError(w, r, appErr)
		return
	}

	for _, userID := range req.UserIDs {
		if userID == uuid.Nil {
			appErr := apperror.NewInvalidInputError("Invalid request: all user_ids must be valid UUIDs")
			apperror.RenderError(w, r, appErr)
			return
		}
	}

	response, err := c.statisticsService.DeactivateUsers(r.Context(), &req)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "User")
		apperror.RenderError(w, r, appErr)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
