package service

import (
	"Git_PR_plugin_avito_tech/internal/dto"
	"context"
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

func (u *UserService) SetIsActive(ctx context.Context, setUserActiveRequest dto.SetUserIsActiveRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.SetIsActive(ctx, setUserActiveRequest.UserID, setUserActiveRequest.IsActive)
	if err != nil {
		return nil, err
	}

	team, err := u.teamRepo.FindByID(ctx, user.TeamID)
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user, team.TeamName), nil
}
