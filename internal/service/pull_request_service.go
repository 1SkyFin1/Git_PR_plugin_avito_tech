package service

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/dto"
	"Git_PR_plugin_avito_tech/internal/model"
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PullRequestService struct {
	prReviewerRepo  PrReviewerRepository
	pullRequestRepo PullRequestRepository
	userRepo        UserRepository
	db              *sqlx.DB
}

func NewPullRequestService(prReviewerRepo PrReviewerRepository, pullRequestRepo PullRequestRepository, userRepo UserRepository, db *sqlx.DB) *PullRequestService {
	return &PullRequestService{
		prReviewerRepo:  prReviewerRepo,
		pullRequestRepo: pullRequestRepo,
		userRepo:        userRepo,
		db:              db,
	}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, request *dto.CreatePullRequestRequest) (*dto.CreatePullRequestResponse, error) {
	log.Printf("[PullRequestService.CreatePullRequest] Starting for PR: %s, Author: %s", request.PullRequestID, request.AuthorID)
	
	author, err := s.userRepo.FindByID(ctx, request.AuthorID)
	if err != nil {
		log.Printf("[PullRequestService.CreatePullRequest] Error finding author: %v", err)
		if apperror.IsNotFoundError(err) {
			return nil, apperror.NewNotFoundError("Author")
		}
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[PullRequestService.CreatePullRequest] Author found: %s, Team: %s", author.Username, author.TeamID)

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[PullRequestService.CreatePullRequest] Error beginning transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer tx.Rollback()

	pr := dto.ToModelFromCreatePullRequestRequest(request)
	err = s.pullRequestRepo.CreateTx(ctx, tx, pr)
	if err != nil {
		log.Printf("[PullRequestService.CreatePullRequest] Error creating PR: %v", err)
		return nil, apperror.HandleDBError(err, "Pull Request")
	}

	log.Printf("[PullRequestService.CreatePullRequest] PR created successfully, finding reviewers...")

	teamUsers, err := s.userRepo.FindAllByTeamID(ctx, author.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team users: %w", err)
	}

	var potentialReviewers []*model.User
	for _, user := range teamUsers {
		if user.UserID != request.AuthorID && user.IsActive {
			potentialReviewers = append(potentialReviewers, user)
		}
	}

	if len(potentialReviewers) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
		return dto.ToCreatePullRequestResponse(pr, []uuid.UUID{}), nil
	}

	reviewerIDs := make([]uuid.UUID, 0, len(potentialReviewers))
	for _, reviewer := range potentialReviewers {
		reviewerIDs = append(reviewerIDs, reviewer.UserID)
	}

	reviewerCounts, err := s.prReviewerRepo.CountByReviewerIDs(ctx, reviewerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to count reviewers: %w", err)
	}

	type reviewerWithCount struct {
		user  *model.User
		count int
	}
	reviewersWithCounts := make([]reviewerWithCount, 0, len(potentialReviewers))
	for _, reviewer := range potentialReviewers {
		count := reviewerCounts[reviewer.UserID]
		reviewersWithCounts = append(reviewersWithCounts, reviewerWithCount{
			user:  reviewer,
			count: count,
		})
	}

	for i := 0; i < len(reviewersWithCounts)-1; i++ {
		for j := i + 1; j < len(reviewersWithCounts); j++ {
			if reviewersWithCounts[j].count < reviewersWithCounts[i].count {
				reviewersWithCounts[i], reviewersWithCounts[j] = reviewersWithCounts[j], reviewersWithCounts[i]
			}
		}
	}

	assignedReviewerIDs := make([]uuid.UUID, 0, 2)
	maxReviewers := 2
	if len(reviewersWithCounts) < maxReviewers {
		maxReviewers = len(reviewersWithCounts)
	}

	for i := 0; i < maxReviewers; i++ {
		prReviewer := &model.PrReviewer{
			PullRequestId: pr.PullRequestID,
			ReviewerId:    reviewersWithCounts[i].user.UserID,
		}
		err = s.prReviewerRepo.CreateTx(ctx, tx, prReviewer)
		if err != nil {
			return nil, fmt.Errorf("failed to assign reviewer: %w", err)
		}
		assignedReviewerIDs = append(assignedReviewerIDs, prReviewer.ReviewerId)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[PullRequestService.CreatePullRequest] Error committing transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to commit transaction: %w", err))
	}

	log.Printf("[PullRequestService.CreatePullRequest] Success! PR: %s, Assigned reviewers: %v", pr.PullRequestID, assignedReviewerIDs)
	return dto.ToCreatePullRequestResponse(pr, assignedReviewerIDs), nil
}

