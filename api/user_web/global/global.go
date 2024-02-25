package global

import (
	"hello_go/mxshop/api/user_web/config"
	"hello_go/mxshop/api/user_web/proto"
)


var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)