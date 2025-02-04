package repo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
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
	cartKey := getCartKey(cart.ProductID.String())
	userCartsKey := getUserCartsKey(cart.UserID.String())

	// perform multiple operations
	pipe := r.Client.Pipeline()

	cartMap := map[string]interface{}{
		"id":                cart.ID.String(),
		"user_id":           cart.UserID.String(),
		"product_id":        cart.ProductID.String(),
		"product_name":      cart.ProductName,
		"product_image_url": cart.ProductImageURL,
		"product_price":     cart.ProductPrice,
		"product_quantity":  cart.ProductQuantity,
		"note":              cart.Note,
	}

	pipe.HSet(ctx, cartKey, cartMap)
	pipe.SAdd(ctx, userCartsKey, cart.ProductID.String())

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
		cartData := cmd.Val()
		if len(cartData) == 0 {
			continue
		}

		cartID, _ := uuid.Parse(cartData["id"])
		userID, _ := uuid.Parse(cartData["user_id"])
		productID, _ := uuid.Parse(cartData["product_id"])
		productQuantity, _ := strconv.ParseInt(cartData["product_quantity"], 10, 64)
		productPrice, _ := strconv.ParseFloat(cartData["product_price"], 64)

		cart := &entity.Cart{
			ID:              cartID,
			UserID:          userID,
			ProductID:       productID,
			ProductName:     cartData["product_name"],
			ProductImageURL: cartData["product_image_url"],
			ProductPrice:    productPrice,
			ProductQuantity: productQuantity,
			Note:            cartData["note"],
		}

		carts = append(carts, cart)
	}

	return carts, nil
}

func (r *CartRedisRepo) UpdateQtyAndNote(ctx context.Context, cart *entity.Cart) error {
	cartKey := getCartKey(cart.ProductID.String())
	exists, err := r.Client.Exists(ctx, cartKey).Result()
	if err != nil {
		return fmt.Errorf("cart is not exist: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("cart not found for updating quantity and note")
	}

	pipe := r.Client.Pipeline()
	pipe.HSet(ctx, cartKey, map[string]interface{}{
		"product_quantity": cart.ProductQuantity,
		"note":             cart.Note,
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}

// TODO: REFACTOR -> bcs already refactor the product_id as a key in redis
func (r *CartRedisRepo) UpdateNameAndPrice(ctx context.Context, cart *entity.Cart) error {
	// get all cart keys from Redis that contain this product
	pattern := fmt.Sprintf("cart:*")
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	pipe := r.Client.Pipeline()
	cartUpdated := false

	for iter.Next(ctx) {
		cartKey := iter.Val()

		// get the product ID from this cart
		cartProductID, err := r.Client.HGet(ctx, cartKey, "product_id").Result()
		if err != nil {
			continue
		}

		// if this cart contains our target product, update it
		if cartProductID == cart.ProductID.String() {
			pipe.HSet(ctx, cartKey, map[string]interface{}{
				"product_name":  cart.ProductName,
				"product_price": cart.ProductPrice,
			})
			cartUpdated = true
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning carts: %w", err)
	}

	// execute pipeline if we found carts to update
	if cartUpdated {
		_, err := pipe.Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update carts: %w", err)
		}
	}

	return nil
}

func (r *CartRedisRepo) DeleteCart(ctx context.Context, userID string, cartID string) error {
	pipe := r.Client.Pipeline()
	pipe.Del(ctx, getCartKey(cartID))
	pipe.SRem(ctx, getUserCartsKey(userID), cartID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete cart from redis: %w", err)
	}

	return nil
}

func (r *CartRedisRepo) DeleteCarts(ctx context.Context, userID string) error {
	cartIDs, err := r.Client.SMembers(ctx, getUserCartsKey(userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get cartid members from redis: %w", err)
	}

	cartKeys := make([]string, len(cartIDs))
	for i, cartID := range cartIDs {
		cartKeys[i] = getCartKey(cartID)
	}

	pipe := r.Client.Pipeline()

	pipe.Del(ctx, cartKeys...)
	pipe.SRem(ctx, getUserCartsKey(userID), cartIDs)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete cart from redis: %w", err)
	}

	return nil
}

func (r *CartRedisRepo) IsProductExistInUserCart(ctx context.Context, userID string, productID string) (bool, error) {
	cartKey := getUserCartsKey(userID)

	exists, err := r.Client.SIsMember(ctx, cartKey, productID).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if cart exists: %w", err)
	}

	return exists, nil
}

func (r *CartRedisRepo) UpdateProductQtyCart(ctx context.Context, cart *entity.Cart) error {
	cartKey := getCartKey(cart.ProductID.String())

	exists, err := r.Client.Exists(ctx, cartKey).Result()
	if err != nil {
		return fmt.Errorf("cart is not exist: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("cart not found for updating quantity")
	}

	pipe := r.Client.Pipeline()
	pipe.HIncrBy(ctx, cartKey, "product_quantity", cart.ProductQuantity)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}
