package main

import (
	"flag"
	"gofemart/internal/accrual"
	"os"
)

var (
	addrServer string
	addrDB     string
)

func parseFlags() {
	flag.StringVar(&addrServer, "a", "localhost:8000", "address and port to run server")
	flag.StringVar(&addrDB, "d", "postgresql://test:test@127.0.0.1:5433/test?sslmode=disable", "address and port to connect db")
	flag.StringVar(&accrual.UrlAccrual, "r", "http://localhost:8080", "log level")
	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		addrServer = envAddr
	}
	if envAddrDB := os.Getenv("DATABASE_URI"); envAddrDB != "" {
		addrDB = envAddrDB
	}
	if envAddrAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAddrAccrual != "" {
		accrual.UrlAccrual = envAddrAccrual
	}

}
