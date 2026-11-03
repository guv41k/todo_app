package tasks_service

import (
	"context"
	"fmt"

	"github.com/guv41k/todo_app/internal/core/domain"
)

func (s *TasksService) PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}
	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}
	updated, err := s.tasksRepository.UpdateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}
	return updated, nil
}
