package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/entity"
)

func createCartRequestToCartEntity(userID uuid.UUID, req createCartRequest) entity.Cart {
	return entity.Cart{
		UserID:          userID,
		ProductID:       req.ProductID,
		ProductName:     req.ProductName,
		ProductPrice:    req.ProductPrice,
		ProductQuantity: req.ProductQuantity,
		Note:            req.Note,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func cartEntityToCreateCartResponse(cart entity.Cart) createCartResponse {
	return createCartResponse{
		ID:              cart.ID,
		UserID:          cart.UserID,
		ProductID:       cart.ProductID,
		ProductName:     cart.ProductName,
		ProductPrice:    cart.ProductPrice,
		ProductQuantity: cart.ProductQuantity,
		Note:            cart.Note,
	}
}

func cartEntitiesToGetCartResponse(cart []*entity.Cart) []getCartResponse {
	var res []getCartResponse
	for _, c := range cart {
		res = append(res, getCartResponse{
			ID:              c.ID,
			UserID:          c.UserID,
			ProductID:       c.ProductID,
			ProductName:     c.ProductName,
			ProductPrice:    c.ProductPrice,
			ProductQuantity: c.ProductQuantity,
			Note:            c.Note,
		})
	}
	return res
}

func updateCartRequestToCartEntity(cartID uuid.UUID, userID uuid.UUID, req updateCartRequest) entity.Cart {
	return entity.Cart{
		ID:              cartID,
		UserID:          userID,
		ProductQuantity: req.ProductQuantity,
		Note:            req.Note,
		UpdatedAt:       time.Now(),
	}
}

func cartEntityToUpdateCartResponse(cart entity.Cart) updateCartResponse {
	return updateCartResponse{
		ID:              cart.ID,
		UserID:          cart.UserID,
		ProductID:       cart.ProductID,
		ProductName:     cart.ProductName,
		ProductPrice:    cart.ProductPrice,
		ProductQuantity: cart.ProductQuantity,
		Note:            cart.Note,
	}
}

func checkoutAddressRequestToCheckoutAddressEntity(req checkoutAddressRequest) entity.CheckoutAddress {
	return entity.CheckoutAddress{
		Street:  req.Street,
		City:    req.City,
		State:   req.State,
		ZipCode: req.ZipCode,
	}
}
