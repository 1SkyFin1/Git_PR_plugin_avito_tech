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
	PullRequestID   uuid.UUID         `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        uuid.UUID         `json:"author_id" db:"author_id"`
	Status          PullRequestStatus `json:"status" db:"status"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	MergedAt        time.Time         `json:"merged_at" db:"merged_at"`
}
