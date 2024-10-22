package config

import (
	"flag"
	"os"
)

type Config struct {
	Server  string
	DB      string
	Accrual string
}

func NewConfig() *Config {
	var config Config
	flag.StringVar(&config.Server, "a", "", "address and port to run server")
	flag.StringVar(&config.DB, "d", "", "address and port to connect db")
	flag.StringVar(&config.Accrual, "r", "", "address and port to connect accrual")
	flag.Parse()
	if config.Server == "" {
		if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
			config.Server = envAddr
		} else {
			config.Server = "localhost:8000"
		}
	}
	if config.DB == "" {
		if envAddrDB := os.Getenv("DATABASE_URI"); envAddrDB != "" {
			config.DB = envAddrDB
		} else {
			config.DB = "postgresql://test:test@127.0.0.1:5433/test?sslmode=disable"
		}
	}
	if config.Accrual == "" {
		if envAddrAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAddrAccrual != "" {
			config.Accrual = envAddrAccrual
		} else {
			config.Accrual = "http://localhost:8080"
		}
	}
	
	return &config
}

