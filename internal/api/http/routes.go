package http

import (
	"net/http"

	"github.com/Dercraker/SearchEngine/internal/api/http/handlers"
	"github.com/Dercraker/SearchEngine/internal/api/http/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	RegisterRoutes(mux)

	return middleware.Logging(mux)
}

func RegisterRoutes(mux *http.ServeMux) {
	health := handlers.NewHealthHandler()

	mux.HandleFunc("GET /health", health.Handle)

}
