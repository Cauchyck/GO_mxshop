package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Host        string        `mapstructure:"host" json:"host"`
	Tags        []string      `mapstructure:"tags" json:"tags"`
	Port        int           `mapstructure:"port" json:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo  AliSmsConfg   `mapstructure:"sms" json:"sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type AliSmsConfg struct {
	ApiKey     string `mapstruct:"key" json:"key"`
	ApiSecrect string `mapstruct:"secrect" json:"secrect"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}
