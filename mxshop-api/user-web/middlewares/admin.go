package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/utils"
	"net/http"
)

func AdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		c, b := context.Get("claims")
		if !b {
			context.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请登录",
			})
			context.Abort()
			return
		}
		claims, ok := c.(*utils.CustomClaims)
		if !ok {
			context.JSON(http.StatusUnauthorized, gin.H{
				"msg": "没有权限",
			})
			context.Abort()
			return
		}
		if claims.AuthorityId != 2 {
			context.JSON(http.StatusUnauthorized, gin.H{
				"msg": "没有权限",
			})
			context.Abort()
			return
		}
		context.Next()

	}
}
