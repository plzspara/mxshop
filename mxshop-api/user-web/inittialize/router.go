package inittialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/middlewares"
	"mxshop-api/router"
)

func Routers() *gin.Engine {
	routers := gin.Default()
	routers.Use(middlewares.Cors())
	apiGroup := routers.Group("/u/v1")
	router.InitUserRouter(apiGroup)
	return routers
}
