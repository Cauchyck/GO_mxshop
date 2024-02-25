package router

import (
	"hello_go/mxshop/api/goods_web/api/category"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitCategoryRouter(Router *gin.RouterGroup) {

	CategorysRouter := Router.Group("categorys")
	zap.S().Info("Init Categorys Router")
	{
		CategorysRouter.GET("", category.GetAllCategoryList)
		CategorysRouter.POST("", category.NewCategory)
		CategorysRouter.GET("/:id", category.GetCategoryDetail)
		CategorysRouter.DELETE("/:id", category.DeleteGoods)
		CategorysRouter.PUT("/:id", category.UpdateCategory)

	}

}
