package retry

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Accrual struct {
	storage storager.AccrualOrderStorager
	b       *backoff.ExponentialBackOff
}

func NewAccrualStorage(storage storager.AccrualOrderStorager, b *backoff.ExponentialBackOff) *Accrual {
	return &Accrual{
		storage: storage,
		b:       b,
	}
}

func (a Accrual) Get(ctx context.Context, number entities.OrderNumber) (*entities.AccrualOrder, error) {
	var order *entities.AccrualOrder

	err := backoff.Retry(func() error {
		var err error
		order, err = a.storage.Get(ctx, number)
		if err != nil {
			if errors.Is(err, entities.NotFoundError{}) {
				return backoff.Permanent(err)
			}
			return err
		}
		return nil
	}, a.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return nil, err.(*backoff.PermanentError).Err
	}

	return order, err
}
