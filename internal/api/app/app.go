package app

import (
	"net/http"

	"github.com/Dercraker/SearchEngine/internal/api/infra/dbx"
	"github.com/Dercraker/SearchEngine/internal/shared/logging"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Dercraker/SearchEngine/internal/DAL"
	"github.com/Dercraker/SearchEngine/internal/api/config"
	httpx "github.com/Dercraker/SearchEngine/internal/api/http"
	"github.com/Dercraker/SearchEngine/internal/api/http/handlers"
	"github.com/Dercraker/SearchEngine/internal/services"
)

type App struct {
	cfg    config.Config
	router http.Handler
}

func New(cfg config.Config) *App {
	logger := logging.New()

	dbConn, err := dbx.Open(logger, dbx.Options{
		DSN:             cfg.DatabaseDSN,
		PingTimeout:     cfg.DBPingTimeout,
		FailFast:        cfg.DBFailFast,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
		ConnMaxIdleTime: cfg.DBConnMaxIdleTime,
	})

	if err != nil {
		panic(err)
	}

	queries := DAL.New(dbConn)

	searchService := search.NewService(logger, queries)

	healthHandler := handlers.NewHealthHandler()
	searchHandler := handlers.NewSearchHandler(searchService, cfg.SearchLimitDefault, cfg.SearchLimitMax)

	r := httpx.NewRouter(httpx.RoutesDependencies{
		Health: healthHandler,
		Search: searchHandler,
	})

	return &App{cfg: cfg, router: r}
}

func (a *App) Router() http.Handler {
	return a.router
}
