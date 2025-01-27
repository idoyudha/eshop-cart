package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/config"
	"github.com/idoyudha/eshop-cart/internal/entity"
	"github.com/idoyudha/eshop-cart/internal/utils"
)

type CartUseCase struct {
	repoRedis    CartRedisRepo
	repoMySQL    CartMySQLRepo
	orderService config.OrderService
}

func NewCartUseCase(
	repoRedis CartRedisRepo,
	repoMySQL CartMySQLRepo,
	orderService config.OrderService,
) *CartUseCase {
	return &CartUseCase{
		repoRedis,
		repoMySQL,
		orderService,
	}
}

func (u *CartUseCase) CreateCart(ctx context.Context, cart *entity.Cart) (entity.Cart, error) {
	err := cart.GenerateCartID()
	if err != nil {
		return entity.Cart{}, err
	}

	exist, errExist := u.repoRedis.IsProductExistInUserCart(ctx, cart.UserID.String(), cart.ProductID.String())
	if errExist != nil {
		return entity.Cart{}, errExist
	}

	var errInsertOrUpdateRedis error

	if !exist {
		if errInsert := u.repoMySQL.Insert(ctx, cart); errInsert != nil {
			return entity.Cart{}, errInsert
		}

		errInsertOrUpdateRedis = u.repoRedis.Save(ctx, cart)
	} else {
		if errUpdateProductQty := u.repoMySQL.UpdateProductQty(ctx, cart); errUpdateProductQty != nil {
			return entity.Cart{}, errUpdateProductQty
		}

		errInsertOrUpdateRedis = u.repoRedis.UpdateProductQtyCart(ctx, cart)
	}

	// delete cart from mysql if insert or update to redis failed
	if errInsertOrUpdateRedis != nil {
		_, _ = u.repoMySQL.DeleteOne(ctx, cart.ID)
		return entity.Cart{}, errInsertOrUpdateRedis
	}

	return *cart, nil
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
	productID, errUpdate := u.repoMySQL.UpdateQtyAndNote(ctx, cart)
	if errUpdate != nil {
		return errUpdate
	}

	cart.ProductID = *productID
	if errSave := u.repoRedis.UpdateQtyAndNote(ctx, cart); errSave != nil {
		return errSave
	}

	return nil
}

func (u *CartUseCase) UpdateProductNameAndPriceCart(ctx context.Context, cart *entity.Cart) error {
	if errUpdate := u.repoMySQL.UpdateNameAndPrice(ctx, cart); errUpdate != nil {
		return errUpdate
	}

	if errSave := u.repoRedis.UpdateNameAndPrice(ctx, cart); errSave != nil {
		return errSave
	}

	return nil
}

func (u *CartUseCase) DeleteCart(ctx context.Context, userID uuid.UUID, cartID uuid.UUID) error {
	productID, errDelete := u.repoMySQL.DeleteOne(ctx, cartID)
	if errDelete != nil {
		return errDelete
	}

	if errDelete := u.repoRedis.DeleteCart(ctx, userID.String(), productID.String()); errDelete != nil {
		return errDelete
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
	if errDelete := u.repoRedis.DeleteCarts(ctx, userID.String()); errDelete != nil {
		return errDelete
	}

	return nil
}

type createOrderRequest struct {
	Items   []createItemsOrderRequest `json:"items"`
	Address createAddressOrderRequest `json:"address"`
}

type createItemsOrderRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Price     float64   `json:"price"`
}

type createAddressOrderRequest struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipcode"`
	Note    string `json:"note"`
}

type restSuccessCreateOrder struct {
	Code    int           `json:"code"`
	Data    orderResponse `json:"data"`
	Message string        `json:"message"`
}

type orderResponse struct {
	Status     string               `json:"status"`
	TotalPrice float64              `json:"total_price"`
	Items      []itemsOrderResponse `json:"items"`
	Address    addressOrderResponse `json:"address"`
}

type itemsOrderResponse struct {
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"price"`
	Quantity  int64     `json:"quantity"`
	Note      string    `json:"note"`
}

type addressOrderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Street  string    `json:"street"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	ZipCode string    `json:"zipcode"`
	Note    string    `json:"note"`
}

func cartToCreateOrderRequest(carts []*entity.Cart, cartIDs uuid.UUIDs, address *entity.CheckoutAddress) createOrderRequest {
	var items []createItemsOrderRequest
	for _, cart := range carts {
		if utils.IDInSliceUUID(cart.ID, cartIDs) {
			items = append(items, createItemsOrderRequest{
				ProductID: cart.ProductID,
				Quantity:  cart.ProductQuantity,
				Price:     cart.ProductPrice,
			})
		}
	}
	return createOrderRequest{
		Items: items,
		Address: createAddressOrderRequest{
			Street:  address.Street,
			City:    address.City,
			State:   address.State,
			ZipCode: address.ZipCode,
		},
	}
}

func (u *CartUseCase) CheckOutCarts(ctx context.Context, userID uuid.UUID, cartIDs uuid.UUIDs, address *entity.CheckoutAddress, token string) error {
	// 1. get cart from redis
	carts, err := u.GetUserCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// 2. request create order
	createOrderURL := fmt.Sprintf("%s/v1/orders", u.orderService.BaseURL)
	createOrderReq := cartToCreateOrderRequest(carts, cartIDs, address)

	requestBody, err := json.Marshal(createOrderReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createOrderURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var successCreateOrder restSuccessCreateOrder
	if err := json.Unmarshal(body, &successCreateOrder); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// 3. delete cart from mysql and redis
	if errDelete := u.DeleteCarts(ctx, userID, cartIDs); errDelete != nil {
		return fmt.Errorf("failed to delete cart: %w", errDelete)
	}

	return nil
}
