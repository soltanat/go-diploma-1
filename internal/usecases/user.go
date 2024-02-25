package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type UserUseCase struct {
	storager storager.UserStorager
}

func NewUserUseCase(storager storager.UserStorager) (*UserUseCase, error) {
	if storager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	return &UserUseCase{
		storager: storager,
	}, nil
}

func (u *UserUseCase) Register(ctx context.Context, login entities.Login, password string) error {
	user := entities.NewUser(login, password)
	if err := user.Validate(); err != nil {
		return err
	}

	if _, err := u.storager.Get(ctx, user.Login, nil); err == nil {
		return entities.ExistUserError{}
	} else if err != nil && !errors.Is(err, entities.NotFoundError{}) {
		return err
	}

	return u.storager.Save(ctx, user)
}

func (u *UserUseCase) Authenticate(ctx context.Context, login entities.Login, password string) (*entities.User, error) {
	if err := login.Validate(); err != nil {
		return nil, err
	}
	if password == "" {
		return nil, entities.ValidationError{Err: fmt.Errorf("password is empty")}
	}

	user, err := u.storager.Get(ctx, login, &password)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, entities.InvalidPasswordError{}
	}

	return user, nil
}
