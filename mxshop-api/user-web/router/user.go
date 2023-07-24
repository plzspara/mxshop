package router

import (
	"github.com/gin-gonic/gin"
	api2 "mxshop-api/user-web/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	{
		userRouter.GET("list", api2.GetUserList)
	}
}
