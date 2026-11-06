package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	if userID != nil && *userID <= 0 {
		return domain.Statistics{}, fmt.Errorf("user_id must be positive: %w", core_errors.ErrInvalidArgument)
	}
	if from != nil && to != nil && !to.After(*from) {
		return domain.Statistics{}, fmt.Errorf("to must be after from: %w", core_errors.ErrInvalidArgument)
	}

	statistics, err := s.repository.GetStatistics(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get statistics: %w", err)
	}
	return statistics, nil
}
