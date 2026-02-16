package crawler

import (
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared"
)

type Stats struct {
	TotalSeeds   int
	InvalidSeeds int
	DedupSkipped int
	Processed    int
	Success      int
	Failed       int
	StartTime    time.Time
	EndTime      time.Time
}

func (s Stats) DurationMs() float64 {
	return shared.DurationMs(s.StartTime, s.EndTime)
}
