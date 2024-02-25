package initialize

import (
	"encoding/json"
	"fmt"
	"hello_go/mxshop/api/user_web/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {

	v := viper.New()
	v.SetConfigFile("/Volumes/KIOXIARC20/GO/mxshop/api/user_web/config-debug-nacos.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	// fmt.Println(global.ServerConfig)
	zap.S().Infof("config info: %v", global.NacosConfig)

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  clientConfig,
		"serverConfigs": serverConfigs,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	
	if err != nil {
		panic(err)
	}
	// fmt.Println(content)
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("Get nacos config failed: %s", err.Error())
	}
	fmt.Println(global.ServerConfig)

}
