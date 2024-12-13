package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-cart/config"
)

const UserIDKey = "userID"

type authSuccessResponse struct {
	Code    int          `json:"code"`
	Data    authResponse `json:"data"`
	Message string       `json:"message"`
}
type authResponse struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
}

func cognitoMiddleware(auth config.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, newUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		authURL := fmt.Sprintf("%s/auth/%s", auth.BaseURL, tokenString)
		response, err := http.Get(authURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
			ctx.Abort()
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusUnauthorized, newUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
			ctx.Abort()
			return
		}

		var authSuccessResponse authSuccessResponse
		if err := json.Unmarshal(body, &authSuccessResponse); err != nil {
			ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
			ctx.Abort()
			return
		}

		ctx.Set(UserIDKey, authSuccessResponse.Data.UserID)
		ctx.Next()
	}
}
