package crawler

import "time"

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

func (s Stats) Duration() time.Duration {
	if s.EndTime.IsZero() || s.StartTime.IsZero() {
		return 0
	}
	return s.EndTime.Sub(s.StartTime)
}
