package service

import (
	"context"
	"fmt"

	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeamService struct {
	teamRepo TeamRepository
	userRepo UserRepository
	db       *sqlx.DB
}

func NewTeamService(teamRepo TeamRepository, userRepo UserRepository, db *sqlx.DB) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		db:       db,
	}
}

func (s *TeamService) CreateTeamWithMembers(ctx context.Context, teamInfo dto.TeamInfo) (*dto.TeamResponse, error) {
	if teamInfo.TeamName == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	team, err := s.teamRepo.CreateTx(ctx, tx, teamInfo.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	users := make([]*model.User, 0, len(teamInfo.Members))
	for _, member := range teamInfo.Members {
		if member.Username == "" {
			continue
		}

		users = append(users, dto.ToUserModelFromMemberDto(&member, team.ID))
	}

	if len(users) > 0 {
		err := s.userRepo.CreateManyTx(ctx, tx, users)
		if err != nil {
			return nil, fmt.Errorf("failed to create users: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return dto.ToTeamResponse(team, users), nil
}

func (s *TeamService) CreateTeam(ctx context.Context, teamName string) (*model.Team, error) {
	if teamName == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}

	team, err := s.teamRepo.Create(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
	if teamID == uuid.Nil {
		return nil, fmt.Errorf("team ID cannot be empty")
	}

	team, err := s.teamRepo.FindByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

func (s *TeamService) GetTeamByName(ctx context.Context, teamName string) (*dto.TeamInfo, error) {
	if teamName == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}

	team, err := s.teamRepo.FindByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	users, err := s.userRepo.FindAllByTeamID(ctx, team.ID)
	if err != nil {
		return nil, err
	}

	return dto.ToTeamInfo(team, dto.ToMembersDto(users)), nil
}
