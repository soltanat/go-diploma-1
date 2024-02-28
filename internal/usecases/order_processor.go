package usecases

import (
	"context"
	"fmt"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type OrderProcessor struct {
	userStorager    storager.UserStorager
	orderStorager   storager.OrderStorager
	accrualStorager storager.AccrualOrderStorager
	orders          chan entities.Order
}

func NewOrderProcessor(
	userStorager storager.UserStorager,
	orderStorager storager.OrderStorager,
	accrualStorager storager.AccrualOrderStorager,
) (*OrderProcessor, error) {
	if userStorager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	if orderStorager == nil {
		return nil, fmt.Errorf("orderStorager is nil")
	}
	if accrualStorager == nil {
		return nil, fmt.Errorf("accrualStorager is nil")
	}

	ordersCh := make(chan entities.Order)

	return &OrderProcessor{
		userStorager:    userStorager,
		orderStorager:   orderStorager,
		accrualStorager: accrualStorager,
		orders:          ordersCh,
	}, nil
}

func (u *OrderProcessor) AddOrder(order entities.Order) {
	u.orders <- order
}

func (u *OrderProcessor) Run(ctx context.Context) {
	for order := range u.orders {
		if ctx.Err() != nil {
			return
		}
		if err := u.ProcessOrder(ctx, order.Number); err != nil {
			//	TODO: что с ошибками
		}
	}
}

func (u *OrderProcessor) Stop() {
	close(u.orders)
}

func (u *OrderProcessor) ProcessOrder(ctx context.Context, number entities.OrderNumber) error {
	if err := number.Validate(); err != nil {
		return err
	}

	tx := u.orderStorager.Tx(ctx)

	order, err := u.orderStorager.Get(ctx, tx, number)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if order.IsProcessed() {
		tx.Rollback(ctx)
		return nil
	}

	user, err := u.userStorager.Get(ctx, tx, order.UserID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	accrualOrder, err := u.accrualStorager.Get(ctx, number)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	updated := order.UpdateWithAccrualOrder(accrualOrder)

	if !updated {
		tx.Rollback(ctx)
		return nil
	}

	user.Balance.Add(&accrualOrder.Accrual)

	err = u.userStorager.Update(ctx, tx, user)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = u.orderStorager.Update(ctx, tx, order)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}
