package statistics_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/guv41k/todo_app/internal/core/domain"
)

func (r *StatisticsRepository) GetStatistics(ctx context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
	SELECT
		COUNT(*),
		COUNT(*) FILTER (WHERE completed),
		EXTRACT(EPOCH FROM AVG(completed_at - created_at) FILTER (WHERE completed))::double precision
	FROM todoapp.tasks
	WHERE ($1::integer IS NULL OR author_user_id = $1)
		AND ($2::timestamptz IS NULL OR created_at >= $2)
		AND ($3::timestamptz IS NULL OR created_at < $3);
	`

	var created, completed int64
	var averageSeconds *float64
	if err := r.pool.QueryRow(ctx, query, userID, from, to).Scan(&created, &completed, &averageSeconds); err != nil {
		return domain.Statistics{}, fmt.Errorf("select task statistics: %w", err)
	}

	var average *time.Duration
	if averageSeconds != nil {
		duration := time.Duration(*averageSeconds * float64(time.Second))
		average = &duration
	}
	return domain.NewStatistics(created, completed, average), nil
}
