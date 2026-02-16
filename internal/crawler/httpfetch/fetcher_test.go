// internal/crawler/httpfetch/fetcher_test.go
package httpfetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Dercraker/SearchEngine/internal/crawler/config"
)

func TestFetch_HTML_ReturnsBodyAndMeta(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); !strings.Contains(ua, "TestBot") {
			t.Fatalf("expected UA to contain TestBot, got %q", ua)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html>ok</html>"))
	}))
	defer srv.Close()

	f := New(Config.FetcherConfig{
		Timeout:         2 * time.Second,
		UserAgent:       "TestBot/1.0",
		FollowRedirects: true,
	})

	res, err := f.Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if !strings.Contains(strings.ToLower(res.ContentType), "text/html") {
		t.Fatalf("expected html content-type, got %q", res.ContentType)
	}
	if string(res.Body) != "<html>ok</html>" {
		t.Fatalf("unexpected body: %q", string(res.Body))
	}
}

func TestFetch_Redirect_Refused(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("target"))
	}))
	defer target.Close()

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer redirector.Close()

	f := New(Config.FetcherConfig{
		Timeout:         2 * time.Second,
		UserAgent:       "TestBot/1.0",
		FollowRedirects: false,
	})

	res, err := f.Fetch(context.Background(), redirector.URL)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.StatusCode != http.StatusFound {
		t.Fatalf("expected 302, got %d", res.StatusCode)
	}
}
