package model

type LeavingMessage struct {
	BaseModel

	User        int32  `gorm:"type:int;index"`
	MessageType int32  `gorm:"type:int comment'留言类型:1,2,3,4,5'"`
	Subject     string `gorm:"type:varchar(100)"`

	Message string
	File    string `gomr:"type:varchar(200)"`
}

func (LeavingMessage) TableName() string{
	return "leavingMessage"
}

type Address struct {
	BaseModel
	
	User        int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City     string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address     string `gorm:"type:varchar(100)"`
	SingerName     string `gorm:"type:varchar(20)"`
	SingerMobile     string `gorm:"type:varchar(11)"`

}

type UserFav struct{
	BaseModel

	User        int32  `gorm:"type:int;index:idx_user_goods,unique"`
	Goods        int32  `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}