package router

import (
	"hello_go/mxshop/api/userop_web/middlewares"
	"hello_go/mxshop/api/userop_web/api/message"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitMessageRouter(Router *gin.RouterGroup) {

	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	zap.S().Info("Init Messages Router")
	{
		MessageRouter.GET("", message.GetMessageList)
		MessageRouter.POST("", message.NewMessage)

	}

}
