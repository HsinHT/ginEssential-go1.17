package middleware

import (
	"net/http"
	"strings"

	"example.com/ginessential/common"
	"example.com/ginessential/model"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 獲取 authorization header
		tokenString := ctx.GetHeader("Authorization")
		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "權限不足"})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "權限不足"})
			ctx.Abort()
			return
		}

		// 驗證通過後獲取 claim 中的 userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 用戶不存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "權限不足"})
			ctx.Abort()
			return
		}

		// 用戶存在 將 user 的訊息寫入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
