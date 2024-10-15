package main

import (
	"gofemart/internal/accrual"
	"gofemart/internal/app"
	"gofemart/internal/app/service"
	"gofemart/internal/config"
	"gofemart/internal/logger"
	"gofemart/internal/storage"
	"net/http"

	"go.uber.org/zap"
)

func run() error {
	logger.InitLogger(false)
	conf := config.NewConfig()
	store := storage.NewStore(conf)
	acc := accrual.NewAccrual(conf)
	srv := service.NewService(store, acc)
	mux := http.NewServeMux()
	app.AddRoute(mux, srv)
	return http.ListenAndServe(conf.Server, mux)
}

func main() {
	if err := run(); err != nil {
		logger.Error("failed to run server", zap.Error(err))
		return
	}
}
