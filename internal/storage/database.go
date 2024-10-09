package storage

import (
	"context"
	"gofemart/internal/storage/order"
	user "gofemart/internal/storage/user"
	withdrawal "gofemart/internal/storage/withdrawal"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	db *pgxpool.Pool

	users       UserRepository
	orders      OrderRepository
	withdrawals WithdrawRepository
}

func InitDB(addr string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil
		// log.Fatalln("Unable to parse DATABASE_URL:", err)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil
	}
	return conn
}

func NewDatabase(db *pgxpool.Pool) *Database {
	return &Database{
		db:          db,
		users:       user.NewRepository(db),
		orders:      order.NewRepository(db),
		withdrawals: withdrawal.NewRepository(db),
	}
}

func (db *Database) Users() UserRepository {
	return db.users
}

func (db *Database) Orders() OrderRepository {
	return db.orders
}

func (db *Database) Withdrawals() WithdrawRepository {
	return db.withdrawals
}
