package initialize

import (
	"fmt"
	"hello_go/mxshop/api/order_web/global"
	"hello_go/mxshop/api/order_web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	goodsConn, err := grpc.Dial(
		// "consul://127.0.0.1:8500/user_srv?wait=14s",
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] connect goods_srv failed", "msg", err.Error())
		return
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)


	orderConn, err := grpc.Dial(
		// "consul://127.0.0.1:8500/user_srv?wait=14s",
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] connect order_srv failed", "msg", err.Error())
		return
	}

	global.OrderSrvClient = proto.NewOrderClient(orderConn)

	invConn, err := grpc.Dial(
		// "consul://127.0.0.1:8500/user_srv?wait=14s",
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.InvSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] connect inventory_srv failed", "msg", err.Error())
		return
	}

	global.InvSrvClient = proto.NewInventoryClient(invConn)
}
