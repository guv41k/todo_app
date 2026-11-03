package tasks_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
)

type taskRepositoryStub struct {
	created  bool
	updated  domain.Task
	getTasks bool
	task     *domain.Task
}

func (r *taskRepositoryStub) CreateTask(_ context.Context, task domain.Task) (domain.Task, error) {
	r.created = true
	return task, nil
}
func (r *taskRepositoryStub) GetTasks(_ context.Context, _, _, _ *int) ([]domain.Task, error) {
	r.getTasks = true
	return []domain.Task{}, nil
}
func (r *taskRepositoryStub) GetTask(_ context.Context, _ int) (domain.Task, error) {
	if r.task != nil {
		return *r.task, nil
	}
	return domain.NewTask(1, 1, "task", nil, false, time.Now().Add(-time.Minute), nil, 1), nil
}
func (r *taskRepositoryStub) UpdateTask(_ context.Context, task domain.Task) (domain.Task, error) {
	r.updated = task
	return task, nil
}
func (r *taskRepositoryStub) DeleteTask(_ context.Context, _ int) error { return nil }

func TestCreateTaskValidatesBeforeSaving(t *testing.T) {
	repo := &taskRepositoryStub{}
	service := NewTasksService(repo)
	_, err := service.CreateTask(context.Background(), domain.NewTaskUninitialized("", nil, 1))
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
	if repo.created {
		t.Fatal("invalid task was saved")
	}
}

func TestPatchTaskSetsAndClearsCompletionTime(t *testing.T) {
	repo := &taskRepositoryStub{}
	service := NewTasksService(repo)
	completed := true
	_, err := service.PatchTask(context.Background(), 1, domain.NewTaskPatch(
		domain.Nullable[string]{}, domain.Nullable[string]{},
		domain.Nullable[bool]{Set: true, Value: &completed},
	))
	if err != nil {
		t.Fatal(err)
	}
	if !repo.updated.Completed || repo.updated.CompletedAt == nil || repo.updated.CompletedAt.Before(repo.updated.CreatedAt) {
		t.Fatalf("completion state is inconsistent: %+v", repo.updated)
	}

	completed = false
	previous := repo.updated
	repo.task = &previous
	_, err = service.PatchTask(context.Background(), 1, domain.NewTaskPatch(
		domain.Nullable[string]{}, domain.Nullable[string]{},
		domain.Nullable[bool]{Set: true, Value: &completed},
	))
	if err != nil {
		t.Fatal(err)
	}
	if repo.updated.Completed || repo.updated.CompletedAt != nil {
		t.Fatalf("completion time was not cleared: %+v", repo.updated)
	}
}

func TestGetTasksRejectsInvalidFilters(t *testing.T) {
	repo := &taskRepositoryStub{}
	service := NewTasksService(repo)
	invalid := -1
	_, err := service.GetTasks(context.Background(), nil, &invalid, nil)
	if !errors.Is(err, core_errors.ErrInvalidArgument) || repo.getTasks {
		t.Fatalf("invalid limit reached repository: %v", err)
	}
}
