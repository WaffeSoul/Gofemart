package service

import (
	"gofemart/internal/jwt"
	"gofemart/internal/storage"
)

type Service struct {
	store      storage.Store
	JwtManager jwt.JWTManager
}

func NewService(store storage.Store, jwtManager jwt.JWTManager) *Service {
	return &Service{
		JwtManager: jwtManager,
		store:      store,
	}
}
