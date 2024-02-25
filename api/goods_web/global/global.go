package global

import (
	"hello_go/mxshop/api/goods_web/config"
	"hello_go/mxshop/api/goods_web/proto"
)


var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
)