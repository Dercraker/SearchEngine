package shared

import (
	"math"
	"time"
)

func DurationMs(start, end time.Time) float64 {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	duration := float64(end.Sub(start)) / float64(time.Millisecond)
	duration = math.Round(duration*1000) / 1000
	return duration
}
