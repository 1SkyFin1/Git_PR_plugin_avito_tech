package repository

import (
	"Git_PR_plugin_avito_tech/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

	reviewerIDStrings := make([]string, len(reviewerIDs))
	for i, id := range reviewerIDs {
		reviewerIDStrings[i] = id.String()
	}

	rows, err := r.db.QueryContext(ctx, query, pq.Array(reviewerIDStrings))
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

func (r *PrReviewerRepository) GetUserStatisticsData(ctx context.Context, limit int) ([]*model.UserStatisticsData, error) {
	query := `
		SELECT 
			u.user_id,
			u.username,
			t.team_name,
			COUNT(pr.pull_request_id) as total_assignments,
			COUNT(CASE WHEN pr.status = 'OPEN' THEN 1 END) as open_assignments,
			COUNT(CASE WHEN pr.status = 'MERGED' THEN 1 END) as completed_assignments
		FROM "user" u
		INNER JOIN team t ON u.team_id = t.id
		LEFT JOIN pr_reviewer prr ON u.user_id = prr.reviewer_id
		LEFT JOIN pull_request pr ON prr.pull_request_id = pr.pull_request_id
		WHERE u.is_active = true
		GROUP BY u.user_id, u.username, t.team_name
		ORDER BY total_assignments DESC
		LIMIT $1
	`

	var results []*model.UserStatisticsData
	err := r.db.SelectContext(ctx, &results, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	return results, nil
}

func (r *PrReviewerRepository) GetPRStatisticsData(ctx context.Context, limit int) ([]*model.PRStatisticsData, error) {
	query := `
		SELECT 
			pr.pull_request_id,
			pr.pull_request_name,
			u.username as author_username,
			pr.status,
			COUNT(prr.reviewer_id) as reviewers_count
		FROM pull_request pr
		INNER JOIN "user" u ON pr.author_id = u.user_id
		LEFT JOIN pr_reviewer prr ON pr.pull_request_id = prr.pull_request_id
		GROUP BY pr.pull_request_id, pr.pull_request_name, u.username, pr.status
		ORDER BY reviewers_count DESC
		LIMIT $1
	`

	var results []*model.PRStatisticsData
	err := r.db.SelectContext(ctx, &results, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR statistics: %w", err)
	}

	return results, nil
}

func (r *PrReviewerRepository) GetTotalAssignments(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM pr_reviewer`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get total assignments: %w", err)
	}

	return count, nil
}

func (r *PrReviewerRepository) DeleteByReviewerIDsTx(ctx context.Context, tx *sqlx.Tx, reviewerIDs []uuid.UUID) (int, error) {
	if len(reviewerIDs) == 0 {
		return 0, nil
	}

	query := `DELETE FROM pr_reviewer WHERE reviewer_id = ANY($1)`

	result, err := tx.ExecContext(ctx, query, pq.Array(reviewerIDs))
	if err != nil {
		return 0, fmt.Errorf("failed to delete pr_reviewers: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}
