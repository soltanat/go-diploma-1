package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/soltanat/go-diploma-1/internal/logger"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type UserUseCase struct {
	storager storager.UserStorager
	hasher   entities.PasswordHasher
}

func NewUserUseCase(storager storager.UserStorager, hasher entities.PasswordHasher) (*UserUseCase, error) {
	if storager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	if hasher == nil {
		return nil, fmt.Errorf("passwordHasher is nil")
	}
	return &UserUseCase{
		storager: storager,
		hasher:   hasher,
	}, nil
}

func (u *UserUseCase) Register(ctx context.Context, login entities.Login, password string) error {
	if password == "" {
		return entities.ValidationError{Err: fmt.Errorf("password is empty")}
	}

	hashPassword, err := u.hasher.Hash([]byte(password))
	if err != nil {
		return err
	}

	user := entities.NewUser(login, hashPassword)
	if err := user.Validate(); err != nil {
		return err
	}

	if _, err := u.storager.Get(ctx, nil, user.Login); err == nil {
		return entities.ExistUserError{}
	} else if !errors.Is(err, entities.NotFoundError{}) {
		return err
	}

	return u.storager.Save(ctx, nil, user)
}

func (u *UserUseCase) Authenticate(ctx context.Context, login entities.Login, password string) (*entities.User, error) {
	if err := login.Validate(); err != nil {
		return nil, err
	}
	if password == "" {
		return nil, entities.ValidationError{Err: fmt.Errorf("password is empty")}
	}

	user, err := u.storager.Get(ctx, nil, login)
	if err != nil {
		return nil, err
	}

	if !u.hasher.Compare(user.Password, []byte(password)) {
		return nil, entities.InvalidPasswordError{}
	}

	return user, nil
}

func (u *UserUseCase) GetUser(ctx context.Context, login entities.Login) (*entities.User, error) {
	l := logger.Get()
	if err := login.Validate(); err != nil {
		return nil, err
	}

	user, err := u.storager.Get(ctx, nil, login)
	if err != nil {
		return nil, err
	}

	l.Debug().Str("usecase", "GetUser").Msgf("found user %s balance %v", user.Login, user.Balance)

	return user, nil
}
