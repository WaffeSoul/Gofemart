package withdrawal

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

func (p *Repository) Create(draw *model.Withdraw) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), createSQL, draw.UserID, draw.OrderNumber, draw.Sum, draw.ProcessedAt)
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
	rows, err := conn.Query(context.Background(), findByOrderSQL, order)
	if err == pgx.ErrNoRows {
		return nil, errors.New("no number in db")
	} else if err != nil {
		return nil, err
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Withdraw])
	if err == pgx.ErrNoRows {
		return nil, errors.New("no user_id in db")
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *Repository) FindByUserID(id int) (*[]model.Withdraw, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), findByUserIDSQL, id)
	if err == pgx.ErrNoRows {
		return nil, errors.New("no user_id in db")
	} else if err != nil {
		return nil, err
	}
	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Withdraw])
	if err == pgx.ErrNoRows {
		return nil, errors.New("no user_id in db")
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *Repository) Migrate() error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		// Add error
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), migrateSQL)
	if err != nil {
		// Add error
		return err
	}
	return nil
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