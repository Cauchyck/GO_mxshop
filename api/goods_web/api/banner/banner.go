package banner

import (
	"context"
	"hello_go/mxshop/api/goods_web/api"
	"hello_go/mxshop/api/goods_web/forms"
	"hello_go/mxshop/api/goods_web/global"
	"hello_go/mxshop/api/goods_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetBannerList(ctx *gin.Context) {
	rsp, err := global.GoodsSrvClient.BannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		result = append(result, map[string]interface{}{
			"id":    value.Id,
			"index": value.Index,
			"image": value.Image,
			"url":   value.Url,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func NewBanner(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Image: bannerForm.Image,
		Url:   bannerForm.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := map[string]interface{}{
		"id":    rsp.Id,
		"index": rsp.Index,
		"url":   rsp.Url,
		"image": rsp.Image,
	}

	ctx.JSON(http.StatusOK, result)
}

func DeleteBanner(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteBanner(context.Background(), &proto.BannerRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateBanner(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(i),
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Update success",
	})
}
