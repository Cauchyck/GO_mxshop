package initialize

import (
	"fmt"
	"hello_go/mxshop/api/user_web/global"
	"hello_go/mxshop/api/user_web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "github.com/mbobakov/grpc-consul-resolver"
)

func InitSrvConn() {
	userConn, err := grpc.Dial(
		// "consul://127.0.0.1:8500/user_srv?wait=14s",
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] connect user_srv failed", "msg", err.Error())
		return
	}
	

	global.UserSrvClient = proto.NewUserClient(userConn)
}

func InitSrvConnOld() {
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Serbvice == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == ""{
		zap.S().Fatal("[InitSrvConn] connect user_srv failed")
		return
	}
	// ip := "127.0.0.1"
	// port := 8888
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] connect user_servicer failed", "msg", err.Error())
		return
	}

	global.UserSrvClient = proto.NewUserClient(userConn)

}
