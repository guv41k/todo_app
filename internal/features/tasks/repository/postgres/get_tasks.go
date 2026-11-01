package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/guv41k/todo_app/internal/core/domain"
)

func (r *TasksRepository) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE $1::INTEGER IS NULL OR author_user_id=$1
	ORDER BY id ASC
	LIMIT $2
	OFFSET $3;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel
		if err := taskModel.Scan(rows); err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return modelsToDomains(taskModels), nil
}
