package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"goods_srv/global"
	"goods_srv/handler"
	"goods_srv/inittialize"
	"goods_srv/proto"
	"goods_srv/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ip := flag.String("ip", "192.168.153.152", "ip地址")
	port := flag.Int("port", 9090, "端口")
	flag.Parse()

	if *port == 9090 {
		freePort, err := utils.GetFreePort()
		if err != nil {
			log.Panic(err)
		}
		*port = freePort
	}
	inittialize.InitLogger()
	inittialize.InitConfig()
	inittialize.InitDb()

	zap.S().Debugf("启动user服务，端口：%d", *port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServers{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Panic("failed to listen: " + err.Error())
	}

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	config := api.DefaultConfig()
	address := fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		log.Panic(err)
	}

	http := fmt.Sprintf("%s:%d", *ip, *port)
	check := &api.AgentServiceCheck{
		GRPC:                           http,
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		Port:    *port,
		Address: *ip,
		Tags:    []string{"srvs", "user_srv"},
		ID:      global.ServerConfig.Name,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(&registration)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		err = server.Serve(listen)
		if err != nil {
			log.Panic("failed to start grpc: " + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err = client.Agent().ServiceDeregister(global.ServerConfig.Name); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
