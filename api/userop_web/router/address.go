package router

import (
	"hello_go/mxshop/api/userop_web/middlewares"
	"hello_go/mxshop/api/userop_web/api/address"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitAddressRouter(Router *gin.RouterGroup) {

	AddressRouter := Router.Group("address").Use(middlewares.JWTAuth())
	zap.S().Info("Init address url")
	{
		AddressRouter.GET("", address.GetAddressList)
		AddressRouter.POST("", address.NewAddress)
		AddressRouter.DELETE("/:id", address.DeleteAddress)
		AddressRouter.PUT("/:id", address.UpdateAddress)

	}

}

