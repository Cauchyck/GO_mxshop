package initialize

import (
	"fmt"
	"hello_go/mxshop/api/userop_web/global"
	"hello_go/mxshop/api/userop_web/proto"

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


	userOpConn, err := grpc.Dial(
		// "consul://127.0.0.1:8500/user_srv?wait=14s",
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UseropSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] connect userop_srv failed", "msg", err.Error())
		return
	}

	global.AddressClient = proto.NewAddressClient(userOpConn)
	global.MessageSrvClient = proto.NewMessageClient(userOpConn)
	global.UserFavClient = proto.NewUserFavClient(userOpConn)

	
}
