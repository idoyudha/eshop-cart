package v1

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-cart/config"
	"github.com/idoyudha/eshop-cart/internal/usecase"
	"github.com/idoyudha/eshop-cart/pkg/logger"
)

func NewRouter(
	handler *gin.Engine,
	ucc usecase.Cart,
	l logger.Interface,
	auth config.AuthService,
) {
	// handler.Use(gin.Logger())
	// handler.Use(gin.Recovery())
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	authMid := cognitoMiddleware(auth)

	h := handler.Group("/v1")
	{
		newCartRoutes(h, ucc, l, authMid)
	}
}
