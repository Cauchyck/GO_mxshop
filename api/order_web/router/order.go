package router

import (
	"hello_go/mxshop/api/order_web/api/order"
	"hello_go/mxshop/api/order_web/api/pay"
	"hello_go/mxshop/api/order_web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitOrderRouter(Router *gin.RouterGroup) {

	OrderRouter := Router.Group("orders").Use(middlewares.JWTAuth())
	zap.S().Info("Init orders Router")
	{
		OrderRouter.GET("", order.GetOrderList)
		OrderRouter.POST("", order.NewOrder)
		OrderRouter.GET("/:id", order.GetOrderDetail)
		OrderRouter.DELETE("/:id", order.DeleteOrder)
		OrderRouter.PUT("/:id", order.UpdateOrder)

	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}

}
