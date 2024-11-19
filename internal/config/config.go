package config

import (
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/models"
)

// ServiceConfig holds the configuration settings for the service
type ServiceConfig struct {
	Host              string
	Port              string
	IdleTimeout       time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	Router            *mux.Router
	DB                *models.InMemoryStore
}

// LoadServiceConfig loads and returns the service configuration with default values
func LoadServiceConfig() (*ServiceConfig, error) {
	return &ServiceConfig{
		Host:              os.Getenv("HOST"),
		Port:              os.Getenv("PORT"),
		IdleTimeout:       5 * time.Second,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Router:            mux.NewRouter(),           // Initialize new router
		DB:                models.NewInMemoryStore(), // Initialize in-memory store
	}, nil
}
