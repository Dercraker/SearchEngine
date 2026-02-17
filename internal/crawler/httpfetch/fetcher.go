package httpfetch

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/config"
	"github.com/Dercraker/SearchEngine/internal/shared/customErrors"
)

type Fetcher struct {
	config config.FetcherConfig
	client *http.Client
}

func New(config config.FetcherConfig) *Fetcher {
	if config.Timeout <= 0 {
		config.Timeout = 10 * time.Second
	}

	if config.UserAgent == "" {
		config.UserAgent = "Dercraker SearchEngineBot/0.1 (+https://github.com/Dercraker/SearchEngine)"
	}

	if config.MaxBodyBytes <= 0 {
		config.MaxBodyBytes = 2 * 1024 * 1024 // 2MB
	}

	if config.MaxRedirects <= 0 {
		config.MaxRedirects = 5
	}

	c := &http.Client{Timeout: config.Timeout}

	if !config.FollowRedirects {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxRedirects {
				config.Logger.Error("too many redirects", slog.Any("URL", req.URL), slog.Int("Redirect", len(via)))
				return customErrors.ErrTooManyRedirects
			}
			req.Header.Set("User-Agent", config.UserAgent)
			return nil
		}
	}

	return &Fetcher{config: config, client: c}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{}, err
	}

	req.Header.Set("User-Agent", f.config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := f.client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	limit := io.LimitReader(resp.Body, f.config.MaxBodyBytes+1)
	body, err := io.ReadAll(limit)
	if err != nil {
		return Result{}, err
	}

	if int64(len(body)) > f.config.MaxBodyBytes {
		f.config.Logger.Error("Body too large", slog.Any("URL", url), slog.Int64("BodySize", int64(len(body))), slog.Int64("MaxBodySize", f.config.MaxBodyBytes))
		return Result{}, customErrors.ErrBodyTooLarge
	}

	ct := resp.Header.Get("Content-Type")

	return Result{StatusCode: resp.StatusCode, ContentType: ct, Body: body, FinalURL: resp.Request.URL.String()}, nil
}
