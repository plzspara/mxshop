package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"mxshop-api/global"
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
