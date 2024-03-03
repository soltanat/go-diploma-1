package main

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/soltanat/go-diploma-1/internal/clients/accrual"
	"github.com/soltanat/go-diploma-1/internal/db"
	"github.com/soltanat/go-diploma-1/internal/http"
	"github.com/soltanat/go-diploma-1/internal/http/api"
	"github.com/soltanat/go-diploma-1/internal/logger"
	"github.com/soltanat/go-diploma-1/internal/storage/external"
	"github.com/soltanat/go-diploma-1/internal/storage/limit"
	"github.com/soltanat/go-diploma-1/internal/storage/postgres"
	"github.com/soltanat/go-diploma-1/internal/storage/retry"
	"github.com/soltanat/go-diploma-1/internal/usecases"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	l := logger.Get()

	parseFlags()

	conn, err := db.New(ctx, flagDBAddr)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to connect to database")
	}
	defer conn.Close()

	err = db.ApplyMigrations(flagDBAddr)
	if err != nil {
		l.Fatal().Err(err).Msg("unable to apply migrations")
	}

	pool, err := db.New(ctx, flagDBAddr)
	if err != nil {
		l.Fatal().Err(err)
	}
	defer pool.Close()

	userStorage := postgres.NewUserStorage(pool)
	orderStorage := postgres.NewOrderStorage(pool)
	withdrawalStorage := postgres.NewWithdrawalStorage(pool)

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Second

	userStorage = retry.NewUserStorage(userStorage, b)
	orderStorage = retry.NewOrderStorage(orderStorage, b)
	withdrawalStorage = retry.NewWithdrawalStorage(withdrawalStorage, b)

	client, err := accrual.NewClientWithResponses(flagAccrualAddr)
	if err != nil {
		l.Fatal().Err(err)
	}
	accrualStorage, err := external.NewAccrualStorage(client)
	if err != nil {
		l.Fatal().Err(err)
	}

	accrualStorage = limit.NewLimitAccrualStorage(accrualStorage, flagAccrualRateLimit)

	accrualStorage = retry.NewAccrualStorage(accrualStorage, b)

	userUseCase, err := usecases.NewUserUseCase(userStorage)
	if err != nil {
		l.Fatal().Err(err)
	}

	orderProcessor, err := usecases.NewOrderProcessor(userStorage, orderStorage, accrualStorage)
	if err != nil {
		l.Fatal().Err(err)
	}

	orderUseCase, err := usecases.NewOrderUseCase(orderStorage, userStorage, orderProcessor)
	if err != nil {
		l.Fatal().Err(err)
	}

	withdrawalUseCase, err := usecases.NewWithdrawUseCase(withdrawalStorage, userStorage)
	if err != nil {
		l.Fatal().Err(err)
	}

	go func() {
		for {
			err := orderProcessor.Produce(ctx)
			if err != nil {
				l.Fatal().Err(err)
			}
		}
	}()

	go func() {
		orderProcessor.Run(ctx)
	}()

	handler := http.NewServerInterfaceWrapper(
		userUseCase,
		orderUseCase,
		withdrawalUseCase,
		http.NewJWTProvider(flagKey, jwt.SigningMethodHS256),
	)
	strictHandler := api.NewStrictHandler(handler, nil)

	e := echo.New()
	e.HideBanner = true
	api.RegisterHandlers(e, strictHandler)

	go func() {
		err := e.Start(flagAddr)
		if err != nil {
			l.Error().Err(err)
		}
	}()

	gracefulShutdown()

	orderProcessor.Stop()

	err = e.Close()
	if err != nil {
		l.Error().Err(err).Msg("unable to close server")
	}

}

func gracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
}
