package main

import (
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/soltanat/go-diploma-1/internal/logger"
)

var flagAddr string
var flagDBAddr string
var flagAccrualAddr string
var flagAccrualRateLimit int
var flagKey string

type Config struct {
	Addr        string `env:"RUN_ADDRESS"`
	DBAddr      string `env:"DATABASE_URI"`
	AccrualAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func parseFlags() {
	l := logger.Get()

	flag.StringVar(&flagAddr, "a", "localhost:8081", "address and port metrics requests server")
	flag.StringVar(&flagDBAddr, "d", "postgres://postgres:postgres@localhost:5432?sslmode=disable", "database dsn")
	flag.StringVar(&flagAccrualAddr, "r", "http://localhost:8080", "accrual system address")
	flag.IntVar(&flagAccrualRateLimit, "l", 2, "accrual rate limit")
	flag.StringVar(&flagKey, "k", "very-secret-key", "key for signature")

	flag.Parse()

	var cfg Config
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			l.Fatal().Err(err)
		}
	}
	err := env.Parse(&cfg)
	if err != nil {
		l.Fatal().Err(err)
	}

	if cfg.Addr != "" {
		flagAddr = cfg.Addr
	}
	if cfg.DBAddr != "" {
		flagDBAddr = cfg.DBAddr
	}
	if cfg.AccrualAddr != "" {
		flagAccrualAddr = cfg.AccrualAddr
	}
}
