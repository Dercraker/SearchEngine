package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dercraker/SearchEngine/internal/api/DTO"
)

type fakeSearchService struct {
	lastQ                 string
	lastLimit, lastOffset int

	results []DTO.SearchDto
	err     error
}

func (f *fakeSearchService) Search(_ any, r *http.Request, q string, limit, offset int) ([]DTO.SearchDto, error) {
	f.lastQ = q
	f.lastLimit = limit
	f.lastOffset = offset

	return f.results, f.err
}

func TestSearchHandler_Return200AndEmptyArray_WithHeaders(t *testing.T) {
	fs := &fakeSearchService{results: nil, err: nil}
	h := NewSearchHandler(fs, 10, 50)

	req := httptest.NewRequest(http.MethodGet, "/search?q=test&page=1&limit=10", nil)
	rec := httptest.NewRecorder()

	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", rec.Code)
	}

	if got := rec.Header().Get("X-Result-Count"); got != "0" {
		t.Fatalf("expected header X-Result-Count to be 0, got %s", got)
	}
	if got := rec.Header().Get("X-Page"); got != "1" {
		t.Fatalf("expected header X-Page to be 1, got %q", got)
	}
	if got := rec.Header().Get("X-Limit"); got != "10" {
		t.Fatalf("expected header X-Limit to be 10, got %q", got)
	}

	var body []DTO.SearchDto
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected body to be empty, got %s", rec.Body.String())
	}
	if len(body) != 0 {
		t.Fatalf("expected body to be empty, got %v", body)
	}

	if fs.lastQ != "test" {
		t.Fatalf("expected q=%q, got %q", "test", fs.lastQ)
	}
	if fs.lastLimit != 10 {
		t.Fatalf("expected limit=%d, got %d", 10, fs.lastLimit)
	}
	if fs.lastOffset != 0 {
		t.Fatalf("expected offset=%d, got %d", 0, fs.lastOffset)
	}

}

func TestSearchHandler_ReturnsResultsSortedAsGiven_WithHeaders(t *testing.T) {
	fs := &fakeSearchService{
		results: []DTO.SearchDto{
			{Title: "A", Url: "https://a", Score: 0.9},
			{Title: "B", Url: "https://b", Score: 0.4},
		},
	}
	h := NewSearchHandler(fs, 10, 50)

	req := httptest.NewRequest(http.MethodGet, "/search?q=test&page=2&limit=5", nil)
	rec := httptest.NewRecorder()

	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if got := rec.Header().Get("X-Result-Count"); got != "2" {
		t.Fatalf("expected X-Result-Count=2, got %q", got)
	}
	if got := rec.Header().Get("X-Page"); got != "2" {
		t.Fatalf("expected X-Page=2, got %q", got)
	}
	if got := rec.Header().Get("X-Limit"); got != "5" {
		t.Fatalf("expected X-Limit=5, got %q", got)
	}

	// offset = (page-1)*limit = (2-1)*5 = 5
	if fs.lastOffset != 5 || fs.lastLimit != 5 {
		t.Fatalf("expected limit=5 offset=5, got limit=%d offset=%d", fs.lastLimit, fs.lastOffset)
	}

	var body []DTO.SearchDto
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json body: %v", err)
	}
	if len(body) != 2 {
		t.Fatalf("expected 2 results, got %d", len(body))
	}
	if body[0].Title != "A" || body[1].Title != "B" {
		t.Fatalf("unexpected body order/content: %+v", body)
	}
}

func TestSearchHandler_ClampsLimitToMax(t *testing.T) {
	fs := &fakeSearchService{results: []DTO.SearchDto{}}
	h := NewSearchHandler(fs, 10, 50)

	req := httptest.NewRequest(http.MethodGet, "/search?q=test&page=1&limit=999", nil)
	rec := httptest.NewRecorder()

	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Limit"); got != "50" {
		t.Fatalf("expected X-Limit=50, got %q", got)
	}
	if fs.lastLimit != 50 {
		t.Fatalf("expected service limit=50, got %d", fs.lastLimit)
	}
}

func TestSearchHandler_MissingQueryReturns400(t *testing.T) {
	fs := &fakeSearchService{}
	h := NewSearchHandler(fs, 10, 50)

	req := httptest.NewRequest(http.MethodGet, "/search?page=1&limit=10", nil)
	rec := httptest.NewRecorder()

	h.Handle(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestSearchHandler_ServiceErrorReturns500(t *testing.T) {
	fs := &fakeSearchService{err: errors.New("db down")}
	h := NewSearchHandler(fs, 10, 50)

	req := httptest.NewRequest(http.MethodGet, "/search?q=test", nil)
	rec := httptest.NewRecorder()

	h.Handle(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}
