package shopcart

import (
	"context"
	"hello_go/mxshop/api/order_web/api"
	"hello_go/mxshop/api/order_web/forms"
	"hello_go/mxshop/api/order_web/global"
	"hello_go/mxshop/api/order_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetShopCartList(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})

	if err != nil {
		zap.S().Errorw("[GetShopCartList] 查询【购物车列表】失败")
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
					"id":         item.Id,
					"good_Id":    good.Id,
					"good_name":  good.Name,
					"good_image": good.GoodsFrontImage,
					"good_price": good.ShopPrice,
					"nums":       item.Nums,
					"checked":    item.Checked,
				})
			}
		}
	}
	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)

}

func NewShopCart(ctx *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewShopCart] 查询【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	invRsp, err := global.InvSrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[NewShopCart] 查询【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "库存不足",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})

	if err != nil {
		zap.S().Errorw("[NewShopCart] 添加到购物车失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func GetShopCartDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: int32(i),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"category": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"is_sale": r.OnSale,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func UpdateShopCart(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url error",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}

	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	request := proto.CartItemRequest{
		UserId: int32(userId.(uint)),
		GoodsId: int32(i),
		Nums: itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	_, err  = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)

	if err != nil {
		zap.S().Errorw("[UpdateShopCart]: 更新购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
}

func DeleteShopCart(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url error",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId: int32(userId.(uint)),
		GoodsId: int32(i),
	})

	if err != nil {
		zap.S().Errorw("[DeleteShopCart]: 删除购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
