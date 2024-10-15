package storage

import (
	"gofemart/internal/config"
)

type Store interface {
	Users() UserRepository
	Orders() OrderRepository
	Withdrawals() WithdrawRepository
}

func NewStore(conf *config.Config) Store {
	db := InitDB(conf.DB)
	store := NewDatabase(db)
	if err := MigrateTables(store); err != nil {
		panic("failed to initialize migrate table")
	}
	return store
}
