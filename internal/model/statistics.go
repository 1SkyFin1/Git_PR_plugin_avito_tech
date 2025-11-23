package model

import "github.com/google/uuid"

type UserStatisticsData struct {
	UserID               uuid.UUID `db:"user_id"`
	Username             string    `db:"username"`
	TeamName             string    `db:"team_name"`
	TotalAssignments     int       `db:"total_assignments"`
	OpenAssignments      int       `db:"open_assignments"`
	CompletedAssignments int       `db:"completed_assignments"`
}

type PRStatisticsData struct {
	PullRequestID   uuid.UUID `db:"pull_request_id"`
	PullRequestName string    `db:"pull_request_name"`
	AuthorUsername  string    `db:"author_username"`
	Status          string    `db:"status"`
	ReviewersCount  int       `db:"reviewers_count"`
}
