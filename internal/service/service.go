package service

import (
	"context"

	"Git_PR_plugin_avito_tech/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, user *model.User) error
	CreateManyTx(ctx context.Context, tx *sqlx.Tx, users []*model.User) error
	SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) (*model.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	FindAllByTeamID(ctx context.Context, teamID uuid.UUID) ([]*model.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, teamName string) (*model.Team, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, teamName string) (*model.Team, error)
	FindByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error)
	FindByName(ctx context.Context, teamName string) (*model.Team, error)
}

type PrReviewerRepository interface {
	FindAllByReviewerID(ctx context.Context, prReviewerID uuid.UUID) ([]*model.PrReviewer, error)
	FindAllByPRID(ctx context.Context, prID uuid.UUID) ([]*model.PrReviewer, error)
	CreateTx(ctx context.Context, tx *sqlx.Tx, prReviewer *model.PrReviewer) error
	CountByReviewerIDs(ctx context.Context, reviewerIDs []uuid.UUID) (map[uuid.UUID]int, error)
	DeleteByPRIDAndReviewerIDTx(ctx context.Context, tx *sqlx.Tx, prID uuid.UUID, reviewerID uuid.UUID) error
}

type PullRequestRepository interface {
	CreateTx(ctx context.Context, tx *sqlx.Tx, pullRequest *model.PullRequest) error
	FindByID(ctx context.Context, pullRequestID uuid.UUID) (*model.PullRequest, error)
	FindAllByPullRequestIDIn(ctx context.Context, prPullRequestIDs []uuid.UUID) ([]*model.PullRequest, error)
	MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (*model.PullRequest, error)
}
