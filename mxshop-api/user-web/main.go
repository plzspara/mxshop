package main

import (
	"go.uber.org/zap"
	"log"
	"mxshop-api/inittialize"
	"strconv"
)

func main() {
	port := 8080
	inittialize.InitLogger()
	inittialize.InitConfig()
	err := inittialize.InitTrans("zh")
	if err != nil {
		log.Panic(err)
	}
	routers := inittialize.Routers()
	zap.S().Debugf("启动user api，端口：%d", port)
	err = routers.Run(":" + strconv.Itoa(port))
	if err != nil {
		zap.S().Panic("启动失败", err)
		return
	}

}
