package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpAdapter "github.com/harlesbayu/fleet-backend/internal/adapter/http_adapter"
	postgresAdapter "github.com/harlesbayu/fleet-backend/internal/adapter/postgres_adapter"
	"github.com/harlesbayu/fleet-backend/internal/config"
	"github.com/harlesbayu/fleet-backend/internal/container"
)

func main() {
	// Load config
	cfg := config.NewFromJSONFile("./configs", "config.json")

	// Init DB
	db, err := postgresAdapter.NewGormDB(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()

	// Init container
	c := container.NewContainer(db, &cfg)

	// Start server
	srv := httpAdapter.Start(c)

	// Graceful shutdown
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	httpAdapter.Shutdown(srv)
}
