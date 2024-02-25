package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;not null;index" json:"goods"`
	Stocks  int32 `gorm:"type:int;not null" json:"stocks"`
	Version int32 `gorm:"type:int;not null" json:"version"`
}

// type InventoryHistory struct {
// 	User int32 `gorm:"type:int;not null;index" json:"user"`
// 	Goods int32`gorm:"type:int;not null;index" json:"goods"`
// 	Nums int32 `gorm:"type:int;not null;index" json:"nums"`
// 	Order int32 `gorm:"type:int;not null;index" json:"order"`
// 	Status int32 `gorm:"type:int;not null;index" json:"status"`
// }