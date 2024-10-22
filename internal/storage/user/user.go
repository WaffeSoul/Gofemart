package user

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

func (p *Repository) Create(user *model.User) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), createSQL, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (p *Repository) FindByName(name string) (*model.User, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), findByNameSQL, name)
	if err != nil {
		return nil, err
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err == pgx.ErrNoRows {
		return nil, errors.New("no name in db")
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func (p *Repository) FindByID(id int) (*model.User, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), findByIDSQL, id)
	if err != nil {
		return nil, err
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err == pgx.ErrNoRows {
		return nil, errors.New("no id in db")
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
	if err != nil {
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
