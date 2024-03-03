package usecases

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	usecasesmocks "github.com/soltanat/go-diploma-1/internal/entities/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager/mocks"
)

func TestOrderUseCase_CreateOrder(t *testing.T) {
	orderStorage := mocks.NewMockOrderStorager(gomock.NewController(t))
	userStorage := mocks.NewMockUserStorager(gomock.NewController(t))

	orderProcessor := usecasesmocks.NewMockOrderProcessorUseCase(gomock.NewController(t))

	orderUseCase, err := NewOrderUseCase(orderStorage, userStorage, orderProcessor)
	assert.NoError(t, err)

	now := time.Now()
	entities.Now = func() time.Time { return now }

	hPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	t.Run("Valid OrderNumber Creation", func(t *testing.T) {
		orderNumber := entities.OrderNumber(4561261212345467)
		userID := entities.Login("user")

		returnUser := &entities.User{
			Login:    userID,
			Password: hPassword,
			Balance: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
		}

		userStorage.EXPECT().Get(gomock.Any(), nil, userID).Return(returnUser, nil)
		orderStorage.EXPECT().Get(gomock.Any(), nil, orderNumber).Return(nil, entities.NotFoundError{})

		order := &entities.Order{
			Number: orderNumber,
			Status: entities.OrderStatusNEW,
			Accrual: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
			UploadedAt: now,
			UserID:     userID,
		}

		orderStorage.EXPECT().Save(gomock.Any(), nil, order).Return(nil)

		orderProcessor.EXPECT().ProcessOrder(gomock.Any(), orderNumber).Return(nil)

		err = orderUseCase.CreateOrder(context.Background(), orderNumber, userID)
		assert.NoError(t, err)
	})

	t.Run("OrderNumber Validation Error", func(t *testing.T) {
		orderNumber := entities.OrderNumber(123)
		userID := entities.Login("user")

		err = orderUseCase.CreateOrder(context.Background(), orderNumber, userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.InvalidOrderNumberError{})
	})

	t.Run("Login Validation Error", func(t *testing.T) {
		orderNumber := entities.OrderNumber(4561261212345467)
		userID := entities.Login("")

		err = orderUseCase.CreateOrder(context.Background(), orderNumber, userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("invalid login: ")})
	})

	t.Run("OrderNumber Creation Error Not found user", func(t *testing.T) {
		orderNumber := entities.OrderNumber(4561261212345467)
		userID := entities.Login("user")

		userStorage.EXPECT().Get(gomock.Any(), nil, userID).Return(nil, entities.NotFoundError{})

		err = orderUseCase.CreateOrder(context.Background(), orderNumber, userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.NotFoundError{})
	})

	t.Run("OrderNumber Creation Error OrderNumber already exists", func(t *testing.T) {
		orderNumber := entities.OrderNumber(4561261212345467)
		userID := entities.Login("user")

		returnUser := &entities.User{
			Login:    userID,
			Password: hPassword,
			Balance: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
		}
		userStorage.EXPECT().Get(gomock.Any(), nil, userID).Return(returnUser, nil)
		returnOrder := &entities.Order{
			Number: orderNumber,
			Status: entities.OrderStatusNEW,
			Accrual: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
			UploadedAt: now,
			UserID:     userID,
		}
		orderStorage.EXPECT().Get(gomock.Any(), nil, orderNumber).Return(returnOrder, nil)

		err = orderUseCase.CreateOrder(context.Background(), orderNumber, userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ExistOrderError{})
	})

	t.Run("OrderNumber Creation Error Already created another user", func(t *testing.T) {
		orderNumber := entities.OrderNumber(4561261212345467)
		userID := entities.Login("user")

		returnUser := &entities.User{
			Login:    userID,
			Password: hPassword,
			Balance: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
		}
		userStorage.EXPECT().Get(gomock.Any(), nil, userID).Return(returnUser, nil)
		returnOrder := &entities.Order{
			Number: orderNumber,
			Status: entities.OrderStatusNEW,
			Accrual: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
			UploadedAt: now,
			UserID:     entities.Login("user2"),
		}
		orderStorage.EXPECT().Get(gomock.Any(), nil, orderNumber).Return(returnOrder, nil)

		err = orderUseCase.CreateOrder(context.Background(), orderNumber, userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.OrderIsCreatedByAnotherUserError{})
	})
}

func TestOrderUseCase_ListOrdersByUserID(t *testing.T) {
	orderStorage := mocks.NewMockOrderStorager(gomock.NewController(t))
	userStorage := mocks.NewMockUserStorager(gomock.NewController(t))

	orderProcessor := usecasesmocks.NewMockOrderProcessorUseCase(gomock.NewController(t))

	orderUseCase, err := NewOrderUseCase(orderStorage, userStorage, orderProcessor)
	assert.NoError(t, err)

	now := time.Now()
	entities.Now = func() time.Time { return now }

	t.Run("Valid Get List Orders", func(t *testing.T) {
		userID := entities.Login("user")

		userStorage.EXPECT().Get(gomock.Any(), nil, userID).Return(&entities.User{}, nil)

		returnOrders := []entities.Order{
			{
				Number: entities.OrderNumber(1),
				Status: entities.OrderStatusNEW,
				Accrual: entities.Currency{
					Whole:   0,
					Decimal: 0,
				},
				UploadedAt: now,
				UserID:     userID,
			},
			{
				Number: entities.OrderNumber(2),
				Status: entities.OrderStatusNEW,
				Accrual: entities.Currency{
					Whole:   0,
					Decimal: 0,
				},
				UploadedAt: now,
				UserID:     userID,
			},
		}

		orderStorage.EXPECT().List(gomock.Any(), nil, &userID, nil).Return(returnOrders, nil)

		result, err := orderUseCase.ListOrdersByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, returnOrders, result)
	})

	t.Run("Validation User Error", func(t *testing.T) {
		userID := entities.Login("")

		_, err := orderUseCase.ListOrdersByUserID(context.Background(), userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.ValidationError{Err: fmt.Errorf("invalid login: ")})
	})

	t.Run("User not found error", func(t *testing.T) {
		userID := entities.Login("user")

		userStorage.EXPECT().Get(gomock.Any(), nil, userID).Return(nil, entities.NotFoundError{})

		_, err := orderUseCase.ListOrdersByUserID(context.Background(), userID)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &entities.NotFoundError{})
	})

}
