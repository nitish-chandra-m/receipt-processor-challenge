package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/config"
	"github.com/nitish-chandra-m/receipt-processor-challenge/internal/models"
)

type Handlers struct {
	ServiceConfig *config.ServiceConfig
}

type ProcessReceiptsRequest struct {
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Total        string        `json:"total"`
	Items        []ItemRequest `json:"items"`
}

type ItemRequest struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ProcessReceiptsResponse struct {
	Id string `json:"id"`
}

type GetPointsResponse struct {
	Points int `json:"points"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Handler function for the ProcessReceipts endpoint
func (h *Handlers) ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt ProcessReceiptsRequest
	err := json.NewDecoder(r.Body).Decode(&receipt)
	defer r.Body.Close()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	rec, err := createReceipt(h.ServiceConfig.DB, &receipt)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Spawn a goroutine to calculate points without blocking
	go rec.CalculatePoints()

	resp := ProcessReceiptsResponse{
		Id: rec.Id.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		jsonError(w, "failed to write response", http.StatusInternalServerError)
	}

}

// Handler function for the GetPoints endpoint
func (h *Handlers) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		jsonError(w, "invalid receipt id", http.StatusBadRequest)
		return
	}

	rec, err := models.FindReceiptById(h.ServiceConfig.DB, id)
	if err != nil {
		jsonError(w, "receipt not found", http.StatusNotFound)
		return
	}

	pts := GetPointsResponse{Points: rec.Points}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&pts); err != nil {
		jsonError(w, "failed to write response", http.StatusInternalServerError)
	}
}

// createReceipt creates a new receipt from the request body and stores it in the inmemory DB.
func createReceipt(db *models.InMemoryStore, receipt *ProcessReceiptsRequest) (*models.Receipt, error) {
	rec := &models.Receipt{}
	rec.Id = uuid.New()
	rec.Retailer = receipt.Retailer
	rec.Total, _ = strconv.ParseFloat(receipt.Total, 64)
	for _, item := range receipt.Items {
		price, _ := strconv.ParseFloat(item.Price, 64)
		it := models.Item{
			Id:               uuid.New(),
			ShortDescription: item.ShortDescription,
			Price:            price,
		}
		rec.Items = append(rec.Items, it)
	}
	rec.PurchaseDate, _ = time.Parse("2006-01-02", receipt.PurchaseDate)
	rec.PurchaseTime, _ = time.Parse("15:04", receipt.PurchaseTime)
	rec.PointsRubric = models.NewRubric(map[string]float64{})
	err := rec.Save(db)
	return rec, err
}

// jsonError sends a JSON-encoded error response to the client.
func jsonError(w http.ResponseWriter, errorString string, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := ErrorResponse{Error: errorString}
	if err := json.NewEncoder(w).Encode(&err); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}
