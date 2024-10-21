package service

import (
	"context"
	"gofemart/internal/accrual"
	"gofemart/internal/config"
	"gofemart/internal/logger"
	"gofemart/internal/model"
	"gofemart/internal/storage"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestService_SetOrder(t *testing.T) {
	logger.InitLogger(false)
	conf := config.Config{
		DB:      "postgresql://test:test@127.0.0.1:5433/test?sslmode=disable",
		Accrual: "http://localhost:8080",
		Server:  "",
	}
	store := storage.NewStore(&conf)
	store.Users().Create(&model.User{
		Username: "test1",
		Password: "test",
	})
	store.Users().Create(&model.User{
		Username: "test",
		Password: "test",
	})
	store.Orders().Create(&model.Order{
		Number:     "3620637573",
		UserID:     2,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: time.Now().Format("2006-01-02T15:04:05Z"),
	})
	store.Orders().Create(&model.Order{
		Number:     "3637279245",
		UserID:     1,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: time.Now().Format("2006-01-02T15:04:05Z"),
	})
	acc := accrual.NewAccrual(&conf, &store)
	defer acc.Finish()
	ser := NewService(store, acc)
	type args struct {
		userID int
		body   string
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
				userID: 1,
				body:   "3637279245",
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "409",
			args: args{
				userID: 1,
				body:   "3620637573",
			},
			want: want{
				code: 409,
			},
		},
		{
			name: "202",
			args: args{
				userID: 1,
				body:   "7025424594",
			},
			want: want{
				code: 202,
			},
		},
		{
			name: "422",
			args: args{
				userID: 1,
				body:   "7025424596789",
			},
			want: want{
				code: 422,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/api/user/orders", strings.NewReader(tt.args.body))
			ctx := context.WithValue(r.Context(), model.UserIDKey, tt.args.userID)
			w := httptest.NewRecorder()
			got := ser.SetOrder()
			got.ServeHTTP(w, r.WithContext(ctx))
			if !reflect.DeepEqual(w.Result().StatusCode, tt.want.code) {
				t.Errorf("Service.SetOrder() %v = %v, want %v", tt.name, w.Result().StatusCode, tt.want.code)
			}
			w.Result().Body.Close()
			r.Body.Close()
		})
	}
}

func TestService_GetOrders(t *testing.T) {
	logger.InitLogger(false)
	conf := config.Config{
		DB:      "postgresql://test:test@127.0.0.1:5433/test?sslmode=disable",
		Accrual: "http://localhost:8080",
		Server:  "",
	}
	store := storage.NewStore(&conf)
	store.Users().Create(&model.User{
		Username: "test3",
		Password: "test",
	})
	store.Users().Create(&model.User{
		Username: "test4",
		Password: "test",
	})
	store.Orders().Create(&model.Order{
		Number:     "1087831903",
		UserID:     3,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: "2020-12-10T15:15:45+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "6751943355",
		UserID:     3,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	acc := accrual.NewAccrual(&conf, &store)
	defer acc.Finish()
	ser := NewService(store, acc)
	type args struct {
		userID int
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "200 OK",
			args: args{
				userID: 3,
			},
			want: want{
				code: 200,
				body: "[{\"number\":\"1087831903\",\"user_id\":3,\"status\":\"NEW\",\"accrual\":0,\"uploaded_at\":\"2020-12-10T15:15:45+03:00\"},{\"number\":\"6751943355\",\"user_id\":3,\"status\":\"NEW\",\"accrual\":0,\"uploaded_at\":\"2020-12-10T15:12:01+03:00\"}]",
			},
		},
		{
			name: "204",
			args: args{
				userID: 4,
			},
			want: want{
				code: 204,
				body: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/api/user/orders", nil)
			ctx := context.WithValue(r.Context(), model.UserIDKey, tt.args.userID)
			w := httptest.NewRecorder()
			got := ser.GetOrders()
			got.ServeHTTP(w, r.WithContext(ctx))
			if !reflect.DeepEqual(w.Result().StatusCode, tt.want.code) {
				t.Errorf("Service.GetOrders() %v = %v, want %v", tt.name, w.Result().StatusCode, tt.want.code)
			}
			if !reflect.DeepEqual(w.Body.String(), tt.want.body) {
				t.Errorf("Service.GetOrders() %v = %v, want %v", tt.name, w.Body.String(), tt.want.body)
			}
			w.Result().Body.Close()
			r.Body.Close()
		})
	}
}

//Остальные тесты также
