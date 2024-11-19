package models

import (
	"errors"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type InMemoryStore struct {
	Receipts map[uuid.UUID]*Receipt
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Receipts: make(map[uuid.UUID]*Receipt),
	}
}

type Receipt struct {
	Id           uuid.UUID `json:"id"`
	Retailer     string    `json:"retailer"`
	PurchaseDate time.Time `json:"purchaseDate"`
	PurchaseTime time.Time `json:"purchaseTime"`
	Total        float64   `json:"total"`
	Items        []Item    `json:"items"`
	PointsRubric Rubric    `json:"rubric"`
	Points       int       `json:"points"`
}

func (r *Receipt) Save(db *InMemoryStore) error {
	_, ok := db.Receipts[r.Id]
	if ok {
		return errors.New("receipt already exists")
	}
	db.Receipts[r.Id] = r
	return nil
}

func FindReceiptById(db *InMemoryStore, id uuid.UUID) (*Receipt, error) {
	r, ok := db.Receipts[id]
	if !ok {
		return nil, errors.New("receipt not found")
	}
	return r, nil
}

func (r *Receipt) CalculatePoints() {

	// Assign a point each for alphanumeric characters in retailer name
	retPts := 0.0
	for _, ch := range r.Retailer {
		if isAlphaNumeric(ch) {
			retPts += r.PointsRubric["retailer"]
		}
	}

	// Assign 50 points if total is a round dollar amount with no cents
	roundDollarPts := 0.0
	if isWholeNumber(r.Total) {
		roundDollarPts += r.PointsRubric["roundedTotal"]
	}

	// Assign 25 points if total is a multiple of 0.25
	totalMultiplePts := 0.0
	multipleOf := 0.25
	div := r.Total / multipleOf
	if isWholeNumber(div) {
		totalMultiplePts += r.PointsRubric["totalMultiple"]
	}

	// Assign 5 points for every pair of items in the receipt
	pairPts := 0.0
	pairs := int(len(r.Items) / 2)
	for range pairs {
		pairPts += r.PointsRubric["pairOfItems"]
	}

	// Assign ceil(price*0.2) points if the item short description is a multiple of 3
	descPts := 0.0
	for _, item := range r.Items {
		// fmt.Printf("%s, %d\n", item.ShortDescription, len(item.ShortDescription))
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			descPts += math.Ceil(item.Price * r.PointsRubric["descriptionLength"])
		}
	}

	// Assign 6 points if day of purchase is odd
	oddPurchaseDayPts := 0.0
	if r.PurchaseDate.Day()%2 != 0 {
		oddPurchaseDayPts += r.PointsRubric["oddPurchaseDay"]
	}

	// Assign 10 points if time of purchase is after 2pm and before 4pm
	timeOfPurchasePts := 0.0
	rangeStart, _ := time.Parse(time.Kitchen, "02:00PM")
	rangeEnd, _ := time.Parse(time.Kitchen, "04:00PM")
	if r.PurchaseTime.After(rangeStart) && r.PurchaseTime.Before(rangeEnd) {
		timeOfPurchasePts += r.PointsRubric["afternoonPurchase"]
	}

	r.Points = int(retPts + roundDollarPts + totalMultiplePts + pairPts + descPts + oddPurchaseDayPts + timeOfPurchasePts)
}

func isAlphaNumeric(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch)
}

func isWholeNumber(num float64) bool {
	_, fracPart := math.Modf(num)
	return fracPart == 0
}
