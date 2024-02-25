package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	PassWord string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
}


type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Tags       []string       `mapstructure:"tags" json:"tags"`
	Host       string       `mapstructure:"host" json:"host"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo GoodsSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64    `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}
