package storage

import "gofemart/internal/model"

type UserRepository interface {
	FindByName(name string) (*model.User, error)
	FindByID(id int) (*model.User, error)
	Create(user *model.User) error
	Delete(name string) error
	Migrate() error
	Drop() error
}

type OrderRepository interface {
	FindByNumber(number string) (*model.Order, error)
	FindByUserID(id int) (*[]model.Order, error)
	Create(order *model.Order) error
	Delete(name string) error
	Update(order *model.Order) error
	Migrate() error
	Drop() error
}

type WithdrawRepository interface {
	FindByOrder(order string) (*model.Withdraw, error)
	FindByUserID(id int) (*[]model.Withdraw, error)
	Create(user *model.Withdraw) error
	Delete(name string) error
	Migrate() error
	Drop() error
}
