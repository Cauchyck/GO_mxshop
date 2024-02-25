package initialize

import (
	"hello_go/mxshop/api/goods_web/middlewares"
	"hello_go/mxshop/api/goods_web/router"

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
	ApiRouter := Router.Group("/g/v1")

	router.InitGoodsRouter(ApiRouter)
	router.InitCategoryRouter(ApiRouter)
	router.InitBannerRouter(ApiRouter)
	router.InitbrandRouter(ApiRouter)

	return Router

}