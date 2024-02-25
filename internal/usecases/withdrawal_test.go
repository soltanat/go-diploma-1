package usecases

import (
	"context"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
		userStorage.EXPECT().GetTx(gomock.Any(), tx, userID, nil).Return(&returnUser, nil)

		updatedUser := returnUser
		updatedUser.Balance = entities.Currency{
			Whole:   80,
			Decimal: 30,
		}
		userStorage.EXPECT().UpdateTx(gomock.Any(), tx, &updatedUser).Return(nil)

		newWithdrawal := entities.Withdrawal{
			Order:       orderNumber,
			UserID:      userID,
			Sum:         sum,
			ProcessedAt: now,
		}
		withdrawalStorage.EXPECT().SaveTx(gomock.Any(), tx, &newWithdrawal).Return(nil)

		tx.EXPECT().Commit()

		err = withdrawalUseCase.Withdraw(context.Background(), userID, orderNumber, sum)
		assert.NoError(t, err)
	})

}
