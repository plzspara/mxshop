package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/api"
	"mxshop-api/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	{
		userRouter.GET("list", middlewares.JwtAuth(), middlewares.AdminAuth(), api.GetUserList)
		userRouter.POST("login", api.PasswordLogin)
		userRouter.GET("captcha", api.GetCaptcha)
	}
}
