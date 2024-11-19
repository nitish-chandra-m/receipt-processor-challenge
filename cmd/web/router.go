package main

import (
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/config"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/handlers"
)

// SetupRouter sets up the API endpoints and their corresponding handler functions
func SetupRouter(cfg *config.ServiceConfig) {
	r := cfg.Router

	// Pass service config which contains the DB to a new Handler (dependency injection)
	handlers := &handlers.Handlers{ServiceConfig: cfg}

	// Define the two required endpoints and their respective handlers and HTTP Methods
	r.HandleFunc("/receipts/process", handlers.ProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")
}
