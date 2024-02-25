package initialize

import (
	"hello_go/mxshop/api/order_web/middlewares"
	"hello_go/mxshop/api/order_web/router"

	"github.com/gin-gonic/gin"
)

func healthCheckHandler(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
    })
}


func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", healthCheckHandler)

	Router.Use(middlewares.Cors())
	ApiRouter := Router.Group("/o/v1")

	router.InitOrderRouter(ApiRouter)
	router.InitShopCartRouter(ApiRouter)


	return Router

}