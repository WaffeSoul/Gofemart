package order

import (
	"context"
	"errors"
	"gofemart/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (p *Repository) Delete(name string) error {

	return nil
}

func (p *Repository) Create(order *model.Order) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	date := time.Now()
	_, err = conn.Exec(context.Background(), `insert into orders(number, user_id, uploaded_at) values ($1, $2, $3)`, order.Number, order.UserId, date)
	if err != nil {
		return err
	}
	return err
}

func (p *Repository) FindByUserId(id int) (*[]model.Order, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	data := &[]model.Order{}
	err = conn.QueryRow(context.Background(), "select * from orders where user_id=$1", id).Scan(data)
	if err == pgx.ErrNoRows {
		return nil, errors.New("no user_id in db")
	} else if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *Repository) FindByNumber(number string) (*model.Order, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	data := &model.Order{}
	err = conn.QueryRow(context.Background(), "select * from orders where number=$1", number).Scan(&data)
	if err == pgx.ErrNoRows {
		return nil, errors.New("no number in db")
	} else if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *Repository) Migrate() error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders(
		number VARCHAR(255)  PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		uploaded_at DATE
	);`)
	return err
}
