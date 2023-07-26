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
		userRouter.POST("pwd_login", api.PasswordLogin)
		userRouter.POST("mobile_login", api.MobileLogin)
		userRouter.POST("register", api.Register)
		userRouter.GET("captcha", api.GetCaptcha)
		userRouter.GET("send_sms", api.SendSms)
	}
}
