package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/entity"
)

type CartUseCase struct {
	repoRedis CartRedisRepo
	repoMySQL CartMySQLRepo
}

func NewCartUseCase(repoRedis CartRedisRepo, repoMySQL CartMySQLRepo) *CartUseCase {
	return &CartUseCase{
		repoRedis,
		repoMySQL,
	}
}

func (u *CartUseCase) CreateCart(ctx context.Context, cart *entity.Cart) error {
	if errInsert := u.repoMySQL.Insert(ctx, cart); errInsert != nil {
		return errInsert
	}

	errSave := u.repoRedis.Save(ctx, cart)

	// delete cart from mysql if save to redis failed
	if errSave != nil {
		_ = u.repoMySQL.DeleteOne(ctx, cart.ID)
		return errSave
	}

	return nil
}

func (u *CartUseCase) GetUserCart(ctx context.Context, userID uuid.UUID) ([]*entity.Cart, error) {
	// get cart from redis
	carts, errGet := u.repoRedis.GetUserCart(ctx, userID.String())
	if errGet != nil {
		return nil, errGet
	}

	// if cart found, return it
	if len(carts) > 0 {
		return carts, nil
	}

	// if cart not found, get cart from mysql
	carts, errGet = u.repoMySQL.GetByUserID(ctx, userID)
	if errGet != nil {
		return nil, errGet
	}

	// save cart to redis
	for _, cart := range carts {
		if errSave := u.repoRedis.Save(ctx, cart); errSave != nil {
			return nil, errSave
		}
	}

	return carts, nil
}

func (u *CartUseCase) UpdateQtyAndNoteCart(ctx context.Context, cart *entity.Cart) error {
	if errUpdate := u.repoMySQL.UpdateQtyAndNote(ctx, cart); errUpdate != nil {
		return errUpdate
	}

	if errSave := u.repoRedis.UpdateQtyAndNote(ctx, cart); errSave != nil {
		return errSave
	}

	return nil
}

func (u *CartUseCase) DeleteCarts(ctx context.Context, userID uuid.UUID, cartIDs uuid.UUIDs) error {
	if errDelete := u.repoMySQL.DeleteMany(ctx, cartIDs); errDelete != nil {
		return errDelete
	}

	var redisCartIDs []string
	for _, cartID := range cartIDs {
		redisCartIDs = append(redisCartIDs, cartID.String())
	}
	if errDelete := u.repoRedis.DeleteCarts(ctx, userID.String(), redisCartIDs); errDelete != nil {
		return errDelete
	}

	return nil
}
