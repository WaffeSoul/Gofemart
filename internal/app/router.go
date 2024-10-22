package app

import (
	"net/http"

	"gofemart/internal/app/middleware"
	"gofemart/internal/app/service"
)

func AddRoute(mux *http.ServeMux, s *service.Service) {
	mux.Handle("POST /api/user/register", middleware.Middleware(s.SignUp()))
	mux.Handle("POST /api/user/login", middleware.Middleware(s.SignIn()))
	mux.Handle("POST /api/user/orders", middleware.MiddlewareWithJWT(&s.JwtManager, s.SetOrder()))
	mux.Handle("GET /api/user/orders", middleware.MiddlewareWithJWT(&s.JwtManager, s.GetOrders()))
	mux.Handle("GET /api/user/balance", middleware.MiddlewareWithJWT(&s.JwtManager, s.GetBalance()))
	mux.Handle("POST /api/user/balance/withdraw", middleware.MiddlewareWithJWT(&s.JwtManager, s.Withdraw()))
	mux.Handle("GET /api/user/withdrawals", middleware.MiddlewareWithJWT(&s.JwtManager, s.Withdrawals()))
}
