package router

import (
	shopcart "hello_go/mxshop/api/order_web/api/shop_cart"
	"hello_go/mxshop/api/order_web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitShopCartRouter(Router *gin.RouterGroup) {

	ShopCartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	zap.S().Info("Init goods url")
	{
		ShopCartRouter.GET("", shopcart.GetShopCartList)
		ShopCartRouter.POST("", shopcart.NewShopCart)
		ShopCartRouter.DELETE("/:id", shopcart.DeleteShopCart)
		ShopCartRouter.PATCH("/:id", shopcart.UpdateShopCart)

	}

}
