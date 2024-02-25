package brands

import (
	"context"
	"hello_go/mxshop/api/goods_web/api"
	"hello_go/mxshop/api/goods_web/forms"
	"hello_go/mxshop/api/goods_web/global"
	"hello_go/mxshop/api/goods_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBrandList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.GoodsSrvClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}

	result := make([]interface{}, 0)

	for _, value := range rsp.Data[pnInt : pnInt*pSizeInt+pSizeInt] {
		result = append(result, map[string]interface{}{
			"id":   value.Id,
			"name": value.Name,
			"logo": value.Logo,
		})
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":  result,
		"total": rsp.Total,
	})
}

func NewBrand(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandForm.Name,
		Logo: brandForm.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := map[string]interface{}{
		"id":   rsp.Id,
		"name": rsp.Name,
		"Logo": rsp.Logo,
	}

	ctx.JSON(http.StatusOK, result)
}

func DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteBrand(context.Background(), &proto.BrandRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateBrand(ctx *gin.Context) {
	brandForm := forms.BrandForm{}
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   int32(i),
		Logo: brandForm.Logo,
		Name: brandForm.Name,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Update success",
	})
}

func GetCategoryBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}

	result := make([]interface{}, 0)

	for _, value := range rsp.Data {
		result = append(result, map[string]interface{}{
			"id":   value.Id,
			"name": value.Name,
			"logo": value.Logo,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func GetAllCategoryBrandList(ctx *gin.Context) {
	rsp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}

	result := make([]interface{}, 0) 

	for _, value := range rsp.Data {
		result = append(result, map[string]interface{}{
			"id":   value.Id,
			"category": map[string]interface{}{
				"id": value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id": value.Brand.Id,
				"logo": value.Brand.Logo,
				"name": value.Brand.Name,
			},
		})
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"total": rsp.Total,
		"data": result,
	})
}

func NewCategoryBrand(ctx *gin.Context) {
	categorybrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categorybrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}



	rsp, err := global.GoodsSrvClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: int32(categorybrandForm.CategoryId),
		BrandId: int32(categorybrandForm.BrandId),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": rsp.Id,
	})
}

func UpdateCategoryBrand(ctx *gin.Context) {
	categorybrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categorybrandForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	id := ctx.Param("id")

	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id: int32(i),
		CategoryId: int32(categorybrandForm.CategoryId),
		BrandId: int32(categorybrandForm.BrandId),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteCategoryBrand(ctx *gin.Context) {
	id := ctx.Param("id")

	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id: int32(i),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}