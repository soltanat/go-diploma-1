package usecases

import (
	"context"
	"fmt"

	"github.com/soltanat/go-diploma-1/internal/model"
	"github.com/soltanat/go-diploma-1/internal/model/entities"
)

type WithdrawUseCase struct {
	withdrawalStorager model.WithdrawalStorager
	userStorager       model.UserStorager
}

func NewWithdrawUseCase(withdrawalStorager model.WithdrawalStorager, userStorager model.UserStorager) (*WithdrawUseCase, error) {
	if withdrawalStorager == nil {
		return nil, fmt.Errorf("withdrawalStorager is nil")
	}
	if userStorager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	return &WithdrawUseCase{
		withdrawalStorager: withdrawalStorager,
		userStorager:       userStorager,
	}, nil
}

func (u *WithdrawUseCase) ListWithdrawals(
	ctx context.Context, userID entities.Login) ([]entities.Withdrawal, error,
) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}
	return u.withdrawalStorager.List(ctx, userID)
}

func (u *WithdrawUseCase) Withdraw(
	ctx context.Context, userID entities.Login, orderNumber entities.OrderNumber, sum entities.Currency,
) error {
	if err := userID.Validate(); err != nil {
		return err
	}
	if err := orderNumber.Validate(); err != nil {
		return err
	}
	if err := sum.Validate(); err != nil {
		return err
	}

	tx := u.userStorager.Tx(ctx)
	user, err := u.userStorager.GetTx(ctx, tx, userID, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := user.Balance.Sub(&sum); err != nil {
		tx.Rollback()
		return err
	}

	err = u.userStorager.UpdateTx(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	withdrawal := entities.NewWithdrawal(orderNumber, sum, userID)
	err = u.withdrawalStorager.SaveTx(ctx, tx, withdrawal)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
