package order

import (
	"context"
	"errors"
	"gofemart/internal/model"

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
	_, err = conn.Exec(context.Background(), createSQL, order.Number, order.UserID, order.Status, order.Accrual, order.UploadedAt)
	if err != nil {
		return err
	}
	return err
}

func (p *Repository) Update(order *model.Order) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), updateSQL, order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}
	return err
}

func (p *Repository) FindByUserID(id int) (*[]model.Order, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), findByIDUserSQL, id)
	if err == pgx.ErrNoRows {
		return nil, errors.New("no user_id in db")
	} else if err != nil {
		return nil, err
	}
	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err == pgx.ErrNoRows {
		return nil, errors.New("no user_id in db")
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *Repository) FindByNumber(number string) (*model.Order, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), findByNumberSQL, number)
	if err == pgx.ErrNoRows {
		return nil, errors.New("no number in db")
	} else if err != nil {
		return nil, err
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Order])
	if err == pgx.ErrNoRows {
		return nil, errors.New("no number in db")
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *Repository) Migrate() error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), migrateSQL)
	return err
}

func (p *Repository) Drop() error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), dropSQL)
	return err
}