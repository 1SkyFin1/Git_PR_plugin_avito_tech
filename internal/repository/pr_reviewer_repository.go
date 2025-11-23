package repository

import (
	"Git_PR_plugin_avito_tech/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PrReviewerRepository struct {
	db *sqlx.DB
}

func NewPrReviewerRepository(db *sqlx.DB) *PrReviewerRepository {
	return &PrReviewerRepository{
		db: db,
	}
}

func (r *PrReviewerRepository) FindAllByReviewerID(ctx context.Context, prReviewerID uuid.UUID) ([]*model.PrReviewer, error) {
	query := `SELECT * FROM pr_reviewer WHERE reviewer_id = $1`

	var prReviewers []*model.PrReviewer
	err := r.db.SelectContext(ctx, &prReviewers, query, prReviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pr_reviewer with reviewer_id %s: %w", prReviewerID, err)
	}
	return prReviewers, nil
}

func (r *PrReviewerRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, prReviewer *model.PrReviewer) error {
	query := `
		INSERT INTO pr_reviewer (pull_request_id, reviewer_id)
		VALUES ($1, $2)
		RETURNING id, assigned_at
	`

	err := tx.QueryRowxContext(ctx, query, prReviewer.PullRequestId, prReviewer.ReviewerId).
		Scan(&prReviewer.ID, &prReviewer.AssignedAt)
	if err != nil {
		return fmt.Errorf("failed to create pr_reviewer: %w", err)
	}

	return nil
}

func (r *PrReviewerRepository) CountByReviewerIDs(ctx context.Context, reviewerIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(reviewerIDs) == 0 {
		return make(map[uuid.UUID]int), nil
	}

	query := `
		SELECT reviewer_id, COUNT(*) as count
		FROM pr_reviewer
		WHERE reviewer_id = ANY($1)
		GROUP BY reviewer_id
	`

	rows, err := r.db.QueryContext(ctx, query, reviewerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to count pr_reviewers: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]int)
	for rows.Next() {
		var reviewerID uuid.UUID
		var count int
		if err := rows.Scan(&reviewerID, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[reviewerID] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

func (r *PrReviewerRepository) FindAllByPRID(ctx context.Context, prID uuid.UUID) ([]*model.PrReviewer, error) {
	query := `SELECT * FROM pr_reviewer WHERE pull_request_id = $1`

	var prReviewers []*model.PrReviewer
	err := r.db.SelectContext(ctx, &prReviewers, query, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pr_reviewer with pull_request_id %s: %w", prID, err)
	}
	return prReviewers, nil
}

func (r *PrReviewerRepository) DeleteByPRIDAndReviewerIDTx(ctx context.Context, tx *sqlx.Tx, prID uuid.UUID, reviewerID uuid.UUID) error {
	query := `DELETE FROM pr_reviewer WHERE pull_request_id = $1 AND reviewer_id = $2`

	result, err := tx.ExecContext(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to delete pr_reviewer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pr_reviewer not found")
	}

	return nil
}
