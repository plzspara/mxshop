package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
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

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panic("new consul client failed ", err)
		return
	}

	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ServerConfig.Host, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	serviceRegistration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		Address: global.ServerConfig.Host,
		Port:    global.ServerConfig.Port,
		Tags:    []string{"api", "user-web"},
		ID:      global.ServerConfig.Name,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		zap.S().Panic(err)
		return
	}
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
