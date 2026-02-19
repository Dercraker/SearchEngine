package shared

import (
	"time"
)

func DurationMs(start, end time.Time) float64 {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	duration := float64(end.Sub(start)) / float64(time.Millisecond)
	return duration
}
func DurationS(start, end time.Time) float64 {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	duration := float64(end.Sub(start)) / float64(time.Second)
	return duration
}
func DurationM(start, end time.Time) float64 {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	duration := float64(end.Sub(start)) / float64(time.Minute)
	return duration
}
