package repository

import (
	"Git_PR_plugin_avito_tech/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PullRequestRepository struct {
	db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) *PullRequestRepository {
	return &PullRequestRepository{
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
	if len(prPullRequestIDs) == 0 {
		return []*model.PullRequest{}, nil
	}

	query := `SELECT * FROM pull_request WHERE pull_request_id = ANY($1)`

	prIDStrings := make([]string, len(prPullRequestIDs))
	for i, id := range prPullRequestIDs {
		prIDStrings[i] = id.String()
	}

	var prs []*model.PullRequest
	err := r.db.SelectContext(ctx, &prs, query, pq.Array(prIDStrings))
	if err != nil {
		return nil, fmt.Errorf("failed to find pull requests by pull request ids: %w", err)
	}
	return prs, nil
}

func (r *PullRequestRepository) MergePullRequest(ctx context.Context, pullRequestID uuid.UUID) (*model.PullRequest, error) {
	query := `
			UPDATE pull_request 
			SET status = $1, merged_at = $2 
			WHERE pull_request_id = $3
			RETURNING *`

	var pullRequest model.PullRequest
	err := r.db.QueryRowxContext(ctx, query, model.PullRequestStatusMerged, time.Now(), pullRequestID).StructScan(&pullRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to merge pull request with id %s: %w", pullRequestID, err)
	}
	return &pullRequest, nil
}

func (r *PullRequestRepository) GetTotalPullRequests(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM pull_request`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get total pull requests: %w", err)
	}

	return count, nil
}

func (r *PullRequestRepository) GetOpenPullRequests(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM pull_request WHERE status = $1`

	var count int
	err := r.db.GetContext(ctx, &count, query, model.PullRequestStatusOpen)
	if err != nil {
		return 0, fmt.Errorf("failed to get open pull requests: %w", err)
	}

	return count, nil
}

func (r *PullRequestRepository) GetMergedPullRequests(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM pull_request WHERE status = $1`

	var count int
	err := r.db.GetContext(ctx, &count, query, model.PullRequestStatusMerged)
	if err != nil {
		return 0, fmt.Errorf("failed to get merged pull requests: %w", err)
	}

	return count, nil
}

func (r *PullRequestRepository) FindOpenPRsByReviewerIDs(ctx context.Context, reviewerIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(reviewerIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	query := `
		SELECT DISTINCT pr.pull_request_id
		FROM pull_request pr
		INNER JOIN pr_reviewer prr ON pr.pull_request_id = prr.pull_request_id
		WHERE pr.status = $1 AND prr.reviewer_id = ANY($2)
	`

	rows, err := r.db.QueryContext(ctx, query, model.PullRequestStatusOpen, pq.Array(reviewerIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to find open PRs by reviewer IDs: %w", err)
	}
	defer rows.Close()

	var prIDs []uuid.UUID
	for rows.Next() {
		var prID uuid.UUID
		if err := rows.Scan(&prID); err != nil {
			return nil, fmt.Errorf("failed to scan PR ID: %w", err)
		}
		prIDs = append(prIDs, prID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return prIDs, nil
}
