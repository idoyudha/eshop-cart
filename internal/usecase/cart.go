package usecase

import (
	"context"

	"github.com/idoyudha/eshop-cart/internal/entity"
)

type (
	CartMySQLRepo interface {
		Insert(context.Context, *entity.Cart) error
		Update(context.Context, *entity.Cart) error
		Get(context.Context, int64) (*entity.Cart, error)
		Delete(context.Context, int64) error
	}

	CartRedisRepo interface {
		Save(context.Context, *entity.Cart) error
		Get(context.Context, int64) (*entity.Cart, error)
		Delete(context.Context, int64) error
	}

	Cart interface {
		CreateCart(context.Context, *entity.Cart) error
		GetCart(context.Context, int64) (*entity.Cart, error)
		UpdateCart(context.Context, *entity.Cart) error
		DeleteCart(context.Context, int64) error
	}
)
