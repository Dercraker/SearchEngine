package Config

import (
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
)

type LimitConfig struct {
	MaxPagesPerRun int64

	GlobalDelay  time.Duration
	GlobalJitter time.Duration

	HostDelay time.Duration
	Jitter    time.Duration
	MaxHost   int
}

func LoadLimitConfig() (limitConfig LimitConfig, err error) {
	maxPages := configHelper.GetEnvInt64("CRAWLER_LIMIT_MAX_PAGES_PER_RUN", 200)

	globalDelay, err := configHelper.ParseDuration("CRAWLER_LIMIT_GLOBAL_DELAY", "200ms")
	if err != nil {
		return LimitConfig{}, err
	}

	globalJitter, err := configHelper.ParseDuration("CRAWLER_LIMIT_GLOBAL_JITTER", "100ms")
	if err != nil {
		return LimitConfig{}, err
	}

	hostDelay, err := configHelper.ParseDuration("CRAWLER_LIMIT_HOST_DELAY", "500ms")
	if err != nil {
		return LimitConfig{}, err
	}

	jitter, err := configHelper.ParseDuration("CRAWLER_LIMIT_JITTER", "250ms")
	if err != nil {
		return LimitConfig{}, err
	}

	maxHost := configHelper.GetEnvInt("CRAWLER_LIMIT_MAX_HOST", 10_000)

	return LimitConfig{
		MaxPagesPerRun: maxPages,
		GlobalDelay:    globalDelay,
		GlobalJitter:   globalJitter,
		HostDelay:      hostDelay,
		Jitter:         jitter,
		MaxHost:        maxHost,
	}, nil
}
