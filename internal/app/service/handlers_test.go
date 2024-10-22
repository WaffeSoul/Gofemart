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
	defer storage.DropTable(store)
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
	defer storage.DropTable(store)
	store.Users().Create(&model.User{
		Username: "test",
		Password: "test",
	})
	store.Users().Create(&model.User{
		Username: "test1",
		Password: "test",
	})
	store.Orders().Create(&model.Order{
		Number:     "1087831903",
		UserID:     1,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: "2020-12-10T15:15:45+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "6751943355",
		UserID:     1,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	ser := NewService(store, nil)
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
				userID: 1,
			},
			want: want{
				code: 200,
				body: "[{\"number\":\"1087831903\",\"user_id\":1,\"status\":\"NEW\",\"accrual\":0,\"uploaded_at\":\"2020-12-10T15:15:45+03:00\"},{\"number\":\"6751943355\",\"user_id\":1,\"status\":\"NEW\",\"accrual\":0,\"uploaded_at\":\"2020-12-10T15:12:01+03:00\"}]",
			},
		},
		{
			name: "204",
			args: args{
				userID: 2,
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
			resp := w.Result()
			defer resp.Body.Close()
			if !reflect.DeepEqual(resp.StatusCode, tt.want.code) {
				t.Errorf("Service.GetOrders() %v = %v, want %v", tt.name, resp.StatusCode, tt.want.code)
			}
			if !reflect.DeepEqual(w.Body.String(), tt.want.body) {
				t.Errorf("Service.GetOrders() %v = %v, want %v", tt.name, w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestService_GetBalance(t *testing.T) {
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
	store.Users().Create(&model.User{
		Username: "test1",
		Password: "test",
	})
	store.Users().Create(&model.User{
		Username: "test2",
		Password: "test",
	})
	store.Orders().Create(&model.Order{
		Number:     "1087831903",
		UserID:     1,
		Status:     "PROCESSED",
		Accrual:    123123,
		UploadedAt: "2020-12-10T15:15:45+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "6751943355",
		UserID:     2,
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "0210087987",
		UserID:     2,
		Status:     "PROCESSED",
		Accrual:    123123,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "4154663639",
		UserID:     2,
		Status:     "PROCESSED",
		Accrual:    123123,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	store.Withdrawals().Create(&model.Withdraw{
		UserID:      2,
		OrderNumber: "4154663639",
		Sum:         1123.23,
		ProcessedAt: "2020-12-10T15:12:01+03:00",
	})
	ser := NewService(store, nil)
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
				userID: 2,
			},
			want: want{
				code: 200,
				body: `{"current":245122.77,"withdrawn":1123.23}`,
			},
		},
		{
			name: "200 OK zero",
			args: args{
				userID: 3,
			},
			want: want{
				code: 200,
				body: `{"current":0,"withdrawn":0}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/api/user/balance", nil)
			ctx := context.WithValue(r.Context(), model.UserIDKey, tt.args.userID)
			w := httptest.NewRecorder()

			got := ser.GetBalance()
			got.ServeHTTP(w, r.WithContext(ctx))
			if !reflect.DeepEqual(w.Result().StatusCode, tt.want.code) {
				t.Errorf("Service.GetBalance() %v = %v, want %v", tt.name, w.Result().StatusCode, tt.want.code)
			}
			if !reflect.DeepEqual(w.Body.String(), tt.want.body) {
				t.Errorf("Service.GetBalance() %v = %v, want %v", tt.name, w.Body.String(), tt.want.body)
			}
		})
	}
}

func TestService_Withdraw(t *testing.T) {
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
	store.Orders().Create(&model.Order{
		Number:     "1087831903",
		UserID:     1,
		Status:     "PROCESSED",
		Accrual:    123123,
		UploadedAt: "2020-12-10T15:15:45+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "0210087987",
		UserID:     1,
		Status:     "PROCESSED",
		Accrual:    123123,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	store.Orders().Create(&model.Order{
		Number:     "4154663639",
		UserID:     1,
		Status:     "PROCESSED",
		Accrual:    123123,
		UploadedAt: "2020-12-10T15:12:01+03:00",
	})
	store.Withdrawals().Create(&model.Withdraw{
		UserID:      1,
		OrderNumber: "4154663639",
		Sum:         1123.23,
		ProcessedAt: "2020-12-10T15:12:01+03:00",
	})
	ser := NewService(store, nil)
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
				body: `{
					"order": "3079027466",
					"sum": 1123.23
				}`,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "402",
			args: args{
				userID: 1,
				body: `{
					"order": "4154663639",
					"sum": 123123123123123123
				}`,
			},
			want: want{
				code: 402,
			},
		},
		{
			name: "422",
			args: args{
				userID: 1,
				body: `{
					"order": "112315345",
					"sum": 1123.23
				}`,
			},
			want: want{
				code: 422,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/api/user/balance/withdraw", strings.NewReader(tt.args.body))
			r.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(r.Context(), model.UserIDKey, tt.args.userID)
			w := httptest.NewRecorder()
			got := ser.Withdraw()
			got.ServeHTTP(w, r.WithContext(ctx))
			if !reflect.DeepEqual(w.Result().StatusCode, tt.want.code) {
				t.Errorf("Service.Withdraw() %v = %v, want %v", tt.name, w.Result().StatusCode, tt.want.code)
			}
		})
	}
}

func TestService_Withdrawals(t *testing.T) {
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
	store.Users().Create(&model.User{
		Username: "test",
		Password: "test",
	})
	store.Withdrawals().Create(&model.Withdraw{
		UserID:      1,
		OrderNumber: "4154663639",
		Sum:         1123.23,
		ProcessedAt: "2020-12-10T15:12:01+03:00",
	})
	store.Withdrawals().Create(&model.Withdraw{
		UserID:      1,
		OrderNumber: "4180189435",
		Sum:         1123.23,
		ProcessedAt: "2020-13-10T15:12:01+03:00",
	})
	ser := NewService(store, nil)
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
				userID: 1,
			},
			want: want{
				code: 200,
				body: `[{"order":"4154663639","user_id":1,"sum":1123.23,"processed_at":"2020-12-10T15:12:01+03:00"},{"order":"4180189435","user_id":1,"sum":1123.23,"processed_at":"2020-13-10T15:12:01+03:00"}]`,
			},
		},
		{
			name: "204",
			args: args{
				userID: 2,
			},
			want: want{
				code: 204,
				body: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/api/user/balance/withdrawls", nil)
			ctx := context.WithValue(r.Context(), model.UserIDKey, tt.args.userID)
			w := httptest.NewRecorder()
			got := ser.Withdrawals()
			got.ServeHTTP(w, r.WithContext(ctx))
			resp := w.Result()
			defer resp.Body.Close()
			if !reflect.DeepEqual(w.Result().StatusCode, tt.want.code) {
				t.Errorf("Service.Withdrawals() %v = %v, want %v", tt.name, w.Result().StatusCode, tt.want.code)
			}
			if !reflect.DeepEqual(w.Body.String(), tt.want.body) {
				t.Errorf("Service.Withdrawals() %v = %v, want %v", tt.name, w.Body.String(), tt.want.body)
			}
		})
	}
}
