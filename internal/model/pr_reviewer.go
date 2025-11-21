package model

import (
	"time"

	"github.com/google/uuid"
)

type PrReviewer struct {
	ID            uuid.UUID `json:"id" db:"id"`
	PullRequestId uuid.UUID `json:"pullRequestId" db:"pull_request_id"`
	ReviewerId    uuid.UUID `json:"reviewerId" db:"reviewer_id"`
	AssignedAt    time.Time `json:"assignedAt" db:"assigned_at"`
}
