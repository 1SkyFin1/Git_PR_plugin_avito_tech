package app

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"Git_PR_plugin_avito_tech/internal/repository"
	"Git_PR_plugin_avito_tech/internal/service"
)

func Run(db *sqlx.DB) {
	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)

	teamService := service.NewTeamService(teamRepo, userRepo, db)

	_ = teamService

	ctx := context.Background()

	log.Println("Application started successfully")
	_ = ctx
}
