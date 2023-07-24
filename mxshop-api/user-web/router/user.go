package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	{
		userRouter.GET("list", api.GetUserList)
	}
}
