package storager

import (
	"context"
	"github.com/soltanat/go-diploma-1/internal/entities"
)

type Tx interface {
	Commit()
	Rollback()
	GetTx() interface{}
}

type OrderStorager interface {
	Save(ctx context.Context, order *entities.Order) error
	Get(ctx context.Context, number entities.OrderNumber) (*entities.Order, error)
	GetTx(ctx context.Context, tx Tx, number entities.OrderNumber) (*entities.Order, error)
	List(ctx context.Context, userID *entities.Login) ([]entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	UpdateTx(ctx context.Context, tx Tx, order *entities.Order) error
	Tx(ctx context.Context) Tx
}

type UserStorager interface {
	Save(ctx context.Context, user *entities.User) error
	Get(ctx context.Context, login entities.Login, password *string) (*entities.User, error)
	GetTx(ctx context.Context, tx Tx, login entities.Login, password *string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	UpdateTx(ctx context.Context, tx Tx, user *entities.User) error
	Tx(ctx context.Context) Tx
}

type WithdrawalStorager interface {
	SaveTx(ctx context.Context, tx Tx, withdraw *entities.Withdrawal) error
	List(ctx context.Context, userID entities.Login) ([]entities.Withdrawal, error)
	Tx(ctx context.Context) Tx
}

type AccrualOrderStorager interface {
	Get(ctx context.Context, number entities.OrderNumber) (*entities.AccrualOrder, error)
}
