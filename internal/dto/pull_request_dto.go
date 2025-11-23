package dto

import (
	"Git_PR_plugin_avito_tech/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreatePullRequestRequest struct {
	PullRequestID   uuid.UUID `json:"pull_request_id"`
	PullRequestName string    `json:"pull_request_name"`
	AuthorID        uuid.UUID `json:"author_id"`
}

type MergePullRequestRequest struct {
	PullRequestID uuid.UUID `json:"pull_request_id"`
}

type ReassignPrReviewerRequest struct {
	PullRequestID uuid.UUID `json:"pull_request_id"`
	OldReviewerID uuid.UUID `json:"old_reviewer_id"`
}

type CreatePullRequestResponse struct {
	PullRequestID     uuid.UUID               `json:"pull_request_id"`
	PullRequestName   string                  `json:"pull_request_name"`
	AuthorID          uuid.UUID               `json:"author_id"`
	Status            model.PullRequestStatus `json:"status"`
	AssignedReviewers []uuid.UUID             `json:"assigned_reviewers"`
}

type GetPullRequestsByReviewerIDResponse struct {
	UserID       uuid.UUID        `json:"user_id"`
	PullRequests []PullRequestDTO `json:"pull_requests"`
}

type MergePullRequestResponse struct {
	PR       PullRequestFullInfoDto `json:"pr"`
	MergedAt *time.Time             `json:"mergedAt,omitempty"`
}

type ReassignPrReviewerResponse struct {
	PR         PullRequestFullInfoDto `json:"pr"`
	ReplacedBy uuid.UUID              `json:"replaced_by"`
}

type PullRequestFullInfoDto struct {
	PullRequestID     uuid.UUID               `json:"pull_request_id"`
	PullRequestName   string                  `json:"pull_request_name"`
	AuthorID          uuid.UUID               `json:"author_id"`
	Status            model.PullRequestStatus `json:"status"`
	AssignedReviewers []uuid.UUID             `json:"assigned_reviewers"`
}

type PullRequestDTO struct {
	PullRequestID   uuid.UUID               `json:"pull_request_id"`
	PullRequestName string                  `json:"pull_request_name"`
	AuthorID        uuid.UUID               `json:"author_id"`
	Status          model.PullRequestStatus `json:"status"`
}
