package repository

import (
	"Git_PR_plugin_avito_tech/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PullRequestRepository struct {
	db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) PullRequestRepository {
	return PullRequestRepository{
		db: db,
	}
}

func (r *PullRequestRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, pullRequest *model.PullRequest) error {
	query := `INSERT INTO pull_request (pull_request_id, pull_request_name, author_id, status)
			VALUES (:pull_request_id, :pull_request_name, :author_id, :status)`

	_, err := tx.NamedExecContext(ctx, query, pullRequest)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}
	return nil
}

func (r *PullRequestRepository) FindByID(ctx context.Context, pullRequestID uuid.UUID) (*model.PullRequest, error) {
	query := `SELECT * FROM pull_request WHERE pull_request_id = $1`

	var pr model.PullRequest
	err := r.db.GetContext(ctx, &pr, query, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pull request by id %s: %w", pullRequestID, err)
	}
	return &pr, nil
}

func (r *PullRequestRepository) FindAllByPullRequestIDIn(ctx context.Context, prPullRequestIDs []uuid.UUID) ([]*model.PullRequest, error) {
	query := `SELECT * FROM pull_request WHERE pull_request_id IN $1`

	var prs []*model.PullRequest
	err := r.db.SelectContext(ctx, &prs, query, prPullRequestIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find pull requests by pull request ids: %w", err)
	}
	return prs, nil
}

func (r *PullRequestRepository) MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (*model.PullRequest, error) {
	query := `
			INSERT INTO pull_request (status, merged_at) VALUES ($1, $2) 
			WHERE pull_request_id = :$3
			RETURNING *`

	var pullRequest model.PullRequest
	err := r.db.QueryRowxContext(ctx, query, model.PullRequestStatusMerged, time.Now(), pullRequestID).StructScan(&pullRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to merge pull request with id %s: %w", pullRequestID, err)
	}
	return &pullRequest, nil
}
