package controller

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type PullRequestController struct {
	pullRequestService *service.PullRequestService
}

func NewPullRequestController(pullRequestService *service.PullRequestService) *PullRequestController {
	return &PullRequestController{pullRequestService: pullRequestService}
}

// CreatePullRequest godoc
// @Summary Создать Pull Request
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.CreatePullRequestRequest true "Запрос на создание PR"
// @Success 201 {object} dto.CreatePullRequestResponse "Созданный PR с назначенными ревьюверами"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} apperror.ErrorResponse "Автор не найден"
// @Failure 409 {object} apperror.ErrorResponse "PR уже существует"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pullRequest/create [post]
func (c *PullRequestController) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.CreatePullRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[PullRequestController.CreatePullRequest] Failed to decode request: %v", err)
		appErr := apperror.NewInvalidInputError("Invalid request: " + err.Error())
		apperror.RenderError(w, r, appErr)
		return
	}

	log.Printf("[PullRequestController.CreatePullRequest] Request received: PR=%s, Author=%s", request.PullRequestID, request.AuthorID)

	response, err := c.pullRequestService.CreatePullRequest(r.Context(), &request)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "Pull Request")
		apperror.RenderError(w, r, appErr)
		return
	}

	log.Printf("[PullRequestController.CreatePullRequest] PR created successfully: %s", request.PullRequestID)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// MergePullRequest godoc
// @Summary Смержить Pull Request
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.MergePullRequestRequest true "Запрос на мерж PR"
// @Success 200 {object} dto.MergePullRequestResponse "Смерженный PR"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} apperror.ErrorResponse "PR не найден"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pullRequest/merge [post]
func (c *PullRequestController) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.MergePullRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[PullRequestController.MergePullRequest] Failed to decode request: %v", err)
		appErr := apperror.NewInvalidInputError("Invalid request: " + err.Error())
		apperror.RenderError(w, r, appErr)
		return
	}

	log.Printf("[PullRequestController.MergePullRequest] Request received: PR=%s", request.PullRequestID)

	response, err := c.pullRequestService.MergePullRequest(r.Context(), &request)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "Pull Request")
		apperror.RenderError(w, r, appErr)
		return
	}

	log.Printf("[PullRequestController.MergePullRequest] PR merged successfully: %s", request.PullRequestID)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// ReassignPrReviewer godoc
// @Summary Переназначить ревьювера PR
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.ReassignPrReviewerRequest true "Запрос на переназначение ревьювера"
// @Success 200 {object} dto.ReassignPrReviewerResponse "PR с обновленными ревьюверами"
// @Failure 400 {object} apperror.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} apperror.ErrorResponse "PR или ревьювер не найден"
// @Failure 409 {object} apperror.ErrorResponse "PR уже смержен или нет доступных ревьюверов"
// @Failure 500 {object} apperror.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pullRequest/reassign [post]
func (c *PullRequestController) ReassignPrReviewer(w http.ResponseWriter, r *http.Request) {
	var request dto.ReassignPrReviewerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[PullRequestController.ReassignPrReviewer] Failed to decode request: %v", err)
		appErr := apperror.NewInvalidInputError("Invalid request: " + err.Error())
		apperror.RenderError(w, r, appErr)
		return
	}

	log.Printf("[PullRequestController.ReassignPrReviewer] Request received: PR=%s, OldReviewer=%s", request.PullRequestID, request.OldReviewerID)

	response, err := c.pullRequestService.ReassignPrReviewer(r.Context(), &request)
	if err != nil {
		appErr := apperror.WrapServiceError(err, "Pull Request")
		apperror.RenderError(w, r, appErr)
		return
	}

	log.Printf("[PullRequestController.ReassignPrReviewer] Reviewer reassigned successfully for PR: %s", request.PullRequestID)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
