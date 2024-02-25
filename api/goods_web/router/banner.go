package router

import (
	"hello_go/mxshop/api/goods_web/api/banner"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitBannerRouter(Router *gin.RouterGroup) {

	BannerRouter := Router.Group("banner")
	zap.S().Info("Init Banner Router")
	{
		BannerRouter.GET("", banner.GetBannerList)
		BannerRouter.POST("", banner.NewBanner)
		BannerRouter.DELETE("/:id", banner.DeleteBanner)
		BannerRouter.PUT("/:id", banner.UpdateBanner)

	}

}
