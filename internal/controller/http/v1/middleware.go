package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const userIDKey = "userID"

func cognitoMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, newUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		// hit to auth service
	}
}
