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

func newCartRoutes(handler *gin.RouterGroup, uc usecase.Cart, l logger.Interface) {
	r := &cartRoutes{uc: uc, l: l}

	h := handler.Group("/carts")
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

	// TODO: get from token
	userID, _ := uuid.NewV7()

	cart, err := CreateCartRequestToCartEntity(userID, req)
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
	// TODO: get from token
	userID, _ := uuid.NewV7()

	carts, err := r.uc.GetUserCart(ctx.Request.Context(), userID)
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

	// TODO: get from token
	userID, _ := uuid.NewV7()
	cart := UpdateCartRequestToCartEntity(cartID, userID, req)

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

	// TODO: get from token
	userID, _ := uuid.NewV7()

	err := r.uc.DeleteCarts(ctx.Request.Context(), userID, req.CartIDs)
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - deleteCarts")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newDeleteSuccess())
}
