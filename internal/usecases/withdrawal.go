package usecases

import (
	"context"
	"fmt"
	"github.com/soltanat/go-diploma-1/internal/logger"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type WithdrawUseCase struct {
	withdrawalStorager storager.WithdrawalStorager
	userStorager       storager.UserStorager
}

func NewWithdrawUseCase(withdrawalStorager storager.WithdrawalStorager, userStorager storager.UserStorager) (*WithdrawUseCase, error) {
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

func (u *WithdrawUseCase) List(
	ctx context.Context, userID entities.Login) ([]entities.Withdrawal, error,
) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}
	return u.withdrawalStorager.List(ctx, nil, userID)
}

func (u *WithdrawUseCase) Count(ctx context.Context, userID entities.Login) (int, error) {
	if err := userID.Validate(); err != nil {
		return 0, err
	}
	return u.withdrawalStorager.Count(ctx, nil, userID)
}

func (u *WithdrawUseCase) Withdraw(
	ctx context.Context, userID entities.Login, orderNumber entities.OrderNumber, amount entities.Currency,
) error {
	l := logger.Get()

	if err := userID.Validate(); err != nil {
		return err
	}
	if err := orderNumber.Validate(); err != nil {
		return err
	}
	if err := amount.Validate(); err != nil {
		return err
	}

	tx := u.userStorager.Tx(ctx)
	err := tx.Begin(ctx)
	if err != nil {
		return err
	}

	user, err := u.userStorager.Get(ctx, tx, userID)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	l.Debug().Str("usecase", "Withdraw").Msgf("found user %s balance %v", user.Login, user.Balance)

	if err := user.Balance.Sub(&amount); err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = u.userStorager.Update(ctx, tx, user)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	l.Debug().Str("usecase", "Withdraw").Msgf("updated user %s balance %v after withdraw %v", user.Login, user.Balance, amount)

	withdrawal := entities.NewWithdrawal(orderNumber, amount, userID)
	err = u.withdrawalStorager.Save(ctx, tx, withdrawal)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	l.Debug().Str("usecase", "Withdraw").Msgf("created withdrawal %v", withdrawal)

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
