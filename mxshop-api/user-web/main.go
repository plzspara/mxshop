package main

import (
	"go.uber.org/zap"
	"log"
	"mxshop-api/global"
	"mxshop-api/inittialize"
	"strconv"
)

func main() {
	inittialize.InitLogger()
	inittialize.InitConfig()
	err := inittialize.InitTrans("zh")
	if err != nil {
		log.Panic(err)
	}
	routers := inittialize.Routers()
	zap.S().Debugf("启动user api，端口：%d", global.ServerConfig.Port)
	err = routers.Run(":" + strconv.Itoa(global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败", err)
		return
	}

}
