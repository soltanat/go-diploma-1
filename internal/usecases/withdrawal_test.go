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

func TestWithdrawUseCase_Withdraw(t *testing.T) {
	withdrawalStorage := mocks.NewMockWithdrawalStorager(gomock.NewController(t))
	userStorage := mocks.NewMockUserStorager(gomock.NewController(t))

	withdrawalUseCase, err := NewWithdrawUseCase(withdrawalStorage, userStorage)
	assert.NoError(t, err)

	now := time.Now()
	entities.Now = func() time.Time { return now }

	t.Run("Valid Withdrawal", func(t *testing.T) {
		userID := entities.Login("user")
		orderNumber := entities.OrderNumber(1)
		sum := entities.Currency{
			Whole:   20,
			Decimal: 20,
		}

		tx := mocks.NewMockTx(gomock.NewController(t))
		userStorage.EXPECT().Tx(gomock.Any()).Return(tx)

		returnUser := entities.User{
			Login:    userID,
			Password: "password",
			Balance: entities.Currency{
				Whole:   100,
				Decimal: 50,
			},
		}
		userStorage.EXPECT().Get(gomock.Any(), tx, userID).Return(&returnUser, nil)

		updatedUser := returnUser
		updatedUser.Balance = entities.Currency{
			Whole:   80,
			Decimal: 30,
		}
		userStorage.EXPECT().Update(gomock.Any(), tx, &updatedUser).Return(nil)

		newWithdrawal := entities.Withdrawal{
			OrderNumber: orderNumber,
			UserID:      userID,
			Sum:         sum,
			ProcessedAt: now,
		}
		withdrawalStorage.EXPECT().Save(gomock.Any(), tx, &newWithdrawal).Return(nil)

		tx.EXPECT().Commit(gomock.Any()).Return(nil)

		err = withdrawalUseCase.Withdraw(context.Background(), userID, orderNumber, sum)
		assert.NoError(t, err)
	})

}
