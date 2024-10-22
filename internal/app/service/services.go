package service

import (
	"gofemart/internal/accrual"
	"gofemart/internal/jwt"
	"gofemart/internal/storage"
)

type Service struct {
	store      storage.Store
	JwtManager jwt.JWTManager
	accrual    *accrual.Accrual
}

func NewService(store storage.Store, acc *accrual.Accrual) *Service {
	jwtManager, err := jwt.NewJWTManager("sfjvpasasdf", "30m", "24h")
	if err != nil {
		panic("failed to initialize jwtManager")
	}
	return &Service{
		JwtManager: *jwtManager,
		store:      store,
		accrual:    acc,
	}
}
