package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId    uuid.UUID `json:"userId" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	TeamId    uuid.UUID `json:"teamId" db:"team_id"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
