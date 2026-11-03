package tasks_service

import (
	"context"

	"github.com/guv41k/todo_app/internal/core/domain"
)

type TasksRepository interface {
	CreateTask(context.Context, domain.Task) (domain.Task, error)
	GetTasks(context.Context, *int, *int, *int) ([]domain.Task, error)
	GetTask(context.Context, int) (domain.Task, error)
	UpdateTask(context.Context, domain.Task) (domain.Task, error)
	DeleteTask(context.Context, int) error
}

type TasksService struct {
	tasksRepository TasksRepository
}

func NewTasksService(repository TasksRepository) *TasksService {
	return &TasksService{tasksRepository: repository}
}
