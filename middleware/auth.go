package middleware

import (
	"backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := utils.VerifyAccessToken(token)
		if err != nil {
			ctx.JSON(401,gin.H{
				"success": false,
				"message": "Token invalid",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", payload)

		ctx.Set("user_id", int64(payload.Id))
		ctx.Set("role", payload.Role)
		ctx.Next()

	}
}