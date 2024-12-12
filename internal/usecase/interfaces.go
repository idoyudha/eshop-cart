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
		Update(context.Context, *entity.Cart) error
		DeleteMany(context.Context, uuid.UUIDs) error
	}

	CartRedisRepo interface {
		Save(context.Context, *entity.Cart) error
		GetUserCart(context.Context, string) ([]*entity.Cart, error)
		DeleteCarts(context.Context, string, []string) error
	}

	Cart interface {
		CreateCart(context.Context, *entity.Cart) error
		GetUserCart(context.Context, uuid.UUID) ([]*entity.Cart, error)
		UpdateCart(context.Context, *entity.Cart) error
		DeleteCarts(context.Context, uuid.UUID, uuid.UUIDs) error
	}
)
