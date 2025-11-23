package service

import (
	"context"
	"fmt"
	"log"

	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/model"

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

func (s *TeamService) CreateTeamWithMembers(ctx context.Context, teamInfo *dto.TeamInfo) (*dto.TeamResponse, error) {
	log.Printf("[TeamService.CreateTeamWithMembers] Starting for team: %s with %d members", teamInfo.TeamName, len(teamInfo.Members))

	if teamInfo.TeamName == "" {
		return nil, apperror.NewInvalidInputError("team name cannot be empty")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[TeamService.CreateTeamWithMembers] Error beginning transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback()

	team, err := s.teamRepo.CreateTx(ctx, tx, teamInfo.TeamName)
	if err != nil {
		log.Printf("[TeamService.CreateTeamWithMembers] Error creating team: %v", err)
		return nil, apperror.HandleDBError(err, "Team")
	}

	log.Printf("[TeamService.CreateTeamWithMembers] Team created: %s (ID: %s)", team.TeamName, team.ID)

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
			log.Printf("[TeamService.CreateTeamWithMembers] Error creating users: %v", err)
			return nil, apperror.NewInternalError(fmt.Errorf("failed to create users: %w", err))
		}
		log.Printf("[TeamService.CreateTeamWithMembers] Created %d users", len(users))
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[TeamService.CreateTeamWithMembers] Error committing transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to commit transaction: %w", err))
	}

	log.Printf("[TeamService.CreateTeamWithMembers] Success! Team: %s", team.TeamName)
	return dto.ToTeamResponse(team, users), nil
}

func (s *TeamService) GetTeamByName(ctx context.Context, teamName string) (*dto.TeamInfo, error) {
	log.Printf("[TeamService.GetTeamByName] Starting for teamName: %s", teamName)
	
	if teamName == "" {
		return nil, apperror.NewInvalidInputError("team name cannot be empty")
	}

	team, err := s.teamRepo.FindByName(ctx, teamName)
	if err != nil {
		log.Printf("[TeamService.GetTeamByName] Error finding team: %v", err)
		if apperror.IsNotFoundError(err) {
			return nil, apperror.NewNotFoundError("Team")
		}
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[TeamService.GetTeamByName] Team found: %s, fetching users...", team.TeamName)

	users, err := s.userRepo.FindAllByTeamID(ctx, team.ID)
	if err != nil {
		log.Printf("[TeamService.GetTeamByName] Error fetching users: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[TeamService.GetTeamByName] Success! Team: %s, Users count: %d", team.TeamName, len(users))
	return dto.ToTeamInfo(team, dto.ToMembersDto(users)), nil
}
