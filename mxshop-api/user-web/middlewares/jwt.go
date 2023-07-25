package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/utils"
	"net/http"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		customClaims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("claims", customClaims)
		c.Set("userId", customClaims.Id)
		c.Next()
	}
}
