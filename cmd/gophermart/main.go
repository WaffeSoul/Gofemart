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
	acc := accrual.NewAccrual(conf, &store)
	defer acc.Finish()
	srv := service.NewService(store, acc)
	mux := http.NewServeMux()
	app.AddRoute(mux, srv)
	res := http.ListenAndServe(conf.Server, mux)
	return res
}

func main() {
	res := run()
	logger.Info("finish run server", zap.Error(res))
}
