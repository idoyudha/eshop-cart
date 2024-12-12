package repo

import (
	"context"
	"fmt"

	"github.com/idoyudha/eshop-cart/internal/entity"
	rClient "github.com/idoyudha/eshop-cart/pkg/redis"
	"github.com/redis/go-redis/v9"
)

const (
	cartKey = "cart"
)

type CartRedisRepo struct {
	*rClient.RedisClient
}

func NewCartRedisRepo(client *rClient.RedisClient) *CartRedisRepo {
	return &CartRedisRepo{
		client,
	}
}

func getCartKey(cartID string) string {
	return fmt.Sprintf("cart:%s", cartID)
}

func getUserCartsKey(userID string) string {
	return fmt.Sprintf("user:%s:carts", userID)
}

// store cart data as hash -> cart:{cartID}
// add cartID to set -> user:{userID}:carts
func (r *CartRedisRepo) Save(ctx context.Context, cart *entity.Cart) error {
	cartKey := getCartKey(cart.ID.String())
	userCartsKey := getUserCartsKey(cart.UserID.String())

	// perform multiple operations
	pipe := r.Client.Pipeline()

	cartMap := map[string]interface{}{
		"id":               cart.ID.String(),
		"user_id":          cart.UserID.String(),
		"product_id":       cart.ProductID.String(),
		"product_name":     cart.ProductName,
		"product_price":    cart.ProductPrice,
		"product_quantity": cart.ProductQuantity,
		"note":             cart.Note,
	}

	pipe.HSet(ctx, cartKey, cartMap)
	pipe.SAdd(ctx, userCartsKey, cart.ID.String())

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save cart to redis: %w", err)
	}

	return nil
}

func (r *CartRedisRepo) GetUserCart(ctx context.Context, userID string) ([]*entity.Cart, error) {
	userCartKey := getUserCartsKey(userID)

	cartIDs, err := r.Client.SMembers(ctx, userCartKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart from redis: %w", err)
	}

	if len(cartIDs) == 0 {
		return nil, nil
	}

	pipe := r.Client.Pipeline()
	commands := make([]*redis.MapStringStringCmd, len(cartIDs))
	for i, cartID := range cartIDs {
		commands[i] = pipe.HGetAll(ctx, getCartKey(cartID))
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart from redis: %w", err)
	}

	carts := make([]*entity.Cart, 0, len(cartIDs))
	for _, cmd := range commands {
		cart := &entity.Cart{}
		err := cmd.Scan(cart)
		if err != nil {
			continue
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

func (r *CartRedisRepo) DeleteCarts(ctx context.Context, userID string, cartIDs []string) error {
	cartKeys := make([]string, len(cartIDs))
	for i, cartID := range cartIDs {
		cartKeys[i] = getCartKey(cartID)
	}

	pipe := r.Client.Pipeline()

	pipe.Del(ctx, cartKeys...)
	pipe.SRem(ctx, getUserCartsKey(userID), cartIDs)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete cart from redis: %w", err)
	}

	return nil
}
