package retry

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type User struct {
	storage storager.UserStorager
	b       *backoff.ExponentialBackOff
}

func NewUserStorage(storage storager.UserStorager, b *backoff.ExponentialBackOff) storager.UserStorager {
	return &User{
		storage: storage,
		b:       b,
	}
}

func (s *User) Save(ctx context.Context, tx storager.Tx, user *entities.User) error {
	err := backoff.Retry(func() error {
		err := s.storage.Save(ctx, tx, user)
		if err != nil {
			if errors.Is(err, entities.ExistUserError{}) {
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

func (s *User) Get(ctx context.Context, tx storager.Tx, login entities.Login) (*entities.User, error) {
	user := &entities.User{}
	err := backoff.Retry(func() error {
		var err error
		user, err = s.storage.Get(ctx, tx, login)
		if err != nil {
			if errors.Is(err, entities.NotFoundError{}) {
				return backoff.Permanent(err)
			}
			return err
		}
		return nil
	}, s.b)

	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *User) Update(ctx context.Context, tx storager.Tx, user *entities.User) error {
	err := backoff.Retry(func() error {
		err := s.storage.Update(ctx, tx, user)
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

func (s *User) Tx(ctx context.Context) storager.Tx {
	return s.storage.Tx(ctx)
}
