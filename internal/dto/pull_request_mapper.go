package dto

import (
	"Git_PR_plugin_avito_tech/internal/model"

	"github.com/google/uuid"
)

func ToModelFromCreatePullRequestRequest(request *CreatePullRequestRequest) *model.PullRequest {
	if request == nil {
		return nil
	}

	return &model.PullRequest{
		PullRequestID:   request.PullRequestID,
		PullRequestName: request.PullRequestName,
		AuthorID:        request.AuthorID,
	}
}

func ToCreatePullRequestResponse(pr *model.PullRequest, reviewerIDs []uuid.UUID) *CreatePullRequestResponse {
	if pr == nil {
		return nil
	}

	return &CreatePullRequestResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: reviewerIDs,
	}
}

func ToPullRequestDTO(pr *model.PullRequest) *PullRequestDTO {
	if pr == nil {
		return nil
	}

	return &PullRequestDTO{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
	}
}

func ToGetPullRequestsByReviewerIDResponse(userID uuid.UUID, prs []*model.PullRequest) *GetPullRequestsByReviewerIDResponse {
	if userID == uuid.Nil || prs == nil {
		return nil
	}

	return &GetPullRequestsByReviewerIDResponse{
		UserID:       userID,
		PullRequests: ToPullRequestDTOs(prs),
	}
}

func ToPullRequestDTOs(prs []*model.PullRequest) []PullRequestDTO {
	pullRequestDTOs := make([]PullRequestDTO, 0, len(prs))

	for _, pr := range prs {
		pullRequestDTOs = append(pullRequestDTOs, *ToPullRequestDTO(pr))
	}
	return pullRequestDTOs
}

func ToPullRequestFullInfoDto(pr *model.PullRequest, assignedReviewers []uuid.UUID) *PullRequestFullInfoDto {
	if pr == nil {
		return nil
	}

	return &PullRequestFullInfoDto{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: assignedReviewers,
	}
}

func ToMergePullRequestResponse(pr *model.PullRequest, assignedReviewers []uuid.UUID) *MergePullRequestResponse {
	if pr == nil {
		return nil
	}

	return &MergePullRequestResponse{
		PR:       *ToPullRequestFullInfoDto(pr, assignedReviewers),
		MergedAt: pr.MergedAt,
	}
}

func ToReassignPrReviewerResponse(pr *model.PullRequest, assignedReviewers []uuid.UUID, newReviewer uuid.UUID) *ReassignPrReviewerResponse {
	if pr == nil {
		return nil
	}

	return &ReassignPrReviewerResponse{
		PR:         *ToPullRequestFullInfoDto(pr, assignedReviewers),
		ReplacedBy: newReviewer,
	}
}
