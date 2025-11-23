package dto

import (
	"github.com/google/uuid"
)

type SetUserIsActiveRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	IsActive bool      `json:"is_active" binding:"required"`
}

type UserResponse struct {
	UserInfo UserInfoDto `json:"user"`
}

type MemberDto struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	IsActive bool      `json:"is_active"`
}

type UserInfoDto struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	TeamName string    `json:"team_name"`
	IsActive bool      `json:"is_active"`
}
