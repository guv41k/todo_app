package statistics_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
)

type statisticsRepositoryStub struct {
	called bool
	userID *int
	from   *time.Time
	to     *time.Time
}

func (r *statisticsRepositoryStub) GetStatistics(_ context.Context, userID *int, from, to *time.Time) (domain.Statistics, error) {
	r.called = true
	r.userID, r.from, r.to = userID, from, to
	return domain.NewStatistics(2, 1, nil), nil
}

func TestGetStatisticsRejectsInvalidFilters(t *testing.T) {
	for _, test := range []struct {
		name   string
		userID *int
		from   *time.Time
		to     *time.Time
	}{
		{name: "invalid user id", userID: intPtr(0)},
		{name: "empty interval", from: timePtr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)), to: timePtr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))},
		{name: "reversed interval", from: timePtr(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)), to: timePtr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))},
	} {
		t.Run(test.name, func(t *testing.T) {
			repository := &statisticsRepositoryStub{}
			_, err := NewStatisticsService(repository).GetStatistics(context.Background(), test.userID, test.from, test.to)
			if !errors.Is(err, core_errors.ErrInvalidArgument) || repository.called {
				t.Fatalf("invalid filters reached repository: %v", err)
			}
		})
	}
}

func TestGetStatisticsPassesFilters(t *testing.T) {
	repository := &statisticsRepositoryStub{}
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 1)
	statistics, err := NewStatisticsService(repository).GetStatistics(context.Background(), intPtr(3), &from, &to)
	if err != nil {
		t.Fatal(err)
	}
	if !repository.called || *repository.userID != 3 || !repository.from.Equal(from) || !repository.to.Equal(to) || statistics.TasksCreated != 2 {
		t.Fatalf("filters or result lost: %+v, %+v", repository, statistics)
	}
}

func intPtr(value int) *int              { return &value }
func timePtr(value time.Time) *time.Time { return &value }
