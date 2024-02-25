package usecases

import (
	"context"
	"fmt"

	"github.com/soltanat/go-diploma-1/internal/model"
	"github.com/soltanat/go-diploma-1/internal/model/entities"
)

type OrderProcessor struct {
	userStorager    model.UserStorager
	orderStorager   model.OrderStorager
	accrualStorager model.AccrualOrderStorager
}

func NewOrderProcessor(userStorager model.UserStorager, orderStorager model.OrderStorager, accrualStorager model.AccrualOrderStorager) (*OrderProcessor, error) {
	if userStorager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	if orderStorager == nil {
		return nil, fmt.Errorf("orderStorager is nil")
	}
	if accrualStorager == nil {
		return nil, fmt.Errorf("accrualStorager is nil")
	}

	return &OrderProcessor{
		userStorager:    userStorager,
		orderStorager:   orderStorager,
		accrualStorager: accrualStorager,
	}, nil
}

func (u *OrderProcessor) ProcessOrder(ctx context.Context, number entities.OrderNumber) error {
	if err := number.Validate(); err != nil {
		return err
	}

	tx := u.orderStorager.Tx(ctx)

	order, err := u.orderStorager.GetTx(ctx, tx, number)
	if err != nil {
		tx.Rollback()
		return err
	}

	if order.IsProcessed() {
		tx.Rollback()
		return nil
	}

	user, err := u.userStorager.GetTx(ctx, tx, order.UserID, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	accrualOrder, err := u.accrualStorager.Get(ctx, number)
	if err != nil {
		tx.Rollback()
		return err
	}

	updated := order.UpdateWithAccrualOrder(accrualOrder)

	if !updated {
		tx.Rollback()
		return nil
	}

	user.Balance.Add(&accrualOrder.Accrual)

	err = u.userStorager.UpdateTx(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = u.orderStorager.UpdateTx(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
