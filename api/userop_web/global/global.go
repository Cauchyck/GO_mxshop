package global

import (
	"hello_go/mxshop/api/userop_web/proto"
	"hello_go/mxshop/api/userop_web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	MessageSrvClient proto.MessageClient

	AddressClient proto.AddressClient

	UserFavClient proto.UserFavClient
	
	GoodsSrvClient proto.GoodsClient
)
