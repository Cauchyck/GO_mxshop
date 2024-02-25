package address

import (
	"context"
	"hello_go/mxshop/api/userop_web/api"
	"hello_go/mxshop/api/userop_web/forms"
	"hello_go/mxshop/api/userop_web/global"
	"hello_go/mxshop/api/userop_web/models"
	"hello_go/mxshop/api/userop_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetAddressList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims")

	request := proto.AddressRequest{}

	model := claims.(*models.CustomClaims)
	if model.AuthorityId != 2 {
		userId, _ := ctx.Get("userId")
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.AddressClient.GetAddressList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	addressList := make([]interface{}, 0)

	for _, item := range rsp.Data {
		addressList = append(addressList, map[string]interface{}{
			"id":            item.Id,
			"user_id":       item.UserId,
			"province":      item.Province,
			"city":          item.City,
			"district":      item.District,
			"address":       item.Address,
			"signer_name":   item.SignerName,
			"signer_mobiel": item.SignerMobile,
		})
	}
	reMap["data"] = addressList
	ctx.JSON(http.StatusOK, reMap)

}

func NewAddress(ctx *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	userId, _ := ctx.Get("userId")

	rsp, err := global.AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       int32(userId.(uint)),
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})

	if err != nil {
		zap.S().Errorw("新建失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})

}
func UpdateAddress(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url error",
		})
		return
	}

	addressForm := forms.AddressForm{}

	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	request := proto.AddressRequest{
		Id: int32(i),
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	}

	_, err = global.AddressClient.UpdateAddress(context.Background(), &request)

	if err != nil {
		zap.S().Errorw("[UpdateShopCart]: 更新地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func DeleteAddress(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url error",
		})
		return
	}

	_, err = global.AddressClient.DeletcAddress(context.Background(), &proto.AddressRequest{
		Id: int32(i),
	})

	if err != nil {
		zap.S().Errorw("[DeleteAddress]: 删地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
