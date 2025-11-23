package service

import (
	"Git_PR_plugin_avito_tech/internal/apperror"
	"Git_PR_plugin_avito_tech/internal/dto"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StatisticsService struct {
	userRepo       UserRepository
	prRepo         PullRequestRepository
	prReviewerRepo PrReviewerRepository
	db             *sqlx.DB
}

func NewStatisticsService(
	userRepo UserRepository,
	prRepo PullRequestRepository,
	prReviewerRepo PrReviewerRepository,
	db *sqlx.DB,
) *StatisticsService {
	return &StatisticsService{
		userRepo:       userRepo,
		prRepo:         prRepo,
		prReviewerRepo: prReviewerRepo,
		db:             db,
	}
}

func (s *StatisticsService) GetStatistics(ctx context.Context) (*dto.GetStatisticsResponse, error) {
	log.Printf("[StatisticsService.GetStatistics] Starting")

	totalUsers, err := s.userRepo.GetTotalUsers(ctx)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting total users: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	activeUsers, err := s.userRepo.GetActiveUsers(ctx)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting active users: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	totalPRs, err := s.prRepo.GetTotalPullRequests(ctx)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting total PRs: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	openPRs, err := s.prRepo.GetOpenPullRequests(ctx)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting open PRs: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	mergedPRs, err := s.prRepo.GetMergedPullRequests(ctx)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting merged PRs: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	totalAssignments, err := s.prReviewerRepo.GetTotalAssignments(ctx)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting total assignments: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	userStatsData, err := s.prReviewerRepo.GetUserStatisticsData(ctx, 10)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting user statistics: %v", err)
		return nil, apperror.NewInternalError(err)
	}
	topReviewers := dto.ToUserStatisticsList(userStatsData)

	prStatsData, err := s.prReviewerRepo.GetPRStatisticsData(ctx, 10)
	if err != nil {
		log.Printf("[StatisticsService.GetStatistics] Error getting PR statistics: %v", err)
		return nil, apperror.NewInternalError(err)
	}
	prWithMostReviewers := dto.ToPullRequestStatisticsList(prStatsData)

	response := &dto.GetStatisticsResponse{
		TotalUsers:                    totalUsers,
		ActiveUsers:                   activeUsers,
		TotalPullRequests:             totalPRs,
		OpenPullRequests:              openPRs,
		MergedPullRequests:            mergedPRs,
		TotalAssignments:              totalAssignments,
		TopReviewers:                  topReviewers,
		PullRequestsWithMostReviewers: prWithMostReviewers,
	}

	log.Printf("[StatisticsService.GetStatistics] Success")
	return response, nil
}

func (s *StatisticsService) DeactivateUsers(ctx context.Context, req *dto.DeactivateUsersRequest) (*dto.DeactivateUsersResponse, error) {
	log.Printf("[StatisticsService.DeactivateUsers] Starting for %d users", len(req.UserIDs))

	if len(req.UserIDs) == 0 {
		return &dto.DeactivateUsersResponse{
			DeactivatedCount:   0,
			ReassignedPRsCount: 0,
			DeactivatedUserIDs: []uuid.UUID{},
		}, nil
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Printf("[StatisticsService.DeactivateUsers] Error starting transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to start transaction: %w", err))
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[StatisticsService.DeactivateUsers] Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	deactivatedIDs, err := s.userRepo.DeactivateUsersTx(ctx, tx, req.UserIDs)
	if err != nil {
		log.Printf("[StatisticsService.DeactivateUsers] Error deactivating users: %v", err)
		return nil, apperror.NewInternalError(err)
	}

	log.Printf("[StatisticsService.DeactivateUsers] Deactivated %d users", len(deactivatedIDs))

	reassignedCount := 0
	if len(deactivatedIDs) > 0 {
		reassignedCount, err = s.prReviewerRepo.DeleteByReviewerIDsTx(ctx, tx, deactivatedIDs)
		if err != nil {
			log.Printf("[StatisticsService.DeactivateUsers] Error deleting assignments: %v", err)
			return nil, apperror.NewInternalError(err)
		}
		log.Printf("[StatisticsService.DeactivateUsers] Removed %d assignments", reassignedCount)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[StatisticsService.DeactivateUsers] Error committing transaction: %v", err)
		return nil, apperror.NewInternalError(fmt.Errorf("failed to commit transaction: %w", err))
	}

	response := &dto.DeactivateUsersResponse{
		DeactivatedCount:   len(deactivatedIDs),
		ReassignedPRsCount: reassignedCount,
		DeactivatedUserIDs: deactivatedIDs,
	}

	log.Printf("[StatisticsService.DeactivateUsers] Success: deactivated %d users, removed %d assignments",
		len(deactivatedIDs), reassignedCount)

	return response, nil
}
