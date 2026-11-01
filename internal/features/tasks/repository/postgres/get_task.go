package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *TasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel
	if err := taskModel.Scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d': %w", id, core_errors.ErrNotFound)
		}

		return domain.Task{}, fmt.Errorf("scan task: %w", err)
	}

	return modelToDomain(taskModel), nil
}
