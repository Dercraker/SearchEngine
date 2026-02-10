package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dercraker/SearchEngine/internal/api/DTO"
)

type Searcher interface {
	Search(ctxKey any, r *http.Request, q string, limit, offset int) ([]DTO.SearchDto, error)
}

type SearchHandler struct {
	searcher               Searcher
	defaultLimit, maxLimit int
}

func NewSearchHandler(searcher Searcher, defaultLimit, maxLimit int) *SearchHandler {
	return &SearchHandler{searcher: searcher, defaultLimit: defaultLimit, maxLimit: maxLimit}
}

func (h *SearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "query param 'q' is required", http.StatusBadRequest)
		return
	}

	page := parseInt(r.URL.Query().Get("page"), 1)
	if page < 1 {
		page = 1
	}

	limit := parseInt(r.URL.Query().Get("limit"), h.defaultLimit)
	if limit < 1 {
		limit = h.defaultLimit
	}
	if limit > h.maxLimit {
		limit = h.maxLimit
	}

	offset := (page - 1) * limit

	results, err := h.searcher.Search(nil, r, q, limit, offset)
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Result-Count", strconv.Itoa(len(results)))
	w.WriteHeader(http.StatusOK)

	if results == nil {
		results = []DTO.SearchDto{}
	}
	_ = json.NewEncoder(w).Encode(results)

}

func parseInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}
