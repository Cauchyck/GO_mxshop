package model

import "time"

type ShoppingCart struct {
	BaseModel
	User   int32 `gorm:"type:int;index"`
	Goods  int32 `gorm:"type:int;index"`
	Nums   int32 `gorm:"type:int;"`
	Checks bool
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct {
	BaseModel
	User    int32  `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30); index"`
	PayType string `gorm:"type:varchar(20) comment 'alipay, wechat'"`

	Status     string `gorm:"type:varchar(20) comment 'PAYING, TRADE_SUCCESS, TRADE_CLOSED, WAIT_BUYER_PAY, TRADE_FINISHED'"`
	TradeNo    string `gorm:"type:varchar(100) comment 'TradeNumber'"`
	OrderMount float32
	PayTime    *time.Time `gorm:"type:datetime"`

	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
	Post         string `gorm:"type:varchar(20)"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

type OrderGoods struct {
	BaseModel
	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	// 字段冗余；镜像；减少服务间调用
	GoodsName  string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
