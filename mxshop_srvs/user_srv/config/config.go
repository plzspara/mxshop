package config

type MysqlConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Db       string `json:"db" mapstructure:"db"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
}

type ConsulConfig struct {
	Host string `json:"host" mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}
type ServerConfig struct {
	Name       string       `json:"name" mapstructure:"name"`
	MysqlInfo  MysqlConfig  `json:"mysql" mapstructure:"mysql"`
	ConsulInfo ConsulConfig `json:"consul" mapstructure:"consul"`
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
