package service

import (
	"gofemart/internal/config"
	"gofemart/internal/crypto"
	"gofemart/internal/logger"
	"gofemart/internal/model"
	"gofemart/internal/storage"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestService_SignIn(t *testing.T) {
	logger.InitLogger(false)
	conf := config.Config{
		DB:      "postgresql://test:test@127.0.0.1:5433/test?sslmode=disable",
		Accrual: "http://localhost:8080",
		Server:  "",
	}
	store := storage.NewStore(&conf)
	defer storage.DropTable(store)
	hashedPassword, _ := crypto.HashPassword("test")
	store.Users().Create(&model.User{
		Username: "test",
		Password: hashedPassword,
	})
	ser := NewService(store, nil)
	type args struct {
		body string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "200 OK",
			args: args{
				body: `{"login":"test","password":"test"}`,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "400",
			args: args{
				body: `{"usernameaaaaa":"user","password":"XXXX"}`,
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "401",
			args: args{
				body: `{"login":"test","password":"XXXX"}`,
			},
			want: want{
				code: 401,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/api/user/login", strings.NewReader(tt.args.body))
			w := httptest.NewRecorder()
			got := ser.SignIn()
			got.ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()
			if !reflect.DeepEqual(resp.StatusCode, tt.want.code) {
				t.Errorf("Service.SignUp() %v = %v, want %v", tt.name, resp.StatusCode, tt.want.code)
			}
		})
	}
}

func TestService_SignUp(t *testing.T) {
	logger.InitLogger(false)
	conf := config.Config{
		DB:      "postgresql://test:test@127.0.0.1:5433/test?sslmode=disable",
		Accrual: "http://localhost:8080",
		Server:  "",
	}
	store := storage.NewStore(&conf)
	defer storage.DropTable(store)
	store.Users().Create(&model.User{
		Username: "test",
		Password: "test",
	})
	ser := NewService(store, nil)
	type args struct {
		body string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "200 OK",
			args: args{
				body: `{"login":"user","password":"XXXX"}`,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "400",
			args: args{
				body: `{"usernameaaaaa":"user","password":"XXXX"}`,
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "409",
			args: args{
				body: `{"login":"test","password":"XXXX"}`,
			},
			want: want{
				code: 409,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/api/user/register", strings.NewReader(tt.args.body))
			w := httptest.NewRecorder()
			got := ser.SignUp()
			got.ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()
			if !reflect.DeepEqual(resp.StatusCode, tt.want.code) {
				t.Errorf("Service.SignUp() %v = %v, want %v", tt.name, resp.StatusCode, tt.want.code)
			}
		})
	}
}
