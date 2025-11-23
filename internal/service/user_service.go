package service

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/dto"
	"context"
	"log"
)

type UserService struct {
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewUserService(userRepo UserRepository, teamRepo TeamRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (u *UserService) SetIsActive(ctx context.Context, setUserActiveRequest *dto.SetUserIsActiveRequest) (*dto.UserResponse, error) {
	log.Printf("[UserService.SetIsActive] Starting for userID: %s, isActive: %v", setUserActiveRequest.UserID, setUserActiveRequest.IsActive)
	
	user, err := u.userRepo.SetIsActive(ctx, setUserActiveRequest.UserID, setUserActiveRequest.IsActive)
	if err != nil {
		log.Printf("[UserService.SetIsActive] Error setting isActive: %v", err)
		if apperror.IsNotFoundError(err) {
			return nil, apperror.NewNotFoundError("User")
		}
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[UserService.SetIsActive] User updated successfully, fetching team for teamID: %s", user.TeamID)
	
	team, err := u.teamRepo.FindByID(ctx, user.TeamID)
	if err != nil {
		log.Printf("[UserService.SetIsActive] Error fetching team: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[UserService.SetIsActive] Success for userID: %s", setUserActiveRequest.UserID)
	return dto.ToUserResponse(user, team.TeamName), nil
}
