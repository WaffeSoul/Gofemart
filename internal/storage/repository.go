package storage

import "gofemart/internal/model"

type UserRepository interface {
	FindByName(name string) (*model.User, error)
	FindById(id int) (*model.User, error)
	Create(user *model.User) error
	Delete(name string) error
	Migrate() error
}

type OrderRepository interface {
	FindByNumber(number string) (*model.Order, error)
	FindByUserId(id int) (*[]model.Order, error)
	Create(user *model.Order) error
	Delete(name string) error
	Migrate() error
}

type WithdrawRepository interface {
	FindByOrder(order string) (*model.Withdraw, error)
	FindByUserId(id int) (*[]model.Withdraw, error)
	Create(user *model.Withdraw) error
	Delete(name string) error
	Migrate() error
}
