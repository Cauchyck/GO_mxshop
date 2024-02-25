package initialize

import (
	"fmt"
	"hello_go/mxshop/order_srv/global"
	"hello_go/mxshop/order_srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvs() {
	consulInfo := global.ServerConfig.ConsulInfo
	invConn, err := grpc.Dial(
		// "consul://127.0.0.1:8500/inventory_srv?wait=14s",
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[InitSrvs] connect inventory_srv failed", "msg", err.Error())
		return
	}

	global.InventorySrvClient = proto.NewInventoryClient(invConn)

	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Errorw("[InitSrvs] connect goods_srv failed", "msg", err.Error())
		return
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

}
