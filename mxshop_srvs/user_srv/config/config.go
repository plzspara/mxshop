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
