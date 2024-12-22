package entity

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	ProductID       uuid.UUID
	ProductName     string
	ProductPrice    float64
	ProductQuantity int64
	Note            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

func (c *Cart) GenerateCartID() error {
	cartID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	c.ID = cartID
	return nil
}
