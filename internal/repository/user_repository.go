package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Git_PR_plugin_avito_tech/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, user *model.User) error {
	query := `
		INSERT INTO "user" (username, team_id, is_active)
		VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query,
		user.Username,
		user.TeamID,
		user.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) CreateManyTx(ctx context.Context, tx *sqlx.Tx, users []*model.User) error {
	if len(users) == 0 {
		return nil
	}

	query := `
		INSERT INTO "user" (username, team_id, is_active)
		VALUES (:username, :team_id, :is_active)
	`

	_, err := tx.NamedExecContext(ctx, query, users)
	if err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}

	return nil
}

func (r *UserRepository) SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) (*model.User, error) {
	query := `
		UPDATE "user"
		SET is_active = $1, updated_at = NOW()
		WHERE user_id = $2
		RETURNING *
	`

	var user model.User
	err := r.db.QueryRowxContext(ctx, query, isActive, userID).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to set isActive: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	query := `
		SELECT user_id, username, team_id, is_active, created_at, updated_at
		FROM "user"
		WHERE user_id = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindAllByTeamID(ctx context.Context, teamID uuid.UUID) ([]*model.User, error) {
	query := `SELECT * FROM "user" WHERE team_id = $1`

	var users []*model.User
	err := r.db.SelectContext(ctx, &users, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to find users with teamID %s: %w", teamID, err)
	}
	return users, nil
}

func (r *UserRepository) GetTotalUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM "user"`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get total users: %w", err)
	}

	return count, nil
}

func (r *UserRepository) GetActiveUsers(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM "user" WHERE is_active = true`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get active users: %w", err)
	}

	return count, nil
}

func (r *UserRepository) DeactivateUsersTx(ctx context.Context, tx *sqlx.Tx, userIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(userIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	query := `
		UPDATE "user"
		SET is_active = false, updated_at = NOW()
		WHERE user_id = ANY($1) AND is_active = true
		RETURNING user_id
	`

	rows, err := tx.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate users: %w", err)
	}
	defer rows.Close()

	var deactivatedIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan user_id: %w", err)
		}
		deactivatedIDs = append(deactivatedIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return deactivatedIDs, nil
}
