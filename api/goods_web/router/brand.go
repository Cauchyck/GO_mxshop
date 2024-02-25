package router

import (
	"hello_go/mxshop/api/goods_web/api/brands"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitbrandRouter(Router *gin.RouterGroup) {

	BrandsRouter := Router.Group("brands")
	zap.S().Info("Init brands Router")
	{
		BrandsRouter.GET("", brands.GetBrandList)
		BrandsRouter.POST("",brands.NewBrand)
		BrandsRouter.DELETE("/:id", brands.DeleteBrand)
		BrandsRouter.PUT("/:id",brands.UpdateBrand)

	}

	CategorysBrandsRouter := Router.Group("categorysBrands")
	zap.S().Info("Init categorysBrands Router")
	{
		CategorysBrandsRouter.GET("", brands.GetAllCategoryBrandList)
		CategorysBrandsRouter.POST("",brands.NewCategoryBrand)
		CategorysBrandsRouter.DELETE("/:id", brands.DeleteCategoryBrand)
		CategorysBrandsRouter.PUT("/:id",brands.UpdateCategoryBrand)

	}

}
