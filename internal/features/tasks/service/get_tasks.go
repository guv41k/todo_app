package tasks_service

import (
	"context"
	"fmt"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
)

func (s *TasksService) GetTasks(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error) {
	if userID != nil && *userID <= 0 {
		return nil, fmt.Errorf("user_id must be positive: %w", core_errors.ErrInvalidArgument)
	}
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}
	return tasks, nil
}
