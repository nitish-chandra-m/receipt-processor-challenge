package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/config"
)

func main() {
	// Load .env file if the app environment is not set
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Load web service configuration
	cfg, err := config.LoadServiceConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up the HTTP server with appropriate timeouts and handler
	srv := http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Bind to the configured address and port
		IdleTimeout:       cfg.IdleTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		Handler:           cfg.Router, // Set the router to handle incoming requests
	}

	// Set up routes and handlers
	SetupRouter(cfg)

	// Start the server in a new goroutine so it doesn't block
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Gracefully shut down the server when interrupt signal is received
	gracefulShutdown(&srv)
}

func gracefulShutdown(srv *http.Server) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Wait for the interrupt signal (Ctrl+C or kill signal)
	<-stopChan
	log.Println("Shutting down server")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully shutdown")
	os.Exit(0)
}
