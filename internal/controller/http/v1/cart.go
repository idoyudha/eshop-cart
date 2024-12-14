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
		h.PUT("/:id", r.updateCart)
		h.PATCH("/deletes", r.deleteCarts)
	}
}

type CreateCartRequest struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductPrice    float64   `json:"product_price"`
	ProductQuantity int64     `json:"product_quantity"`
	Note            string    `json:"note"`
}

func (r *cartRoutes) createCart(ctx *gin.Context) {
	var req CreateCartRequest
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

	cart, err := CreateCartRequestToCartEntity(userID.(uuid.UUID), req)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	err = r.uc.CreateCart(ctx.Request.Context(), &cart)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - createCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, newCreateSuccess(cart))
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

	ctx.JSON(http.StatusOK, newGetSuccess(carts))
}

type UpdateCartRequest struct {
	ProductQuantity int64  `json:"product_quantity"`
	Note            string `json:"note"`
}

func (r *cartRoutes) updateCart(ctx *gin.Context) {
	cartID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - updateCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	var req UpdateCartRequest
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

	cart := UpdateCartRequestToCartEntity(cartID, userID.(uuid.UUID), req)

	err = r.uc.UpdateCart(ctx.Request.Context(), &cart)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - updateCart")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newUpdateSuccess(cart))
}

type DeleteCartsRequest struct {
	CartIDs uuid.UUIDs `json:"cart_ids"`
}

func (r *cartRoutes) deleteCarts(ctx *gin.Context) {
	var req DeleteCartsRequest
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
