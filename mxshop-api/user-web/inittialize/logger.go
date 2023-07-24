package inittialize

import "go.uber.org/zap"

func InitLogger() {
	development, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(development)
}
