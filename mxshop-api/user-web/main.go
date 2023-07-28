package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"mxshop-api/global"
	"mxshop-api/inittialize"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	inittialize.InitLogger()
	inittialize.InitConfig()
	err := inittialize.InitTrans("zh")
	if err != nil {
		log.Panic(err)
	}
	routers := inittialize.Routers()
	inittialize.InitConn()
	go func() {
		zap.S().Debugf("启动user api，端口：%d", global.ServerConfig.Port)
		err = routers.Run(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
		if err != nil {
			zap.S().Panic("启动失败", err)
			return
		}
	}()
	zap.S().Info("启动user api successfully")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
