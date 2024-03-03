package retry

import (
	"context"

	"github.com/soltanat/go-diploma-1/internal/backoff"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Withdrawal struct {
	storage storager.WithdrawalStorager
}

func NewWithdrawalStorage(storage storager.WithdrawalStorager) storager.WithdrawalStorager {
	return &Withdrawal{
		storage: storage,
	}
}

func (w Withdrawal) Save(ctx context.Context, tx storager.Tx, withdraw *entities.Withdrawal) (err error) {
	err = backoff.Backoff(func() error {
		return w.storage.Save(ctx, tx, withdraw)
	}, "WithdrawalStorage.Save", entities.ExistWithdrawalError{})
	return err
}

func (w Withdrawal) List(ctx context.Context, tx storager.Tx, userID entities.Login) (ww []entities.Withdrawal, err error) {
	err = backoff.Backoff(func() error {
		ww, err = w.storage.List(ctx, tx, userID)
		return err
	}, "WithdrawalStorage.List")
	return ww, err
}

func (w Withdrawal) Count(ctx context.Context, tx storager.Tx, userID entities.Login) (count int, err error) {
	err = backoff.Backoff(func() error {
		count, err = w.storage.Count(ctx, tx, userID)
		return err
	}, "WithdrawalStorage.Count")
	return count, err
}

func (w Withdrawal) Tx(ctx context.Context) storager.Tx {
	return w.storage.Tx(ctx)
}
