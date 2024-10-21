package model

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID       int    `json:"-"`
	Username string `json:"login"`
	Password string `json:"password"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type Balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}

type UserKey int

const (
	UserIDKey UserKey = iota
)
