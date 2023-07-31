package config

type UserSrvConfig struct {
	Name string `json:"name"`
}

type Jwt struct {
	Key string `json:"key"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

type Consul struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type NacosInfo struct {
	Ip                  string `mapstructure:"ip"`
	Port                uint64 `mapstructure:"port"`
	NamespaceId         string `mapstructure:"namespaceid"`
	TimeoutMs           uint64 `mapstructure:"timeoutms"`
	NotLoadCacheAtStart bool   `mapstructure:"notloadcacheatstart"`
	LogDir              string `mapstructure:"logdir"`
	CacheDir            string `mapstructure:"cachedir"`
	LogLevel            string `mapstructure:"loglevel"`
}

type NacosConfig struct {
	Nacos NacosInfo `json:"nacos"`
}

type ServerConfig struct {
	Name         string        `json:"name"`
	Port         int           `json:"port"`
	Host         string        `json:"host"`
	UserSrvInfo  UserSrvConfig `json:"user_srv"`
	JwtInfo      Jwt           `json:"jwt"`
	RedisConfig  Redis         `json:"redis"`
	ConsulConfig Consul        `json:"consul"`
}
