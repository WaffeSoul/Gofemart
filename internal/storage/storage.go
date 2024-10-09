package storage

type Store interface {
	Users() UserRepository
	Orders() OrderRepository
	Withdrawals() WithdrawRepository
}
