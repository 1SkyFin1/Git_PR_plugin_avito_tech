package dto

import "github.com/google/uuid"

type UserStatistics struct {
	UserID               uuid.UUID `json:"user_id"`
	Username             string    `json:"username"`
	TeamName             string    `json:"team_name"`
	TotalAssignments     int       `json:"total_assignments"`
	OpenAssignments      int       `json:"open_assignments"`
	CompletedAssignments int       `json:"completed_assignments"`
}

type PullRequestStatistics struct {
	PullRequestID   uuid.UUID `json:"pull_request_id"`
	PullRequestName string    `json:"pull_request_name"`
	AuthorUsername  string    `json:"author_username"`
	Status          string    `json:"status"`
	ReviewersCount  int       `json:"reviewers_count"`
}

type GetStatisticsResponse struct {
	TotalUsers                    int                     `json:"total_users"`
	ActiveUsers                   int                     `json:"active_users"`
	TotalPullRequests             int                     `json:"total_pull_requests"`
	OpenPullRequests              int                     `json:"open_pull_requests"`
	MergedPullRequests            int                     `json:"merged_pull_requests"`
	TotalAssignments              int                     `json:"total_assignments"`
	TopReviewers                  []UserStatistics        `json:"top_reviewers"`
	PullRequestsWithMostReviewers []PullRequestStatistics `json:"pull_requests_with_most_reviewers"`
}

type DeactivateUsersRequest struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}

type DeactivateUsersResponse struct {
	DeactivatedCount   int         `json:"deactivated_count"`
	ReassignedPRsCount int         `json:"reassigned_prs_count"`
	DeactivatedUserIDs []uuid.UUID `json:"deactivated_user_ids"`
}
