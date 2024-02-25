package main

import (
	"fmt"
	"hello_go/mxshop/api/userop_web/global"
	"hello_go/mxshop/api/userop_web/initialize"
	"hello_go/mxshop/api/userop_web/utils/register/consul"
	myvalidator "hello_go/mxshop/api/userop_web/validator"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	Router := initialize.Routers()
	initialize.InitSrvConn()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
	// }
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
	}

	// 服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)
	if err != nil {
		zap.S().Panic("Service registry failed", err.Error())
	}

	zap.S().Infof("start service, port: %d", global.ServerConfig.Port)

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("start failed", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = register_client.DeRegister(serviceID); err != nil {
		zap.S().Panic("DeRegister faided", err.Error())
	} else {
		zap.S().Info("DeRegister success")
	}

}
