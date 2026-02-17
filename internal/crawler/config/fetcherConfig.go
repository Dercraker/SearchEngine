package config

import (
	"log/slog"
	"time"

	"github.com/Dercraker/SearchEngine/internal/shared/configHelper"
)

type FetcherConfig struct {
	Timeout         time.Duration
	UserAgent       string
	MaxBodyBytes    int64
	FollowRedirects bool
	MaxRedirects    int
	Logger          *slog.Logger
}

func LoadFetcherConfig() (fetcherConfig FetcherConfig, err error) {
	readTimeout, err := configHelper.ParseDuration("CRAWLER_READ_TIMEOUT", "30s")
	if err != nil {
		return FetcherConfig{}, err
	}

	userAgent := configHelper.GetEnv("CRAWLER_USER_AGENT", "")

	maxBodyBytes := configHelper.GetEnvInt64("CRAWLER_MAX_BODY_BYTES", 2*1024*1024)

	maxRedirects := configHelper.GetEnvInt("CRAWLER_MAX_REDIRECTS", 5)

	return FetcherConfig{
		Timeout:         readTimeout,
		UserAgent:       userAgent,
		MaxBodyBytes:    maxBodyBytes,
		FollowRedirects: true,
		MaxRedirects:    maxRedirects,
	}, nil
}
