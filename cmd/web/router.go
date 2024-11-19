package main

import (
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/config"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/handlers"
)

func SetupRouter(cfg *config.ServiceConfig) {
	r := cfg.Router

	handlers := &handlers.Handlers{ServiceConfig: cfg}

	r.HandleFunc("/receipts/process", handlers.ProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")
}
