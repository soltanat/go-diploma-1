package retry

import (
	"context"

	"github.com/soltanat/go-diploma-1/internal/backoff"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type User struct {
	storage storager.UserStorager
}

func NewUserStorage(storage storager.UserStorager) storager.UserStorager {
	return &User{
		storage: storage,
	}
}

func (s *User) Save(ctx context.Context, tx storager.Tx, user *entities.User) (err error) {
	err = backoff.Backoff(func() error {
		err = s.storage.Save(ctx, tx, user)
		return nil
	}, "UserStorage.Save", entities.ExistUserError{})
	return err
}

func (s *User) Get(ctx context.Context, tx storager.Tx, login entities.Login) (user *entities.User, err error) {
	err = backoff.Backoff(func() error {
		user, err = s.storage.Get(ctx, tx, login)
		return err
	}, "UserStorage.Get", entities.NotFoundError{})
	return user, err
}

func (s *User) Update(ctx context.Context, tx storager.Tx, user *entities.User) (err error) {
	err = backoff.Backoff(func() error {
		return s.storage.Update(ctx, tx, user)
	}, "UserStorage.Update", entities.NotFoundError{})
	return err
}

func (s *User) Tx(ctx context.Context) storager.Tx {
	return s.storage.Tx(ctx)
}
