package goods

import (
	"context"
	"hello_go/mxshop/api/goods_web/api"
	"hello_go/mxshop/api/goods_web/forms"
	"hello_go/mxshop/api/goods_web/global"
	"hello_go/mxshop/api/goods_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetGoodsList(ctx *gin.Context) {
	request := &proto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}

	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	IsTab := ctx.DefaultQuery("it", "0")
	if IsTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.TopCategory = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.TopCategory = int32(perNumsInt)

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.TopCategory = int32(brandIdInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	zap.S().Debug("Get goods list")
	rsp, err := global.GoodsSrvClient.GoodsList(context.Background(), request)

	if err != nil {
		zap.S().Errorw("[GetGoodsList] Get goods list failed")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// zap.S().Debug("Get user list")

	result := map[string]interface{}{
		"total": rsp.Total,
	}
	goodsList := make([]interface{}, 0)
	for _, value := range rsp.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"is_sale": value.OnSale,
		})
	}
	result["data"] = goodsList
	ctx.JSON(http.StatusOK, result)
}

func NewGoods(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := global.GoodsSrvClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		GoodsDesc:       goodsForm.GoodsDesc,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func GetGoodsDetail(ctx *gin.Context) {
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

func DeleteGoods(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetGoodsStocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Todo Get goods stocks
	return
}

func UpdateGoodsStatus(ctx *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,
		IsNew:  *goodsStatusForm.IsNew,
		OnSale: *goodsStatusForm.OnSale,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Update success",
	})
}

func UpdateGoods(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		GoodsDesc:       goodsForm.GoodsDesc,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Update success",
	})
}
