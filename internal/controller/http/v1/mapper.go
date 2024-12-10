package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/entity"
)

func CreateCartRequestToCartEntity(userID uuid.UUID, req CreateCartRequest) (entity.Cart, error) {
	cartID, err := uuid.NewV7()
	if err != nil {
		return entity.Cart{}, err
	}
	return entity.Cart{
		ID:              cartID,
		UserID:          userID,
		ProductID:       req.ProductID,
		ProductName:     req.ProductName,
		ProductPrice:    req.ProductPrice,
		ProductQuantity: req.ProductQuantity,
		Note:            req.Note,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}
