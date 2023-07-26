package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"mxshop-api/global"
	"strconv"
	"time"
)

type CustomClaims struct {
	Id          uint   `json:"id"`
	Nickname    string `json:"nickname"`
	AuthorityId uint   `json:"authority_id"`
	jwt.RegisteredClaims
}

func NewJwt(claims CustomClaims, b []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(b)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.ServerConfig.JwtInfo.Key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("CustomClaims error")
	}
	return claims, nil
}

func GetConnClient() (*grpc.ClientConn, error) {
	fmt.Println(global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorf("Error grpc Dial error: %v", err.Error())
		return nil, err
	}
	return conn, nil
}

func GetRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		Password: global.ServerConfig.RedisConfig.Password, // no password set
		DB:       global.ServerConfig.RedisConfig.Db,       // use default DB
	})
	return rdb
}

func GetMsmCode() string {
	n := rand.Int63n(time.Now().UnixNano()) + 123789
	n = n << 2
	return strconv.FormatInt(n%1000000, 10)
}