func (s *PullRequestService) GetPullRequestsByReviewerID(ctx context.Context, reviewerID uuid.UUID) (*dto.GetPullRequestsByReviewerIDResponse, error) {
	log.Printf("[PullRequestService.GetPullRequestsByReviewerID] Starting for reviewerID: %s", reviewerID)
	
	prReviewers, err := s.prReviewerRepo.FindAllByReviewerID(ctx, reviewerID)
	if err != nil {
		log.Printf("[PullRequestService.GetPullRequestsByReviewerID] Error finding pr_reviewers: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[PullRequestService.GetPullRequestsByReviewerID] Found %d pr_reviewers", len(prReviewers))

	// Если пользователь не является ревьювером ни на одном PR, возвращаем пустой список
	if len(prReviewers) == 0 {
		log.Printf("[PullRequestService.GetPullRequestsByReviewerID] No PRs found for reviewer, returning empty list")
		return dto.ToGetPullRequestsByReviewerIDResponse(reviewerID, []*model.PullRequest{}), nil
	}

	prReviewerIDs := make([]uuid.UUID, 0, len(prReviewers))
	for _, prReviewer := range prReviewers {
		prReviewerIDs = append(prReviewerIDs, prReviewer.PullRequestId)
	}

	pullRequests, err := s.pullRequestRepo.FindAllByPullRequestIDIn(ctx, prReviewerIDs)
	if err != nil {
		log.Printf("[PullRequestService.GetPullRequestsByReviewerID] Error finding pull requests: %v", err)
		return nil, apperror.NewInternalError(err)
	}
	
	log.Printf("[PullRequestService.GetPullRequestsByReviewerID] Success! Found %d pull requests", len(pullRequests))
	return dto.ToGetPullRequestsByReviewerIDResponse(reviewerID, pullRequests), nil
}

func (s *PullRequestService) MergePullRequest(ctx context.Context, request *dto.MergePullRequestRequest) (*dto.MergePullRequestResponse, error) {
	log.Printf("[PullRequestService.MergePullRequest] Starting for PR: %s", request.PullRequestID)

	pr, err := s.pullRequestRepo.MergePullRequest(ctx, request.PullRequestID)
	if err != nil {
		log.Printf("[PullRequestService.MergePullRequest] Error merging PR: %v", err)
		if apperror.IsNotFoundError(err) {
			return nil, apperror.NewNotFoundError("Pull Request")
		}
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[PullRequestService.MergePullRequest] PR merged: %s, fetching reviewers...", pr.PullRequestName)

	prReviewers, err := s.prReviewerRepo.FindAllByPRID(ctx, pr.PullRequestID)
	if err != nil {
		log.Printf("[PullRequestService.MergePullRequest] Error fetching reviewers: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	assignedReviewerIDs := make([]uuid.UUID, 0, 2)
	for _, prReviewer := range prReviewers {
		assignedReviewerIDs = append(assignedReviewerIDs, prReviewer.ReviewerId)
	}

	log.Printf("[PullRequestService.MergePullRequest] Success! PR: %s, Reviewers: %d", pr.PullRequestID, len(assignedReviewerIDs))
	return dto.ToMergePullRequestResponse(pr, assignedReviewerIDs), nil
}

func (s *PullRequestService) ReassignPrReviewer(ctx context.Context, request *dto.ReassignPrReviewerRequest) (*dto.ReassignPrReviewerResponse, error) {
	log.Printf("[PullRequestService.ReassignPrReviewer] Starting for PR: %s, OldReviewer: %s", request.PullRequestID, request.OldReviewerID)
	
	pr, err := s.pullRequestRepo.FindByID(ctx, request.PullRequestID)
	if err != nil {
		log.Printf("[PullRequestService.ReassignPrReviewer] Error finding PR: %v", err)
		if apperror.IsNotFoundError(err) {
			return nil, apperror.NewNotFoundError("Pull Request")
		}
		return nil, apperror.NewInternalError(err)
	}

	// Проверяем, что PR не смержен
	if pr.Status == model.PullRequestStatusMerged {
		log.Printf("[PullRequestService.ReassignPrReviewer] Cannot reassign reviewers: PR is already merged")
		return nil, apperror.NewPRMergedError()
	}

	log.Printf("[PullRequestService.ReassignPrReviewer] PR found: %s, Status: %s", pr.PullRequestName, pr.Status)

	oldReviewer, err := s.userRepo.FindByID(ctx, request.OldReviewerID)
	if err != nil {
		log.Printf("[PullRequestService.ReassignPrReviewer] Error finding old reviewer: %v", err)
		if apperror.IsNotFoundError(err) {
			return nil, apperror.NewNotFoundError("Reviewer")
		}
		return nil, apperror.NewInternalError(err)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	currentReviewers, err := s.prReviewerRepo.FindAllByPRID(ctx, request.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current reviewers: %w", err)
	}

	currentReviewerIDs := make(map[uuid.UUID]bool)
	for _, reviewer := range currentReviewers {
		currentReviewerIDs[reviewer.ReviewerId] = true
	}

	if !currentReviewerIDs[request.OldReviewerID] {
		return nil, fmt.Errorf("old reviewer is not assigned to this pull request")
	}

	teamUsers, err := s.userRepo.FindAllByTeamID(ctx, oldReviewer.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team users: %w", err)
	}

	var potentialReviewers []*model.User
	for _, user := range teamUsers {
		if user.UserID != pr.AuthorID &&
			!currentReviewerIDs[user.UserID] &&
			user.IsActive {
			potentialReviewers = append(potentialReviewers, user)
		}
	}

	if len(potentialReviewers) == 0 {
		return nil, fmt.Errorf("no available reviewers to reassign")
	}

	randomIndex := rand.Intn(len(potentialReviewers))
	newReviewer := potentialReviewers[randomIndex]

	err = s.prReviewerRepo.DeleteByPRIDAndReviewerIDTx(ctx, tx, request.PullRequestID, request.OldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old reviewer: %w", err)
	}

	prReviewer := &model.PrReviewer{
		PullRequestId: request.PullRequestID,
		ReviewerId:    newReviewer.UserID,
	}
	err = s.prReviewerRepo.CreateTx(ctx, tx, prReviewer)
	if err != nil {
		return nil, fmt.Errorf("failed to assign new reviewer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[PullRequestService.ReassignPrReviewer] Error committing transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to commit transaction: %w", err))
	}

	updatedReviewers, err := s.prReviewerRepo.FindAllByPRID(ctx, request.PullRequestID)
	if err != nil {
		log.Printf("[PullRequestService.ReassignPrReviewer] Error getting updated reviewers: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	assignedReviewerIDs := make([]uuid.UUID, 0, len(updatedReviewers))
	for _, reviewer := range updatedReviewers {
		assignedReviewerIDs = append(assignedReviewerIDs, reviewer.ReviewerId)
	}

	log.Printf("[PullRequestService.ReassignPrReviewer] Success! PR: %s, Old: %s, New: %s", request.PullRequestID, request.OldReviewerID, newReviewer.UserID)
	return dto.ToReassignPrReviewerResponse(pr, assignedReviewerIDs, newReviewer.UserID), nil
}
