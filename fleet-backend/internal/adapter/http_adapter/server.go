package http_adapter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harlesbayu/fleet-backend/internal/container"
)

// Start runs the HTTP server asynchronously
func Start(c *container.Container) *http.Server {
	engine, addr := initServer(c)

	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// Run the server in a goroutine (non-blocking)
	go func() {
		log.Printf("Server running on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped: %v", err)
		}
	}()

	return srv
}

// initServer initializes the router, middleware, and routes
func initServer(c *container.Container) (*gin.Engine, string) {
	r := gin.Default()

	// Register global middleware (if any in the future)
	registerMiddleware(r, c)

	// Register all routes
	registerRoutes(r, c)

	port := fmt.Sprintf(":%d", c.Config.Server.Port)
	return r, port
}

// Shutdown performs a graceful shutdown
func Shutdown(srv *http.Server) {
	log.Println("Stopping HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("HTTP server stopped gracefully.")
}

func registerMiddleware(r *gin.Engine, c *container.Container) {}

func registerRoutes(r *gin.Engine, c *container.Container) {
	r.GET("/ping", c.PingHandler.Ping)
	r.GET("/vehicles/:vehicle_id/location", c.VehicleHandler.GetLastLocation)
	r.GET("/vehicles/:vehicle_id/history", c.VehicleHandler.GetHistory)
}
