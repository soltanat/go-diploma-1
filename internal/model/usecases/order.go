package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/soltanat/go-diploma-1/internal/model"
	"github.com/soltanat/go-diploma-1/internal/model/entities"
)

type OrderUseCase struct {
	orderStorager model.OrderStorager
	userStorager  model.UserStorager
}

func NewOrderUseCase(orderStorager model.OrderStorager, userStorager model.UserStorager) (*OrderUseCase, error) {
	if orderStorager == nil {
		return nil, fmt.Errorf("orderStorager is nil")
	}
	if userStorager == nil {
		return nil, fmt.Errorf("userStorager is nil")
	}
	return &OrderUseCase{
		orderStorager: orderStorager,
		userStorager:  userStorager,
	}, nil
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, orderNumber entities.OrderNumber, userID entities.Login) error {
	if err := orderNumber.Validate(); err != nil {
		return err
	}

	if err := userID.Validate(); err != nil {
		return err
	}

	_, err := u.userStorager.Get(ctx, userID, nil)
	if err != nil {
		return err
	}

	if order, err := u.orderStorager.Get(ctx, orderNumber); err == nil {
		if order.UserID == userID {
			return entities.ExistOrderError{}
		}
		return entities.OrderIsCreatedByAnotherUserError{}
	} else if err != nil {
		if !errors.Is(err, entities.NotFoundError{}) {
			return err
		}
	}

	order := entities.NewOrder(orderNumber, userID)
	return u.orderStorager.Save(ctx, order)
}

func (u *OrderUseCase) ListOrdersByUserID(ctx context.Context, userID entities.Login) ([]entities.Order, error) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}

	_, err := u.userStorager.Get(ctx, userID, nil)
	if err != nil {
		return nil, err
	}

	return u.orderStorager.List(ctx, &userID)
}
