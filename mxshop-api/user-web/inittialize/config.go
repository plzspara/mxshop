package inittialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
	err = v.Unmarshal(&global.NacosConfig)
	if err != nil {
		log.Panic(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Infow("配置文件产生变化：", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		zap.S().Infow("配置信息：", global.ServerConfig)
	})

	nacos := global.NacosConfig.Nacos
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacos.NamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           nacos.TimeoutMs,
		NotLoadCacheAtStart: nacos.NotLoadCacheAtStart,
		LogDir:              nacos.LogDir,
		CacheDir:            nacos.CacheDir,
		LogLevel:            nacos.LogLevel,
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacos.Ip,
			Port:   nacos.Port,
		},
	}

	iClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.S().Panic(err)
		return
	}
	config, err := iClient.GetConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev",
	})
	if err != nil {
		log.Panic(err)
		return
	}

	err = json.Unmarshal([]byte(config), &global.ServerConfig)
	if err != nil {
		log.Panic(err)
	}
	zap.S().Info(global.ServerConfig)
}
