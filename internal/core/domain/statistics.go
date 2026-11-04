package domain

import "time"

type Statistics struct {
	TasksCreated               int64
	TasksCompleted             int64
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(created, completed int64, average *time.Duration) Statistics {
	statistics := Statistics{
		TasksCreated:               created,
		TasksCompleted:             completed,
		TasksAverageCompletionTime: average,
	}
	if created > 0 {
		rate := float64(completed) / float64(created) * 100
		statistics.TasksCompletedRate = &rate
	}
	return statistics
}
