package retry

import (
	"context"

	"github.com/soltanat/go-diploma-1/internal/backoff"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Order struct {
	storage storager.OrderStorager
}

func NewOrderStorage(storage storager.OrderStorager) storager.OrderStorager {
	return &Order{
		storage: storage,
	}
}

func (s *Order) Save(ctx context.Context, tx storager.Tx, order *entities.Order) (err error) {
	err = backoff.Backoff(func() error {
		return s.storage.Save(ctx, tx, order)
	}, "OrderStorage.Save", entities.ExistOrderError{})
	return err
}

func (s *Order) Get(ctx context.Context, tx storager.Tx, number entities.OrderNumber) (order *entities.Order, err error) {
	err = backoff.Backoff(func() error {
		order, err = s.storage.Get(ctx, tx, number)
		return err
	}, "OrderStorage.Get", entities.NotFoundError{})
	return order, err
}

func (s *Order) List(ctx context.Context, tx storager.Tx, userID *entities.Login, status *[]entities.OrderStatus) (oo []entities.Order, err error) {
	err = backoff.Backoff(func() error {
		oo, err = s.storage.List(ctx, tx, userID, status)
		return err
	}, "OrderStorage.List")
	return oo, err
}

func (s *Order) Update(ctx context.Context, tx storager.Tx, order *entities.Order) (err error) {
	err = backoff.Backoff(func() error {
		err = s.storage.Update(ctx, tx, order)
		return err
	}, "OrderStorage.Update", entities.NotFoundError{})
	return err
}

func (s *Order) Tx(ctx context.Context) storager.Tx {
	return s.storage.Tx(ctx)
}
