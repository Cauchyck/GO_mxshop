package router

import (
	"hello_go/mxshop/api/goods_web/api/goods"
	"hello_go/mxshop/api/goods_web/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func InitGoodsRouter(Router *gin.RouterGroup) {

	GoodsRouter := Router.Group("goods")
	zap.S().Info("Init goods url")
	{
		GoodsRouter.GET("list", goods.GetGoodsList)
		GoodsRouter.POST("new", middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.NewGoods)
		GoodsRouter.GET("detail/:id", goods.GetGoodsDetail)
		GoodsRouter.DELETE("delete/:id", middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.DeleteGoods)
		GoodsRouter.GET("detail/:id/stocks", goods.GetGoodsStocks)
		GoodsRouter.PATCH("detail/:id/update", middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.UpdateGoods)
		GoodsRouter.PATCH("detail/:id", middlewares.JWTAuth(),middlewares.IsAdminAuth(),goods.UpdateGoodsStatus)
	}
	
} 


