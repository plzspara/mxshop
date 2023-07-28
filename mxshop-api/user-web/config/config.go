package config

type UserSrvConfig struct {
	Name string `json:"name"`
}

type Jwt struct {
	Key string `mapstructure:"key"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type Consul struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	UserSrvInfo  UserSrvConfig `mapstructure:"user_srv"`
	JwtInfo      Jwt           `mapstructure:"jwt"`
	RedisConfig  Redis         `mapstructure:"redis"`
	ConsulConfig Consul        `mapstructure:"consul"`
}
