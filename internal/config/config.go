package config

import (
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/models"
)

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

func NewServiceConfig() (*ServiceConfig, error) {
	return &ServiceConfig{
		Host:              os.Getenv("HOST"),
		Port:              os.Getenv("PORT"),
		IdleTimeout:       5 * time.Second,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Router:            mux.NewRouter(),
		DB:                models.NewInMemoryStore(),
	}, nil
}
