package main

import (
	"context"
	"core-ticket/config"
	database_constants "core-ticket/constants/database"
	"core-ticket/database"
	"core-ticket/helpers/error"
	"core-ticket/helpers/shutdown"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	// Version will be overridden by ldflags at build time
	Version = "dev"
	// CommitSHA will be overridden by ldflags at build time
	CommitSHA = "none"
	// BuildTime will be overridden by ldflags at build time
	BuildTime = "unknown"

	// GlobalShutdownManager is the application-wide shutdown coordinator
	// Manages graceful shutdown for HTTP server and background workers
	GlobalShutdownManager *shutdown.Manager
)

func main() {
	cfg := config.AppConfig{}
	err := config.LoadConfig(&cfg)

	if err != nil {
		panic(fmt.Sprintf("error load config: %s", err.Error()))
	}
	errSetEnv := os.Setenv("TZ", cfg.Timezone)
	error.PanicIfError(errSetEnv)

	errSetEnv = os.Setenv("APP_JWT_SECRET", cfg.AppJwtSecret)
	error.PanicIfError(errSetEnv)

	errSetEnv = os.Setenv("APP_JWT_TOKEN_LIFE_SPAN", cfg.AppJwtTokenLifeSpan)
	error.PanicIfError(errSetEnv)

	errSetEnv = os.Setenv("APP_FRONTEND_URL", cfg.AppFrontendUrl)
	error.PanicIfError(errSetEnv)

	errSetEnv = os.Setenv("APP_REFRESH_SECRET", cfg.AppRefreshSecret)
	error.PanicIfError(errSetEnv)

	// Skip database initialization for health checks
	var db map[string]*sqlx.DB
	if os.Getenv("SKIP_DB_INIT") != "true" {
		db = database.NewDB(cfg)
		database_constants.InitDatabaseName(cfg)
	} else {
		fmt.Println("Skipping database initialization for health check")
		db = make(map[string]*sqlx.DB)
	}

	// Initialize Global Shutdown Manager
	GlobalShutdownManager = shutdown.NewManager()

	r := SetupRouter(db, GlobalShutdownManager)

	// HTTP Server Graceful Shutdown Implementation
	// Create HTTP server with explicit configuration
	srv := &http.Server{
		Addr:    cfg.AppHost + ":" + cfg.AppPort,
		Handler: r,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Server is starting on %s:%s", cfg.AppHost, cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started successfully. Press Ctrl+C to gracefully shutdown.")

	// Block until signal received
	<-quit
	log.Println("Shutdown signal received. Starting graceful shutdown...")

	// Signal shutdown to all background workers
	GlobalShutdownManager.Shutdown()

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown of HTTP server
	log.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("HTTP server shutdown completed successfully")
	}

	// Wait for background workers to complete
	log.Println("Waiting for background workers to complete...")
	if GlobalShutdownManager.Wait(30 * time.Second) {
		log.Println("All background workers completed successfully")
	} else {
		log.Println("Background workers did not complete within timeout")
	}

	log.Println("Application stopped")
}
