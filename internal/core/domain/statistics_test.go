package domain

import (
	"math"
	"testing"
	"time"
)

func TestNewStatistics(t *testing.T) {
	if empty := NewStatistics(0, 0, nil); empty.TasksCompletedRate != nil || empty.TasksAverageCompletionTime != nil {
		t.Fatalf("empty statistics should have no rate or average: %+v", empty)
	}

	average := 90 * time.Second
	statistics := NewStatistics(4, 3, &average)
	if statistics.TasksCompletedRate == nil || math.Abs(*statistics.TasksCompletedRate-75) > 0.000001 {
		t.Fatalf("completion rate = %v, want 75", statistics.TasksCompletedRate)
	}
	if statistics.TasksAverageCompletionTime == nil || *statistics.TasksAverageCompletionTime != average {
		t.Fatalf("average completion time = %v", statistics.TasksAverageCompletionTime)
	}
}
