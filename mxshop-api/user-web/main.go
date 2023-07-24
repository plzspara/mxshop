package main

import (
	"go.uber.org/zap"
	inittialize2 "mxshop-api/user-web/inittialize"
	"strconv"
)

func main() {
	port := 8080
	inittialize2.InitLogger()
	routers := inittialize2.Routers()
	zap.S().Debugf("启动user api，端口：%d", port)
	err := routers.Run(":" + strconv.Itoa(port))
	if err != nil {
		zap.S().Panic("启动失败", err)
		return
	}
}
