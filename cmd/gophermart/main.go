package main

import (
	"fmt"
	"gofemart/internal/app"
	"gofemart/internal/app/service"
	"gofemart/internal/jwt"
	"gofemart/internal/storage"
	"net/http"
)

func run() error {
	parseFlags()
	db := storage.InitDB(addrDB)
	if db == nil {
		return fmt.Errorf("failed to connect to database")
	}
	store := storage.NewDatabase(db)
	if err := storage.MigrateTables(store); err != nil {
		return err
	}
	jwtManager, err := jwt.NewJWTManager("sfjvpasasdf", "30m", "24h")
	if err != nil {
		return err
	}
	s := service.NewService(store, *jwtManager)
	mux := http.NewServeMux()
	app.AddRoute(mux, s)
	return http.ListenAndServe(addrServer, mux)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		return
	}
}
