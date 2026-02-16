package http

import (
	"net/http"
	"time"

	"github.com/Dercraker/SearchEngine/internal/api/config"
)

func NewServer(cfg config.ApiConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}
}
