package obs

import (
	"sync/atomic"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared"
)

type Stats struct {
	TotalSeeds   int
	InvalidSeeds int
	DedupSkipped int

	Processed atomic.Int64
	Success   atomic.Int64
	Failed    atomic.Int64

	SkippedNonHTML atomic.Int64
	Inserted       atomic.Int64
	Updated        atomic.Int64
	Unchanged      atomic.Int64
	Touched        atomic.Int64

	FetchFailed atomic.Int64
	DBFailed    atomic.Int64
	Retries     atomic.Int64

	StartTime time.Time
	EndTime   time.Time
}

func (s Stats) DurationMs() float64 {
	return shared.DurationMs(s.StartTime, s.EndTime)
}
