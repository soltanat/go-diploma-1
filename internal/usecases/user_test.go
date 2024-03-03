package usecases

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/soltanat/go-diploma-1/internal/entities"
	entitiesMocks "github.com/soltanat/go-diploma-1/internal/entities/mocks"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager/mocks"
)

func TestUserUseCase_Register(t *testing.T) {
	mockStorage := mocks.NewMockUserStorager(gomock.NewController(t))
	mockHasher := entitiesMocks.NewMockPasswordHasher(gomock.NewController(t))

	userUseCase, err := NewUserUseCase(mockStorage, mockHasher)
	assert.NoError(t, err)

	login := entities.Login("login")
	password := "password"

	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	t.Run("Valid User Creation", func(t *testing.T) {
		mockHasher.EXPECT().Hash(gomock.Any()).Return(hPassword, nil)
		mockStorage.EXPECT().Get(gomock.Any(), nil, login).Return(nil, entities.NotFoundError{})
		saveUser := &entities.User{
			Login:    entities.Login(login),
			Password: hPassword,
		}
		mockStorage.EXPECT().Save(gomock.Any(), nil, saveUser).Return(nil)

		err = userUseCase.Register(context.Background(), entities.Login(login), password)
		assert.NoError(t, err)
	})

	t.Run("User Validation Error", func(t *testing.T) {
		mockHasher.EXPECT().Hash(gomock.Any()).Return(hPassword, nil)
		err = userUseCase.Register(context.Background(), "", password)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("invalid login: ")})
	})

	t.Run("Existing User Error", func(t *testing.T) {
		mockHasher.EXPECT().Hash(gomock.Any()).Return(hPassword, nil)
		mockStorage.EXPECT().Get(gomock.Any(), nil, login).Return(nil, entities.ExistUserError{})

		err := userUseCase.Register(context.Background(), entities.Login(login), password)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ExistUserError{})
	})

	t.Run("Storage Save Error", func(t *testing.T) {
		mockHasher.EXPECT().Hash(gomock.Any()).Return(hPassword, nil)
		mockStorage.EXPECT().Get(gomock.Any(), nil, login).Return(nil, entities.NotFoundError{})
		saveU1ser := &entities.User{
			Login:    entities.Login(login),
			Password: hPassword,
		}
		mockStorage.EXPECT().Save(gomock.Any(), nil, saveU1ser).Return(entities.StorageError{})

		err := userUseCase.Register(context.Background(), entities.Login(login), password)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.StorageError{})
	})
}

func TestUserUseCase_Authenticate(t *testing.T) {
	mockStorage := mocks.NewMockUserStorager(gomock.NewController(t))
	mockHasher := entitiesMocks.NewMockPasswordHasher(gomock.NewController(t))

	userUseCase, err := NewUserUseCase(mockStorage, mockHasher)
	assert.NoError(t, err)

	login := entities.Login("login")
	password := "password"
	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	t.Run("Valid User Authentication", func(t *testing.T) {
		returnUser := &entities.User{
			Login:    "login",
			Password: hPassword,
			Balance: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
		}
		mockStorage.EXPECT().Get(gomock.Any(), nil, login).Return(returnUser, nil)
		mockHasher.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(true)

		result, err := userUseCase.Authenticate(context.Background(), login, password)
		assert.NoError(t, err)
		assert.Equal(t, returnUser, result)
	})

	t.Run("User Login Error", func(t *testing.T) {
		result, err := userUseCase.Authenticate(context.Background(), "", "password")
		assert.Nil(t, result)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("invalid login: ")})
	})

	t.Run("User Password Error", func(t *testing.T) {
		result, err := userUseCase.Authenticate(context.Background(), "login", "")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("password is empty")})
	})

	t.Run("User Not Found Error", func(t *testing.T) {
		mockStorage.EXPECT().Get(gomock.Any(), nil, login).Return(nil, entities.NotFoundError{})
		result, err := userUseCase.Authenticate(context.Background(), "login", "password")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.NotFoundError{})
	})

	t.Run("Storage Get Error", func(t *testing.T) {
		mockStorage.EXPECT().Get(gomock.Any(), nil, login).Return(nil, entities.StorageError{})
		result, err := userUseCase.Authenticate(context.Background(), "login", "password")
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.StorageError{})
		assert.Nil(t, result)
	})
}
