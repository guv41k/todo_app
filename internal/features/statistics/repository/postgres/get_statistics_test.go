package statistics_postgres_repository

import (
	"context"
	"testing"
	"time"

	core_postgres_pool "github.com/guv41k/todo_app/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
)

type statisticsPoolStub struct {
	core_postgres_pool.Pool
	args []any
	row  pgx.Row
}

func (p *statisticsPoolStub) QueryRow(_ context.Context, _ string, args ...any) pgx.Row {
	p.args = args
	return p.row
}

func (p *statisticsPoolStub) OpTimeout() time.Duration { return time.Second }

type statisticsRowStub struct {
	created, completed int64
	average            *float64
}

func (r statisticsRowStub) Scan(dest ...any) error {
	*dest[0].(*int64) = r.created
	*dest[1].(*int64) = r.completed
	*dest[2].(**float64) = r.average
	return nil
}

func TestGetStatisticsMapsAggregate(t *testing.T) {
	seconds := 90.0
	pool := &statisticsPoolStub{row: statisticsRowStub{created: 4, completed: 3, average: &seconds}}
	userID := 2
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 1)
	statistics, err := NewStatisticsRepository(pool).GetStatistics(context.Background(), &userID, &from, &to)
	if err != nil {
		t.Fatal(err)
	}
	if len(pool.args) != 3 || pool.args[0] != &userID || pool.args[1] != &from || pool.args[2] != &to {
		t.Fatalf("filters lost: %v", pool.args)
	}
	if statistics.TasksCreated != 4 || statistics.TasksCompleted != 3 || statistics.TasksCompletedRate == nil || *statistics.TasksCompletedRate != 75 || statistics.TasksAverageCompletionTime == nil || *statistics.TasksAverageCompletionTime != 90*time.Second {
		t.Fatalf("incorrect statistics: %+v", statistics)
	}
}

func TestGetStatisticsEmptyAggregate(t *testing.T) {
	pool := &statisticsPoolStub{row: statisticsRowStub{}}
	statistics, err := NewStatisticsRepository(pool).GetStatistics(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if statistics.TasksCreated != 0 || statistics.TasksCompleted != 0 || statistics.TasksCompletedRate != nil || statistics.TasksAverageCompletionTime != nil {
		t.Fatalf("incorrect empty statistics: %+v", statistics)
	}
}
