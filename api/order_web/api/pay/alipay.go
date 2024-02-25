package pay

import (
	"context"
	"hello_go/mxshop/api/order_web/global"
	"hello_go/mxshop/api/order_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func Notify(ctx *gin.Context) {
	client, err := alipay.New(global.ServerConfig.AlipayInfo.AppId, global.ServerConfig.AlipayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝url失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(global.ServerConfig.AlipayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	_, err = global.OrderSrvClient.UpdateOrderStstus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	}
	ctx.String(http.StatusOK, "success")

}
