package user

import (
	"context"
	"gofemart/internal/model"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(addr string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		panic("parse config db")
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic("failed to initialize db")
	}
	return conn
}

func TestRepository_Create(t *testing.T) {
	db := InitDB("postgresql://test:test@127.0.0.1:5433/test?sslmode=disable")
	res := NewRepository(db)
	res.Migrate()
	defer res.Drop()
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				user: &model.User{
					Username: "user1",
					Password: "xxxxxxxxxx",
				},
			},
			wantErr: false,
		},
		{
			name: "Exsist",
			args: args{
				user: &model.User{
					Username: "user1",
					Password: "xxxxxxxxxx",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := res.Create(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_FindByName(t *testing.T) {
	db := InitDB("postgresql://test:test@127.0.0.1:5433/test?sslmode=disable")
	res := NewRepository(db)
	res.Migrate()
	defer res.Drop()
	res.Create(&model.User{
		Username: "asdasdas",
		Password: "asdasdasd",
	})
	type args struct {
		name string
	}
	type want struct {
		orders  *model.User
		wantErr bool
		err     string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name: "OK",
			args: args{
				name: "asdasdas",
			},
			want: want{
				orders: &model.User{
					ID:       1,
					Username: "asdasdas",
					Password: "asdasdasd",
				},
				wantErr: false,
				err:     "",
			},
		},
		{
			name: "No username",
			args: args{
				name: "2",
			},
			want: want{
				orders:  nil,
				wantErr: true,
				err:     "no name in db",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := res.FindByName(tt.args.name)
			if (err != nil) != tt.want.wantErr && tt.want.err != "" {
				if err.Error() != tt.want.err {
					t.Errorf("Repository.FindByNumber() error = %v, wantErr %v", err, tt.want.err)
					return
				}
			}
			if (err != nil) && tt.want.err != "" {
				if err.Error() != tt.want.err {
					t.Errorf("Repository.FindByNumber() error = %v, wantErr %v", err, tt.want.err)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want.orders) {
				t.Errorf("Repository.FindByNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_FindByID(t *testing.T) {
	db := InitDB("postgresql://test:test@127.0.0.1:5433/test?sslmode=disable")
	res := NewRepository(db)
	res.Migrate()
	defer res.Drop()
	res.Create(&model.User{
		Username: "asdasdas",
		Password: "asdasdasd",
	})
	type args struct {
		id int
	}
	type want struct {
		orders  *model.User
		wantErr bool
		err     string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			want: want{
				orders: &model.User{
					ID:       1,
					Username: "asdasdas",
					Password: "asdasdasd",
				},
				wantErr: false,
				err:     "",
			},
		},
		{
			name: "No username",
			args: args{
				id: 2,
			},
			want: want{
				orders:  nil,
				wantErr: true,
				err:     "no id in db",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := res.FindByID(tt.args.id)
			if (err != nil) != tt.want.wantErr && tt.want.err != "" {
				if err.Error() != tt.want.err {
					t.Errorf("Repository.FindByNumber() error = %v, wantErr %v", err, tt.want.err)
					return
				}
			}
			if (err != nil) && tt.want.err != "" {
				if err.Error() != tt.want.err {
					t.Errorf("Repository.FindByNumber() error = %v, wantErr %v", err, tt.want.err)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want.orders) {
				t.Errorf("Repository.FindByNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
