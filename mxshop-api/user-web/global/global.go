package global

import (
	ut "github.com/go-playground/universal-translator"
	"google.golang.org/grpc"
	"mxshop-api/config"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig
	GrpcClient   *grpc.ClientConn
)
