package retry

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Withdrawal struct {
	storage storager.WithdrawalStorager
	b       *backoff.ExponentialBackOff
}

func NewWithdrawalStorage(storage storager.WithdrawalStorager, b *backoff.ExponentialBackOff) storager.WithdrawalStorager {
	return &Withdrawal{
		storage: storage,
		b:       b,
	}
}

func (w Withdrawal) Save(ctx context.Context, tx storager.Tx, withdraw *entities.Withdrawal) error {
	err := backoff.Retry(func() error {
		err := w.storage.Save(ctx, tx, withdraw)
		if err != nil {
			if errors.Is(err, entities.ExistWithdrawalError{}) {
				return backoff.Permanent(err)
			}
			return err
		}
		return nil
	}, w.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return err.(*backoff.PermanentError).Err
	}

	return err
}

func (w Withdrawal) List(ctx context.Context, tx storager.Tx, userID entities.Login) ([]entities.Withdrawal, error) {
	var ww []entities.Withdrawal

	err := backoff.Retry(func() error {
		var err error
		ww, err = w.storage.List(ctx, tx, userID)
		if err != nil {
			return err
		}
		return nil
	}, w.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return nil, err.(*backoff.PermanentError).Err
	}

	return ww, err
}

func (w Withdrawal) Count(ctx context.Context, tx storager.Tx, userID entities.Login) (int, error) {
	count := 0
	err := backoff.Retry(func() error {
		var err error
		count, err = w.storage.Count(ctx, tx, userID)
		if err != nil {
			return err
		}
		return nil
	}, w.b)

	if err != nil && err.(*backoff.PermanentError).Err != nil {
		return 0, err.(*backoff.PermanentError).Err
	}

	return count, err
}

func (w Withdrawal) Tx(ctx context.Context) storager.Tx {
	return w.storage.Tx(ctx)
}
