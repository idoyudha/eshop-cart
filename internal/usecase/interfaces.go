package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/entity"
)

type (
	CartMySQLRepo interface {
		Insert(context.Context, *entity.Cart) error
		GetByUserID(context.Context, uuid.UUID) ([]*entity.Cart, error)
		UpdateQtyAndNote(context.Context, *entity.Cart) (*uuid.UUID, error)
		UpdateNameAndPrice(context.Context, *entity.Cart) error
		DeleteMany(context.Context, uuid.UUIDs) error
		DeleteOne(context.Context, uuid.UUID) error
		UpdateProductQty(context.Context, *entity.Cart) error
	}

	CartRedisRepo interface {
		Save(context.Context, *entity.Cart) error
		GetUserCart(context.Context, string) ([]*entity.Cart, error)
		UpdateQtyAndNote(context.Context, *entity.Cart) error
		UpdateNameAndPrice(context.Context, *entity.Cart) error
		DeleteCart(context.Context, string, string) error
		DeleteCarts(context.Context, string, []string) error
		IsProductExistInUserCart(context.Context, string, string) (bool, error)
		UpdateProductQtyCart(context.Context, *entity.Cart) error
	}

	Cart interface {
		CreateCart(context.Context, *entity.Cart) (entity.Cart, error)
		GetUserCart(context.Context, uuid.UUID) ([]*entity.Cart, error)
		UpdateProductNameAndPriceCart(context.Context, *entity.Cart) error
		UpdateQtyAndNoteCart(context.Context, *entity.Cart) error
		DeleteCart(context.Context, uuid.UUID, uuid.UUID) error
		DeleteCarts(context.Context, uuid.UUID, uuid.UUIDs) error
		CheckOutCarts(context.Context, uuid.UUID, uuid.UUIDs, *entity.CheckoutAddress, string) error
	}
)
