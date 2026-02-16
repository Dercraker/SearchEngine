package processors

import (
	"context"
	"net/url"
	"strings"

	"github.com/Dercraker/SearchEngine/internal/crawler/httpfetch"
)

type Downloader struct {
	Fetcher *httpfetch.Fetcher
	Store   DocumentStore
}

func (d Downloader) Process(ctx context.Context, u *url.URL) error {

	res, err := d.Fetcher.Fetch(ctx, u.String())
	if err != nil {
		return err
	}

	if !strings.Contains(strings.ToLower(res.ContentType), "text/html") {
		return nil
	}

	//res.body == html
	//ICI on parse / extract / store

	hash := sha256Hex(res.Body)

	if err := d.Store.UpsertFetch(ctx, res.FinalURL, res.StatusCode, res.ContentType, hash, res.Body); err != nil {
		return err
	}

	return nil
}
