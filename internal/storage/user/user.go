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

func (p *Repository) Create(user *model.User) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `insert into users(username, password) values ($1, $2)`, user.Username, user.Password)
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
	data := &model.User{}
	err = conn.QueryRow(context.Background(), "select * from users where name=$1", name).Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *Repository) FindById(id int) (*model.User, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	data := &model.User{}
	err = conn.QueryRow(context.Background(), "select * from users where id=$1", id).Scan(&data)
	if err != nil {
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
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		username VARCHAR(255) UNIQUE,
		password VARCHAR(255)
);`)
	if err != nil {
		return err
	}
	return nil
}
