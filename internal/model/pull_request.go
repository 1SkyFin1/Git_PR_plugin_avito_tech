package model

import (
	"time"

	"github.com/google/uuid"
)

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestId   uuid.UUID         `json:"pullRequestId" db:"pull_request_id"`
	PullRequestName string            `json:"pullRequestName" db:"pull_request_name"`
	AuthorId        uuid.UUID         `json:"authorId" db:"author_id"`
	Status          PullRequestStatus `json:"status" db:"status"`
	CreatedAt       time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time         `json:"updatedAt" db:"updated_at"`
}
