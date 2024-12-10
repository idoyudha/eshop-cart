package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductPrice    float64   `json:"product_price"`
	ProductQuantity int64     `json:"product_quantity"`
	Note            string    `json:"note"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at"`
}
