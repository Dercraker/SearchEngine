package app

import (
	"net/http"

	"github.com/Dercraker/SearchEngine/internal/api/config"
	httpx "github.com/Dercraker/SearchEngine/internal/api/http"
)

type App struct {
	cfg    config.Config
	router http.Handler
}

func New(cfg config.Config) *App {
	//Ici on instancies les services métier  + repos d'infra
	// ex : searchService := search.NewService(indexRepo, ....)
	// handlers := handlers.NewHandlers(searchService, ....)

	r := httpx.NewRouter() //Transport Layer

	return &App{cfg: cfg, router: r}
}

func (a *App) Router() http.Handler {
	return a.router
}
