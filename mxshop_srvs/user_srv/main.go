package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/inittialize"
	"mxshop_srvs/user_srv/proto"
	"net"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "ip地址")
	port := flag.Int("port", 9090, "端口")
	flag.Parse()
	inittialize.InitLogger()
	zap.S().Debugf("启动user服务，端口：%d", *port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServers{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		log.Panic("failed to listen: " + err.Error())
	}
	err = server.Serve(listen)

	if err != nil {
		log.Panic("failed to start grpc: " + err.Error())
	}
}
