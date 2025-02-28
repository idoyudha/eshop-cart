package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/usecase"
	"github.com/idoyudha/eshop-cart/pkg/logger"
)

type cartRoutes struct {
	uc usecase.Cart
	l  logger.Interface
}

func newCartRoutes(handler *gin.RouterGroup, uc usecase.Cart, l logger.Interface, authMid gin.HandlerFunc) {
	r := &cartRoutes{uc: uc, l: l}

	h := handler.Group("/carts").Use(authMid)
	{
		h.POST("", r.createCart)
		h.GET("/user", r.getCartByUserID)
		h.PATCH("/:id", r.updateCart)
		h.DELETE("/:id", r.deleteCart)
		h.PATCH("/deletes", r.deleteCarts)
		h.POST("/checkout", r.checkOutCarts)
	}
}

type createCartRequest struct {
	ProductID       uuid.UUID `json:"product_id" binding:"required"`
	ProductName     string    `json:"product_name" binding:"required"`
	ProductImageURL string    `json:"product_image_url" inding:"required"`
	ProductPrice    float64   `json:"product_price" binding:"required"`
	ProductQuantity int64     `json:"product_quantity" binding:"required"`
	Note            string    `json:"note"`
}

type createCartResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductImageURL string    `json:"product_image_url"`
	ProductPrice    float64   `json:"product_price"`
	ProductQuantity int64     `json:"product_quantity"`
	Note            string    `json:"note"`
}

func (r *cartRoutes) createCart(ctx *gin.Context) {
	var req createCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	cartEntity := createCartRequestToCartEntity(userID.(uuid.UUID), req)

	cart, err := r.uc.CreateCart(ctx.Request.Context(), &cartEntity)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	cartResponse := cartEntityToCreateCartResponse(cart)

	ctx.JSON(http.StatusCreated, newCreateSuccess(cartResponse))
}

type getCartResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductImageURL string    `json:"product_image_url"`
	ProductPrice    float64   `json:"product_price"`
	ProductQuantity int64     `json:"product_quantity"`
	Note            string    `json:"note"`
}

// get cart by user id
func (r *cartRoutes) getCartByUserID(ctx *gin.Context) {
	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	carts, err := r.uc.GetUserCart(ctx.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - getCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	cartsResponse := cartEntitiesToGetCartResponse(carts)

	ctx.JSON(http.StatusOK, newGetSuccess(cartsResponse))
}

type updateCartRequest struct {
	ProductQuantity int64  `json:"product_quantity" binding:"required"`
	Note            string `json:"note"`
}

type updateCartResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductImageURL string    `json:"product_image_url"`
	ProductPrice    float64   `json:"product_price"`
	ProductQuantity int64     `json:"product_quantity"`
	Note            string    `json:"note"`
}

func (r *cartRoutes) updateCart(ctx *gin.Context) {
	cartID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - updateCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	var req updateCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	cart := updateCartRequestToCartEntity(cartID, userID.(uuid.UUID), req)

	err = r.uc.UpdateQtyAndNoteCart(ctx.Request.Context(), &cart)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - updateCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	cartResponse := cartEntityToUpdateCartResponse(cart)

	ctx.JSON(http.StatusOK, newUpdateSuccess(cartResponse))
}

func (r *cartRoutes) deleteCart(ctx *gin.Context) {
	cartID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - deleteCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	err = r.uc.DeleteCart(ctx.Request.Context(), userID.(uuid.UUID), cartID)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - deleteCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newDeleteSuccess())
}

type deleteCartsRequest struct {
	CartIDs uuid.UUIDs `json:"cart_ids"`
}

func (r *cartRoutes) deleteCarts(ctx *gin.Context) {
	var req deleteCartsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - deleteCarts")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	err := r.uc.DeleteCarts(ctx.Request.Context(), userID.(uuid.UUID), req.CartIDs)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - deleteCarts")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newDeleteSuccess())
}

type checkoutCartsRequest struct {
	CartIDs uuid.UUIDs             `json:"cart_ids" binding:"required"`
	Address checkoutAddressRequest `json:"address" binding:"required"`
}

type checkoutAddressRequest struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	ZipCode string `json:"zipcode" binding:"required"`
}

func (r *cartRoutes) checkOutCarts(ctx *gin.Context) {
	var req checkoutCartsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - checkOutCarts")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	token, exist := ctx.Get(TokenKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - orderRoutes - createOrder")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("token not exist"))
		return
	}

	address := checkoutAddressRequestToCheckoutAddressEntity(req.Address)

	err := r.uc.CheckOutCarts(
		ctx.Request.Context(),
		userID.(uuid.UUID),
		req.CartIDs,
		&address,
		token.(string),
	)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - checkOutCarts")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newCheckoutSuccess())
}
