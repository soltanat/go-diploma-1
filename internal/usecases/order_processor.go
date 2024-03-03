package usecases

import (
	"context"
	"fmt"

	"github.com/soltanat/go-diploma-1/internal/logger"

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

func (u *OrderProcessor) Produce(ctx context.Context) error {
	l := logger.Get()
	l.Debug().Msg("start produce orders")
	orders, err := u.orderStorager.List(
		ctx, nil, nil, &[]entities.OrderStatus{entities.OrderStatusNEW, entities.OrderStatusPROCESSING},
	)
	if err != nil {
		return err
	}
	for _, order := range orders {
		u.orders <- order
	}
	l.Debug().Msgf("produced %d orders", len(orders))
	return nil
}

func (u *OrderProcessor) AddOrder(order entities.Order) {
	u.orders <- order
}

func (u *OrderProcessor) Run(ctx context.Context) {
	for order := range u.orders {
		if ctx.Err() != nil {
			return
		}
		l := logger.Get()
		if err := u.ProcessOrder(ctx, order.Number); err != nil {
			l.Error().Err(err).Msgf("failed process order %d", order.Number)
		}
	}
}

func (u *OrderProcessor) Stop() {
	close(u.orders)
}

func (u *OrderProcessor) ProcessOrder(ctx context.Context, number entities.OrderNumber) error {
	l := logger.Get()
	l.Debug().Str("usecase", "ProcessOrder").Msgf("start process order %d", number)

	if err := number.Validate(); err != nil {
		return err
	}

	tx := u.orderStorager.Tx(ctx)
	err := tx.Begin(ctx)
	if err != nil {
		return err
	}

	order, err := u.orderStorager.Get(ctx, tx, number)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	if order.IsProcessed() {
		l.Debug().Str("usecase", "ProcessOrder").Msgf("order %d is already processed", number)
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	user, err := u.userStorager.Get(ctx, tx, order.UserID)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	l.Debug().Str("usecase", "ProcessOrder").Msgf("found user %s balance %v", user.Login, user.Balance)

	accrualOrder, err := u.accrualStorager.Get(ctx, number)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	l.Debug().Str("usecase", "ProcessOrder").Msgf("found accrual order %d, currency %v", accrualOrder.Number, accrualOrder.Accrual)

	updated := order.UpdateWithAccrualOrder(accrualOrder)

	if !updated {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	if accrualOrder.Accrual != nil {
		user.Balance.Add(accrualOrder.Accrual)
	}

	err = u.userStorager.Update(ctx, tx, user)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}
	l.Debug().Str("usecase", "ProcessOrder").Msgf("updated user %s balance %v", user.Login, user.Balance)

	err = u.orderStorager.Update(ctx, tx, order)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	l.Debug().Str("usecase", "ProcessOrder").Msgf("processed order %d", number)
	return nil
}
