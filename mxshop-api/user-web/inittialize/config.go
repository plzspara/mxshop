package inittialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"mxshop-api/global"
)

func InitConfig() {
	configPrefix := "debug"
	filePath := fmt.Sprintf("user-web/config-%s.yaml", configPrefix)

	v := viper.New()
	v.SetConfigFile(filePath)
	err := v.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	err = v.Unmarshal(&global.ServerConfig)
	if err != nil {
		log.Panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infow("配置文件产生变化：", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		zap.S().Infow("配置信息：", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port)
	})

}
