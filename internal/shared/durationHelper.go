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
	duration = math.Round(duration*1) / 1
	return duration
}
func DurationS(start, end time.Time) float64 {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	duration := float64(end.Sub(start)) / float64(time.Second)
	duration = math.Round(duration*1) / 1
	return duration
}
func DurationM(start, end time.Time) float64 {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	duration := float64(end.Sub(start)) / float64(time.Minute)
	duration = math.Round(duration*1) / 1
	return duration
}
