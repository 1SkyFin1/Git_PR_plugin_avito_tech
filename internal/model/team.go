package model

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TeamName  string    `json:"teamName" db:"team_name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
