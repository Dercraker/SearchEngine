package Config

import (
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
)

type LimitConfig struct {
	DelayBetweenRequests time.Duration
	MaxPagesPerRun       int
}

func LoadLimitConfig() (limitConfig LimitConfig, err error) {
	delayBetweenRequest, err := configHelper.ParseDuration("CRAWLER_LIMIT_DELAY_BETWEEN_REQUEST", "1s")
	if err != nil {
		return LimitConfig{}, err
	}

	maxPages := configHelper.GetEnvInt("CRAWLER_MAX_REDIRECTS", 5)

	return LimitConfig{
		DelayBetweenRequests: delayBetweenRequest,
		MaxPagesPerRun:       maxPages,
	}, nil
}
