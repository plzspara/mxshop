package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"mxshop_srvs/user_srv/model"
	"os"
	"strings"
	"time"
)

func main() {
	pwd := passwordEndCode("qwerty")
	fmt.Println(len(pwd))
	b := verifyPwd("qwerty", pwd)
	fmt.Println(b)
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
		log.Panic(err)
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

var options = &password.Options{
	SaltLen:      16,
	Iterations:   100,
	KeyLen:       32,
	HashFunction: sha512.New,
}

func passwordEndCode(pwd string) (encode string) {
	salt, encodedPwd := password.Encode(pwd, options)
	encode = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	return
}

func verifyPwd(rawPwd string, encode string) bool {
	strs := strings.Split(encode, "$")
	return password.Verify(rawPwd, strs[1], strs[2], options)
}
