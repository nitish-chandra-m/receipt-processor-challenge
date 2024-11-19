package models

import (
	"github.com/google/uuid"
)

type Item struct {
	Id               uuid.UUID `json:"id"`
	ShortDescription string    `json:"shortDescription"`
	Price            float64   `json:"price"`
}
