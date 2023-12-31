package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)
