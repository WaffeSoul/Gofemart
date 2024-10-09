package model

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       int    `json:"-"`
	Username string `json:"login"`
	Password string `json:"password"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int
}

type Balance struct {
	Current  int `json:"current"`
	Withdraw int `json:"withdraw"`
}
