package user

import (
	"context"
	"gofemart/internal/model"

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

func (p *Repository) Create(draw *model.Withdraw) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `insert into withdrawals(user_id, order_mumber,sum,processed_at) values ($1, $2)`, draw.UserId, draw.Order, draw.Sum, draw.ProcessedAt)
	if err != nil {
		// Add error
		return err
	}
	return nil
}

func (p *Repository) FindByOrder(order string) (*model.Withdraw, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return nil, err
	}
	defer conn.Release()
	data := &model.Withdraw{}
	err = conn.QueryRow(context.Background(), "select * from withdrawals where order_mumber=$1", order).Scan(&data)
	if err != nil {
		// Add error
		return nil, err
	}
	return data, nil
}

func (p *Repository) FindByUserId(id int) (*[]model.Withdraw, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return nil, err
	}
	defer conn.Release()
	data := &[]model.Withdraw{}
	err = conn.QueryRow(context.Background(), "select * from withdrawals where user_id=$1", id).Scan(&data)
	if err != nil {
		// Add error
		return nil, err
	}
	return data, nil
}

func (p *Repository) Migrate() error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS withdrawals(
		order_mumber VARCHAR(255) PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		sum INTEGER,
		processed_at DATE
	);`)
	if err != nil {
		// Add error
		return err
	}
	return nil
}
