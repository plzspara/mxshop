package inittialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop-api/global"
)

func InitConn() {
	config := api.DefaultConfig()
	consul := global.ServerConfig.ConsulConfig
	config.Address = fmt.Sprintf("%s:%d", consul.Host, consul.Port)
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panic(err)
	}
	http := fmt.Sprintf("http://%s:%d/health", global.ServerConfig.Host, global.ServerConfig.Port)
	check := &api.AgentServiceCheck{
		HTTP:                           http,
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	registration := &api.AgentServiceRegistration{
		ID:      global.ServerConfig.Name,
		Name:    global.ServerConfig.Name,
		Tags:    []string{"api", "user-web"},
		Port:    global.ServerConfig.Port,
		Address: global.ServerConfig.Host,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Panic(err)
	}
	filter := fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name)
	services, err := client.Agent().ServicesWithFilter(filter)
	if err != nil {
		zap.S().Panic(err)
	}
	grpcHost := ""
	grpcPort := 0
	for s := range services {
		grpcHost = services[s].Address
		grpcPort = services[s].Port
	}
	global.GrpcClient, err = grpc.Dial(fmt.Sprintf("%s:%d", grpcHost, grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panic(err)
	}

}
