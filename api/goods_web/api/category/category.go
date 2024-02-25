package category

import (
	"context"
	"encoding/json"
	"hello_go/mxshop/api/goods_web/api"
	"hello_go/mxshop/api/goods_web/forms"
	"hello_go/mxshop/api/goods_web/global"
	"hello_go/mxshop/api/goods_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetAllCategoryList(ctx *gin.Context) {
	r, err := global.GoodsSrvClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[GetAllCategoryList]: faild", err.Error())
	}

	ctx.JSON(http.StatusOK, data)
}

func NewCategory(ctx *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name: categoryForm.Name,
		IsTab: *categoryForm.IsTab,
		Level: categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := map[string]interface{}{
		"id": rsp.Id,
		"name": rsp.Name,
		"parent": rsp.ParentCategory,
		"level": rsp.Level,
		"is_tab": rsp.IsTab,
	}

	ctx.JSON(http.StatusOK, result)
}

func GetCategoryDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}

	subCategorys := make([]interface{}, 0)

	r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	for _, value := range r.SubCategorys {
		subCategorys = append(subCategorys, map[string]interface{}{
			"id":              value.Id,
			"name":            value.Name,
			"level":           value.Level,
			"parent_category": value.ParentCategory,
			"is_tab":          value.IsTab,
		})
	}
	rsp := map[string]interface{}{
		"id":          r.Info.Id,
		"name":        r.Info.Name,
		"level": r.Info.Level,
		"parent_category":   r.Info.ParentCategory,
		"is_tab":      r.Info.IsTab,
		"sub_categorys": subCategorys,
	}


	ctx.JSON(http.StatusOK, rsp)
}

func DeleteGoods(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}


func UpdateCategory(ctx *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request :=&proto.CategoryInfoRequest{
		Id:              int32(i),
		Name:           categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}

	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Update success",
	})
}
