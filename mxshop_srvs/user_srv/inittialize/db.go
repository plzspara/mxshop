package inittialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	. "mxshop_srvs/user_srv/global"
	"os"
	"time"
)

func InitDb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ServerConfig.MysqlInfo.User, ServerConfig.MysqlInfo.Password, ServerConfig.MysqlInfo.Host, ServerConfig.MysqlInfo.Port, ServerConfig.MysqlInfo.Db)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 生成的表名不加S
		},
		Logger: newLogger,
	})
	if err != nil {
		log.Panic(err)
	}
}
