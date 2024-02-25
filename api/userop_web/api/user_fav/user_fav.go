package userfav

import (
	"context"
	"hello_go/mxshop/api/userop_web/api"
	"hello_go/mxshop/api/userop_web/forms"
	"hello_go/mxshop/api/userop_web/global"
	"hello_go/mxshop/api/userop_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUserFavList(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	rsp, err := global.UserFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: int32(userId.(uint)),
	})

	if err != nil {
		zap.S().Errorw("[GetUserFavList] 查询【用户收藏】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 请求商品服务 获取商品信息

	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[BatchGetGoods]: 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}
	goodsList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				goodsList = append(goodsList, map[string]interface{}{
					"good_Id":    good.Id,
					"good_name":  good.Name,
					"good_image": good.GoodsFrontImage,
					"good_price": good.ShopPrice,
				})
			}
		}
	}
	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)

}

func NewUserFav(ctx *gin.Context) {

	userFavForm := forms.UserFavForm{}
	if err := ctx.ShouldBindJSON(&userFavForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")

	_, err := global.UserFavClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId: int32(userId.(uint)),
		GoodsId: userFavForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewShopCart] 添加【用户收藏】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func GetUserFavDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserFavClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		GoodsId: int32(i),
		UserId: int32(userId.(uint)),
		// UserId: 1,
	})

	if err != nil {
		zap.S().Errorw("查询收藏状态失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	
	ctx.JSON(http.StatusOK, gin.H{})
}

func DeleteUserFav(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url error",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.UserFavClient.DelectUserFav(context.Background(), &proto.UserFavRequest{
		UserId: int32(userId.(uint)),
		GoodsId: int32(i),
	})

	if err != nil {
		zap.S().Errorw("[DeleteShopCart]: 删除收藏记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
