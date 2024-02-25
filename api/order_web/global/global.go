package global

import (
	"hello_go/mxshop/api/order_web/config"
	"hello_go/mxshop/api/order_web/proto"
)


var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	
	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient

	OrderSrvClient proto.OrderClient

	InvSrvClient proto.InventoryClient
)