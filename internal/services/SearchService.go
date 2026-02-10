package search

import (
	"log/slog"
	"net/http"

	"github.com/Dercraker/SearchEngine/internal/DAL"
	"github.com/Dercraker/SearchEngine/internal/api/DTO"
	"github.com/Dercraker/SearchEngine/internal/shared/requestId"
)

type searchService struct {
	logger *slog.Logger
	q      *DAL.Queries
}

func NewService(logger *slog.Logger, q *DAL.Queries) *searchService {
	return &searchService{logger: logger, q: q}
}

func (s *searchService) Search(_ any, r *http.Request, q string, limit, offset int) ([]DTO.SearchDto, error) {
	ctx := r.Context()

	rows, err := s.q.SearchDocuments(r.Context(), DAL.SearchDocumentsParams{
		Q:          q,
		PageLimit:  int32(limit),
		PageOffset: int32(offset),
	})

	if err != nil {
		attributes := []any{
			slog.String("operation", "SearchDocuments"),
			slog.String("q", q),
			slog.Int("limit", limit),
			slog.Int("offset", offset),
			slog.Any("error", err),
		}

		if reqId, isValid := requestId.Get(ctx); isValid {
			attributes = append(attributes, slog.String("request_id", reqId))
		}

		s.logger.Error("DB : Search Failed", attributes...)

		return nil, err
	}

	out := make([]DTO.SearchDto, 0, len(rows))
	for _, row := range rows {
		title := ""
		if row.Title.Valid {
			title = row.Title.String
		}

		out = append(out, DTO.SearchDto{
			Title: title,
			Url:   row.Url,
			Score: row.Score,
		})
	}

	return out, nil
}
