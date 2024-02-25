package initialize

import (
	"hello_go/mxshop/api/user_web/middlewares"
	"hello_go/mxshop/api/user_web/router"

	"github.com/gin-gonic/gin"
)

func healthCheckHandler(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
    })
}


func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())

	Router.GET("/health", healthCheckHandler)
	
	ApiRouter := Router.Group("/u/v1")
	router.InitUserRouter(ApiRouter)
	router.InitBaseRouter(ApiRouter)

	return Router

}