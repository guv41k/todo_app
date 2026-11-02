package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
	core_postgres_pool "github.com/guv41k/todo_app/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)

	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if core_postgres_pool.IsForeignKeyViolation(err) {
			return domain.Task{}, fmt.Errorf("user with id='%d': %w", task.AuthorUserID, core_errors.ErrNotFound)
		}

		return domain.Task{}, fmt.Errorf("scan created task: %w", err)
	}

	return modelToDomain(taskModel), nil
}
