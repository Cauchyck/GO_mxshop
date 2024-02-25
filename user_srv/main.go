package main

import (
	"flag"
	"fmt"
	"hello_go/mxshop/user_srv/global"
	"hello_go/mxshop/user_srv/handler"
	"hello_go/mxshop/user_srv/initialize"
	"hello_go/mxshop/user_srv/proto"
	"hello_go/mxshop/user_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
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
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("127.0.0.1:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"imooc", "bobby"}
	registration.Address = "127.0.0.1"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)

		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("exit failed")
	}
	zap.S().Info("exit success")
}
