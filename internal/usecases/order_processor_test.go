package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager/mocks"
)

func TestOrderProcessor_ProcessOrder(t *testing.T) {
	orderStorage := mocks.NewMockOrderStorager(gomock.NewController(t))
	userStorage := mocks.NewMockUserStorager(gomock.NewController(t))
	accrualStorage := mocks.NewMockAccrualOrderStorager(gomock.NewController(t))

	orderUseCase, err := NewOrderProcessor(userStorage, orderStorage, accrualStorage)
	assert.NoError(t, err)

	now := time.Now()
	entities.Now = func() time.Time { return now }

	t.Run("Valid OrderNumber Process (Accrual Processed)", func(t *testing.T) {
		orderNumber := entities.OrderNumber(1)

		tx := mocks.NewMockTx(gomock.NewController(t))

		orderStorage.EXPECT().Tx(gomock.Any()).Return(tx)

		returnOrder := entities.Order{
			Number: orderNumber,
			Status: entities.OrderStatusNEW,
			Accrual: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
			UploadedAt: now,
			UserID:     entities.Login("user"),
		}
		orderStorage.EXPECT().Get(gomock.Any(), tx, orderNumber).Return(&returnOrder, nil)

		returnUser := entities.User{
			Login:    entities.Login("user"),
			Password: "password",
			Balance: entities.Currency{
				Whole:   0,
				Decimal: 0,
			},
		}
		userStorage.EXPECT().Get(gomock.Any(), tx, returnOrder.UserID).Return(&returnUser, nil)

		returnAccrual := entities.AccrualOrder{
			Number: orderNumber,
			Status: entities.AccrualOrderStatusPROCESSED,
			Accrual: &entities.Currency{
				Whole:   20,
				Decimal: 20,
			},
		}
		accrualStorage.EXPECT().Get(gomock.Any(), orderNumber).Return(&returnAccrual, nil)

		updatedUser := returnUser
		updatedUser.Balance = entities.Currency{
			Whole:   20,
			Decimal: 20,
		}
		userStorage.EXPECT().Update(gomock.Any(), tx, &updatedUser).Return(nil)

		updatedOrder := returnOrder
		updatedOrder.Status = entities.OrderStatusPROCESSED
		updatedOrder.Accrual = entities.Currency{
			Whole:   20,
			Decimal: 20,
		}
		orderStorage.EXPECT().Update(gomock.Any(), tx, &updatedOrder).Return(nil)

		tx.EXPECT().Commit(gomock.Any()).Return(nil)

		err := orderUseCase.ProcessOrder(context.Background(), orderNumber)
		assert.NoError(t, err)
	})
}
