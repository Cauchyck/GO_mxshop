package router

import (
	"hello_go/mxshop/api/user_web/api"
	"hello_go/mxshop/api/user_web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func InitUserRouter(Router *gin.RouterGroup){
	
	UserRouter := Router.Group("user")
	zap.S().Info("Init User url")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}
}