package withdrawal

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
		order *model.Withdraw
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				order: &model.Withdraw{
					OrderNumber: "4154663639",
					UserID:      1,
					Sum:         123123,
					ProcessedAt: "2020-12-10T16:12:01+03:00",
				},
			},
			wantErr: false,
		},
		{
			name: "Exsist",
			args: args{
				order: &model.Withdraw{
					OrderNumber: "4154663639",
					UserID:      1,
					Sum:         123123,
					ProcessedAt: "2020-12-10T16:12:01+03:00",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := res.Create(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_FindByUserID(t *testing.T) {
	db := InitDB("postgresql://test:test@127.0.0.1:5433/test?sslmode=disable")
	res := NewRepository(db)
	res.Migrate()
	defer res.Drop()
	res.Create(&model.Withdraw{
		OrderNumber: "4154663639",
		UserID:      1,
		Sum:         123,
		ProcessedAt: "2020-12-10T15:12:01+03:00",
	})
	res.Create(&model.Withdraw{
		OrderNumber: "1231251234",
		UserID:      1,
		Sum:         123,
		ProcessedAt: "2020-12-10T16:12:01+03:00",
	})
	type args struct {
		userID int
	}
	type want struct {
		orders  *[]model.Withdraw
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
				userID: 1,
			},
			want: want{
				orders: &[]model.Withdraw{
					{
						OrderNumber: "4154663639",
						UserID:      1,
						Sum:         123,
						ProcessedAt: "2020-12-10T15:12:01+03:00",
					},
					{
						OrderNumber: "1231251234",
						UserID:      1,
						Sum:         123,
						ProcessedAt: "2020-12-10T16:12:01+03:00",
					},
				},
				wantErr: false,
				err:     "",
			},
		},
		{
			name: "No user id",
			args: args{
				userID: 2,
			},
			want: want{
				orders:  nil,
				wantErr: true,
				err:     "no user_id in db",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := res.FindByUserID(tt.args.userID)
			if (err != nil) != tt.want.wantErr && tt.want.err != "" {
				if err.Error() != tt.want.err {
					t.Errorf("Repository.FindByUserID() error = %v, wantErr %v", err, tt.want.err)
					return
				}
			}
			if (err != nil) && tt.want.err != "" {
				if err.Error() != tt.want.err {
					t.Errorf("Repository.FindByUserID() error = %v, wantErr %v", err, tt.want.err)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want.orders) {
				t.Errorf("Repository.FindByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_FindByOrder(t *testing.T) {
	db := InitDB("postgresql://test:test@127.0.0.1:5433/test?sslmode=disable")
	res := NewRepository(db)
	res.Migrate()
	defer res.Drop()
	res.Create(&model.Withdraw{
		OrderNumber: "4154663639",
		UserID:      1,
		Sum:         123,
		ProcessedAt: "2020-12-10T15:12:01+03:00",
	})
	res.Create(&model.Withdraw{
		OrderNumber: "1231251234",
		UserID:      1,
		Sum:         123,
		ProcessedAt: "2020-12-10T16:12:01+03:00",
	})
	type args struct {
		number string
	}
	type want struct {
		orders  *model.Withdraw
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
				number: "4154663639",
			},
			want: want{
				orders: &model.Withdraw{
					OrderNumber: "4154663639",
					UserID:      1,
					Sum:         123,
					ProcessedAt: "2020-12-10T15:12:01+03:00",
				},
				wantErr: false,
				err:     "",
			},
		},
		{
			name: "No  number",
			args: args{
				number: "2",
			},
			want: want{
				orders:  nil,
				wantErr: true,
				err:     "no number in db",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := res.FindByOrder(tt.args.number)
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
