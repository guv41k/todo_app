package statistics_service

import (
	"context"
	"time"

	"github.com/guv41k/todo_app/internal/core/domain"
)

type StatisticsRepository interface {
	GetStatistics(context.Context, *int, *time.Time, *time.Time) (domain.Statistics, error)
}

type StatisticsService struct {
	repository StatisticsRepository
}

func NewStatisticsService(repository StatisticsRepository) *StatisticsService {
	return &StatisticsService{repository: repository}
}
