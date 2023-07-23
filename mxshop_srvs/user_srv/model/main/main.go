package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"mxshop_srvs/user_srv/model"
	"os"
	"time"
)

func main() {
	createTable()
}

func createTable() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 生成的表名不加S
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
