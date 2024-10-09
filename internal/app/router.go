package app

import (
	"net/http"

	"gofemart/internal/app/middleware"
	"gofemart/internal/app/service"
)

func AddRoute(mux *http.ServeMux, s *service.Service) {
	mux.Handle("POST /api/user/register", middleware.GzipMiddleware(s.SignUp()))
	mux.Handle("POST /api/user/login", middleware.GzipMiddleware(s.SignIn()))
	mux.Handle("POST /api/user/orders", middleware.GzipMiddleware(middleware.JWTInterceptor(&s.JwtManager, s.SetOrder())))
	mux.Handle("GET /api/user/orders", middleware.GzipMiddleware(middleware.JWTInterceptor(&s.JwtManager, s.GetOrders())))
	mux.Handle("GET /api/user/balance", middleware.GzipMiddleware(middleware.JWTInterceptor(&s.JwtManager, s.GetBalance())))
	mux.Handle("POST /api/user/balance/withdraw", middleware.GzipMiddleware(middleware.JWTInterceptor(&s.JwtManager, s.Withdraw())))
	mux.Handle("GET /api/user/withdrawals", middleware.GzipMiddleware(middleware.JWTInterceptor(&s.JwtManager, s.Withdrawals())))
}
