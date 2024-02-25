package main

import (
	"flag"
	"fmt"
	"hello_go/mxshop/inventory_srv/global"
	"hello_go/mxshop/inventory_srv/handler"
	"hello_go/mxshop/inventory_srv/initialize"
	"hello_go/mxshop/inventory_srv/proto"
	"hello_go/mxshop/inventory_srv/utils"
	"hello_go/mxshop/inventory_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "127.0.0.1", "ip")
	Port := flag.Int("port", 0, "port")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip: ", *IP)
	zap.S().Info("port:", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)
	if err != nil {
		zap.S().Panic("Service registry failed", err.Error())
	}

	zap.S().Infof("start service, port: %d", *Port)

	go func() {
		err = server.Serve(lis)

		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	if err = register_client.DeRegister(serviceID); err != nil {
		zap.S().Panic("DeRegister faided", err.Error())
	}else{
		zap.S().Info("DeRegister success")
	}

}
