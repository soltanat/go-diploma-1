package storager

import (
	"context"

	"github.com/soltanat/go-diploma-1/internal/entities"
)

type Tx interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type OrderStorager interface {
	Save(ctx context.Context, tx Tx, order *entities.Order) error
	Get(ctx context.Context, tx Tx, number entities.OrderNumber) (*entities.Order, error)
	List(ctx context.Context, tx Tx, userID *entities.Login) ([]entities.Order, error)
	Update(ctx context.Context, tx Tx, order *entities.Order) error
	Tx(ctx context.Context) Tx
}

type UserStorager interface {
	Save(ctx context.Context, tx Tx, user *entities.User) error
	Get(ctx context.Context, tx Tx, login entities.Login) (*entities.User, error)
	Update(ctx context.Context, tx Tx, user *entities.User) error
	Tx(ctx context.Context) Tx
}

type WithdrawalStorager interface {
	Save(ctx context.Context, tx Tx, withdraw *entities.Withdrawal) error
	List(ctx context.Context, tx Tx, userID entities.Login) ([]entities.Withdrawal, error)
	Tx(ctx context.Context) Tx
}

type AccrualOrderStorager interface {
	Get(ctx context.Context, number entities.OrderNumber) (*entities.AccrualOrder, error)
}
