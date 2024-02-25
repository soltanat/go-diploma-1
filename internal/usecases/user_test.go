package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager/mocks"
)

func TestUserUseCase_Register(t *testing.T) {
	mockStorage := mocks.NewMockUserStorager(gomock.NewController(t))

	userUseCase, err := NewUserUseCase(mockStorage)
	assert.NoError(t, err)

	t.Run("Valid User Creation", func(t *testing.T) {
		user := &entities.User{
			Login:    "login",
			Password: "password",
		}

		mockStorage.EXPECT().Get(gomock.Any(), user.Login, nil).Return(nil, entities.NotFoundError{})
		mockStorage.EXPECT().Save(gomock.Any(), user).Return(nil)

		err = userUseCase.Register(context.Background(), user.Login, user.Password)
		assert.NoError(t, err)
	})

	t.Run("User Validation Error", func(t *testing.T) {
		err = userUseCase.Register(context.Background(), "", "")
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("invalid login: ")})
	})

	t.Run("Existing User Error", func(t *testing.T) {
		user := &entities.User{
			Login:    "login",
			Password: "password",
		}

		mockStorage.EXPECT().Get(gomock.Any(), user.Login, nil).Return(nil, entities.ExistUserError{})

		err := userUseCase.Register(context.Background(), user.Login, user.Password)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ExistUserError{})
	})

	t.Run("Storage SaveTx Error", func(t *testing.T) {
		user := &entities.User{
			Login:    "login",
			Password: "password",
		}

		mockStorage.EXPECT().Get(gomock.Any(), user.Login, nil).Return(nil, entities.NotFoundError{})
		mockStorage.EXPECT().Save(gomock.Any(), user).Return(entities.StorageError{})

		err := userUseCase.Register(context.Background(), user.Login, user.Password)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.StorageError{})
	})
}

func TestUserUseCase_Authenticate(t *testing.T) {
	mockStorage := mocks.NewMockUserStorager(gomock.NewController(t))

	userUseCase, err := NewUserUseCase(mockStorage)
	assert.NoError(t, err)

	t.Run("Valid User Authentication", func(t *testing.T) {
		user := &entities.User{
			Login:    "login",
			Password: "password",
		}

		returnUser := &entities.User{
			Login:    "login",
			Password: "password",
			Balance: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
		}
		mockStorage.EXPECT().Get(gomock.Any(), user.Login, &user.Password).Return(returnUser, nil)

		result, err := userUseCase.Authenticate(context.Background(), user.Login, user.Password)
		assert.NoError(t, err)
		assert.Equal(t, returnUser, result)
	})

	t.Run("User Login Error", func(t *testing.T) {
		user := &entities.User{
			Password: "password",
		}
		result, err := userUseCase.Authenticate(context.Background(), user.Login, user.Password)
		assert.Nil(t, result)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("invalid login: ")})
	})

	t.Run("User Password Error", func(t *testing.T) {
		user := &entities.User{
			Login: "login",
		}
		result, err := userUseCase.Authenticate(context.Background(), user.Login, user.Password)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("password is empty")})
	})

	t.Run("User Not Found Error", func(t *testing.T) {
		user := &entities.User{
			Login:    "login",
			Password: "password",
		}
		mockStorage.EXPECT().Get(gomock.Any(), user.Login, &user.Password).Return(nil, entities.NotFoundError{})
		result, err := userUseCase.Authenticate(context.Background(), user.Login, user.Password)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.NotFoundError{})
	})

	t.Run("Storage Get Error", func(t *testing.T) {
		user := &entities.User{
			Login:    "login",
			Password: "password",
		}
		mockStorage.EXPECT().Get(gomock.Any(), user.Login, &user.Password).Return(nil, entities.StorageError{})
		result, err := userUseCase.Authenticate(context.Background(), user.Login, user.Password)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.StorageError{})
		assert.Nil(t, result)
	})
}
