package retry

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Order struct {
	storage storager.OrderStorager
	b       *backoff.ExponentialBackOff
}

func NewOrderStorage(storage storager.OrderStorager, b *backoff.ExponentialBackOff) storager.OrderStorager {
	return &Order{
		storage: storage,
		b:       b,
	}
}

func (s *Order) Save(ctx context.Context, tx storager.Tx, order *entities.Order) error {
	err := backoff.Retry(func() error {
		err := s.storage.Save(ctx, tx, order)
		if err != nil {
			if errors.Is(err, entities.ExistOrderError{}) {
				return backoff.Permanent(err)
			}
			return err
		}
		return nil
	}, s.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return err.(*backoff.PermanentError).Err
	}

	return err
}

func (s *Order) Get(ctx context.Context, tx storager.Tx, number entities.OrderNumber) (*entities.Order, error) {
	order := &entities.Order{}
	err := backoff.Retry(func() error {
		var err error
		order, err = s.storage.Get(ctx, tx, number)
		if err != nil {
			if errors.Is(err, entities.NotFoundError{}) {
				return backoff.Permanent(err)
			}
			return err
		}
		return nil
	}, s.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return nil, err.(*backoff.PermanentError).Err
	}

	return order, err
}

func (s *Order) List(ctx context.Context, tx storager.Tx, userID *entities.Login, status *[]entities.OrderStatus) ([]entities.Order, error) {
	var orders []entities.Order

	err := backoff.Retry(func() error {
		var err error
		orders, err = s.storage.List(ctx, tx, userID, status)
		if err != nil {
			return err
		}
		return nil
	}, s.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return nil, err.(*backoff.PermanentError).Err
	}

	return orders, err
}

func (s *Order) Update(ctx context.Context, tx storager.Tx, order *entities.Order) error {
	err := backoff.Retry(func() error {
		var err error
		err = s.storage.Update(ctx, tx, order)
		if err != nil {
			if errors.Is(err, entities.NotFoundError{}) {
				return backoff.Permanent(err)
			}
			return err
		}
		return nil
	}, s.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return err.(*backoff.PermanentError).Err
	}

	return err
}

func (s *Order) Tx(ctx context.Context) storager.Tx {
	return s.storage.Tx(ctx)
}
