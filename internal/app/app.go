package app

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"Git_PR_plugin_avito_tech/internal/controller"
	"Git_PR_plugin_avito_tech/internal/repository"
	"Git_PR_plugin_avito_tech/internal/service"
)

func Run(db *sqlx.DB) {
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prReviewerRepo := repository.NewPrReviewerRepository(db)
	pullRequestRepo := repository.NewPullRequestRepository(db)

	teamService := service.NewTeamService(teamRepo, userRepo, db)
	userService := service.NewUserService(userRepo, teamRepo)
	pullRequestService := service.NewPullRequestService(prReviewerRepo, pullRequestRepo, userRepo, db)

	teamController := controller.NewTeamController(teamService)
	userController := controller.NewUserController(userService, pullRequestService)
	pullRequestController := controller.NewPullRequestController(pullRequestService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/team", func(r chi.Router) {
		r.Post("/add", teamController.CreateTeamWithMembers)
		r.Get("/get", teamController.GetTeamByName)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", userController.SetIsActive)
		r.Get("/getReview", userController.GetPullRequestsByReviewerID)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", pullRequestController.CreatePullRequest)
		r.Post("/merge", pullRequestController.MergePullRequest)
		r.Post("/reassign", pullRequestController.ReassignPrReviewer)
	})

	_ = teamService

	log.Println("Application started successfully on port 8080")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
