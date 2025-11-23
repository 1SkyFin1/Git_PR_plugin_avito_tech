package dto

import (
	"Git_PR_plugin_avito_tech/internal/model"

	"github.com/google/uuid"
)

func ToUserResponse(user *model.User, teamName string) *UserResponse {
	return &UserResponse{*ToUserInfoDto(user, teamName)}
}

func ToUserModelFromRequest(req *CreateUserRequest) *model.User {
	if req == nil {
		return nil
	}

	return &model.User{
		Username: req.Username,
		TeamID:   req.TeamID,
		IsActive: true,
	}
}

func ToUserModelFromMemberDto(member *MemberDto, teamID uuid.UUID) *model.User {
	if member == nil {
		return nil
	}

	return &model.User{
		UserID:   member.UserID,
		Username: member.Username,
		TeamID:   teamID,
		IsActive: member.IsActive,
	}
}

func ToMemberDto(user *model.User) *MemberDto {
	if user == nil {
		return nil
	}

	return &MemberDto{
		UserID:   user.UserID,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

func ToMembersDto(users []*model.User) []MemberDto {
	result := make([]MemberDto, 0, len(users))

	for _, user := range users {
		if user == nil {
			continue
		}
		result = append(result, *ToMemberDto(user))
	}
	return result
}

func ToUserInfoDto(user *model.User, teamName string) *UserInfoDto {
	if user == nil {
		return nil
	}

	return &UserInfoDto{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: teamName,
		IsActive: user.IsActive,
	}
}
