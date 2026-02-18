package processors

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/httpfetch"
	"github.com/Dercraker/SearchEngine/internal/crawler/obs"
)

type Downloader struct {
	Fetcher *httpfetch.Fetcher
	Store   DocumentStore
	Stats   *obs.Stats
}

func (d Downloader) Process(ctx context.Context, u *url.URL) error {
	perURLCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := d.Fetcher.Fetch(perURLCtx, u.String())
	if err != nil {
		d.Stats.FetchFailed.Add(1)
		return err
	}

	if !strings.Contains(strings.ToLower(res.ContentType), "text/html") {
		d.Stats.SkippedNonHTML.Add(1)
		return nil
	}

	//res.body == html
	//ICI on parse / extract / store

	hash := sha256Hex(res.Body)
	finalUrl := res.FinalURL

	oldHash, err := d.Store.GetHashByURL(ctx, finalUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := d.Store.UpsertFetch(ctx, res.FinalURL, res.StatusCode, res.ContentType, hash, res.Body); err != nil {
				d.Stats.DBFailed.Add(1)
				return err
			}
			d.Stats.Inserted.Add(1)
		}
		d.Stats.DBFailed.Add(1)
		return err
	}

	if oldHash == hash {
		if err := d.Store.TouchFetchAt(ctx, finalUrl, res.StatusCode, res.ContentType); err != nil {
			d.Stats.DBFailed.Add(1)
			return err
		}
		d.Stats.Touched.Add(1)
		d.Stats.Unchanged.Add(1)
		return nil
	}

	if err := d.Store.UpsertFetch(ctx, res.FinalURL, res.StatusCode, res.ContentType, hash, res.Body); err != nil {
		d.Stats.DBFailed.Add(1)
		return err
	}

	d.Stats.Updated.Add(1)
	return nil
}
