package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Id          uint   `json:"id"`
	Nickname    string `json:"nickname"`
	AuthorityId uint   `json:"authority_id"`
	jwt.RegisteredClaims
}
