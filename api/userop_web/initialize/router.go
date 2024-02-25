package initialize

import (
	"hello_go/mxshop/api/userop_web/middlewares"
	"hello_go/mxshop/api/userop_web/router"

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
	ApiRouter := Router.Group("/op/v1")

	router.InitAddressRouter(ApiRouter)
	router.InitMessageRouter(ApiRouter)
	router.InitUserFavtRouter(ApiRouter)


	return Router

}