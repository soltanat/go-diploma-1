package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/soltanat/go-diploma-1/internal/logger"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type OrderUseCase struct {
	orderStorager  storager.OrderStorager
	userStorager   storager.UserStorager
	OrderProcessor entities.OrderProcessorUseCase
}

func NewOrderUseCase(orderStorager storager.OrderStorager, userStorager storager.UserStorager, orderProcessor entities.OrderProcessorUseCase) (*OrderUseCase, error) {
	if orderStorager == nil {
		return nil, fmt.Errorf("orderStorager is nil")
	}
	if userStorager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	if orderProcessor == nil {
		return nil, fmt.Errorf("orderProcessor is nil")
	}
	return &OrderUseCase{
		orderStorager:  orderStorager,
		userStorager:   userStorager,
		OrderProcessor: orderProcessor,
	}, nil
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, orderNumber entities.OrderNumber, userID entities.Login) error {
	l := logger.Get()

	if err := orderNumber.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	_, err := u.userStorager.Get(ctx, nil, userID)
	if err != nil {
		return err
	}

	if order, err := u.orderStorager.Get(ctx, nil, orderNumber); err == nil {
		if order.UserID == userID {
			return entities.ExistOrderError{}
		}
		return entities.OrderIsCreatedByAnotherUserError{}
	} else if !errors.Is(err, entities.NotFoundError{}) {
		return err
	}

	order := entities.NewOrder(orderNumber, userID)
	err = u.orderStorager.Save(ctx, nil, order)
	if err != nil {
		return err
	}

	err = u.OrderProcessor.ProcessOrder(ctx, orderNumber)
	if err != nil {
		return err
	}

	l.Debug().Str("usecase", "CreateOrder").Msgf("created order %d", orderNumber)

	return nil
}

func (u *OrderUseCase) ListOrdersByUserID(ctx context.Context, userID entities.Login) ([]entities.Order, error) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}

	_, err := u.userStorager.Get(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	return u.orderStorager.List(ctx, nil, &userID, nil)
}
