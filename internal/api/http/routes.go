package http

import (
	"net/http"

	"github.com/Dercraker/SearchEngine/internal/api/http/handlers"
	"github.com/Dercraker/SearchEngine/internal/api/http/middleware"
	"github.com/Dercraker/SearchEngine/internal/shared/logging"
)

type RoutesDependencies struct {
	Health *handlers.HealthHandler
	Search *handlers.SearchHandler
}

func NewRouter(dependencies RoutesDependencies) http.Handler {
	mux := http.NewServeMux()

	RegisterRoutes(mux, dependencies)

	logger := logging.New()
	return middleware.Logging(logger)(mux)
}

func RegisterRoutes(mux *http.ServeMux, dependencies RoutesDependencies) {
	mux.HandleFunc("GET /health", dependencies.Health.Handle)
	mux.HandleFunc("GET /search", dependencies.Search.Handle)
}
