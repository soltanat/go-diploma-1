package retry

import (
	"context"

	"github.com/soltanat/go-diploma-1/internal/backoff"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Accrual struct {
	storage storager.AccrualOrderStorager
}

func NewAccrualStorage(storage storager.AccrualOrderStorager) *Accrual {
	return &Accrual{
		storage: storage,
	}
}

func (a Accrual) Get(ctx context.Context, number entities.OrderNumber) (order *entities.AccrualOrder, err error) {
	err = backoff.Backoff(func() error {
		var err error
		order, err = a.storage.Get(ctx, number)
		return err
	}, "AccrualStorage.Get", entities.NotFoundError{})

	return order, err
}
