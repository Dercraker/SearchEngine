package processors

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/httpfetch"
)

type Downloader struct {
	Fetcher *httpfetch.Fetcher
	Store   DocumentStore
}

func (d Downloader) Process(ctx context.Context, u *url.URL) error {
	perURLCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := d.Fetcher.Fetch(perURLCtx, u.String())
	if err != nil {
		return err
	}

	if !strings.Contains(strings.ToLower(res.ContentType), "text/html") {
		return nil
	}

	//res.body == html
	//ICI on parse / extract / store

	hash := sha256Hex(res.Body)
	finalUrl := res.FinalURL

	oldHash, err := d.Store.GetHashByURL(ctx, finalUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return d.Store.UpsertFetch(ctx, res.FinalURL, res.StatusCode, res.ContentType, hash, res.Body)
		}
		return err
	}

	if oldHash == hash {
		return d.Store.TouchFetchAt(ctx, finalUrl, res.StatusCode, res.ContentType)
	}

	return d.Store.UpsertFetch(ctx, res.FinalURL, res.StatusCode, res.ContentType, hash, res.Body)
}
