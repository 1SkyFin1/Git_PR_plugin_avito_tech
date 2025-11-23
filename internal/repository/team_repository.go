package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Git_PR_plugin_avito_tech/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, teamName string) (*model.Team, error) {
	query := `
		INSERT INTO team (team_name)
		VALUES ($1)
		RETURNING *
	`

	var createdTeam model.Team
	err := tx.QueryRowxContext(ctx, query, teamName).StructScan(&createdTeam)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return &createdTeam, nil
}

func (r *TeamRepository) FindByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
	query := `
		SELECT *
		FROM team
		WHERE id = $1
	`

	var team model.Team
	err := r.db.GetContext(ctx, &team, query, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to find team with id %s: %w", teamID, err)
	}

	return &team, nil
}

func (r *TeamRepository) FindByName(ctx context.Context, teamName string) (*model.Team, error) {
	query := `
		SELECT * from team WHERE team_name = $1`

	var team model.Team
	err := r.db.GetContext(ctx, &team, query, teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to find team with name %s: %w", teamName, err)
	}

	return &team, nil
}
