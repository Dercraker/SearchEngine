package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dercraker/SearchEngine/internal/api/app"
	"github.com/Dercraker/SearchEngine/internal/api/config"
	httpx "github.com/Dercraker/SearchEngine/internal/api/http"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("config: %v", err)
	}

	application := app.New(cfg)

	srv := httpx.NewServer(cfg, application.Router())

	go func() {
		log.Printf("Search-API listening on %s\n", cfg.Addr)

		if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("Http serve: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down...")
	_ = srv.Shutdown(ctx)
}
