package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/config"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig
)
