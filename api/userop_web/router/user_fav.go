package router

import (
	"hello_go/mxshop/api/userop_web/middlewares"
	userfav "hello_go/mxshop/api/userop_web/api/user_fav"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserFavtRouter(Router *gin.RouterGroup) {

	UserFavRouter := Router.Group("userfavs").Use(middlewares.JWTAuth())
	zap.S().Info("Init UserFavs url")
	{
		UserFavRouter.GET("", userfav.GetUserFavList)
		UserFavRouter.POST("", userfav.NewUserFav)
		UserFavRouter.DELETE("/:id", userfav.DeleteUserFav)
		UserFavRouter.GET("/:id", userfav.GetUserFavDetail)

	}

}
