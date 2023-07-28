package inittialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/middlewares"
	"mxshop-api/router"
	"net/http"
)

func Routers() *gin.Engine {
	routers := gin.Default()
	routers.Use(middlewares.Cors())
	routers.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	})
	apiGroup := routers.Group("/u/v1")
	router.InitUserRouter(apiGroup)
	return routers
}
