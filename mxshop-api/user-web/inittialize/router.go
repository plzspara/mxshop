package inittialize

import (
	"github.com/gin-gonic/gin"
	router2 "mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	routers := gin.Default()
	apiGroup := routers.Group("/u/v1")
	router2.InitUserRouter(apiGroup)
	return routers
}
