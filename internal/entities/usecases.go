package entities

import "context"

type OrderUseCase interface {
	CreateOrder(ctx context.Context, orderNumber OrderNumber, userID Login) error
	ListOrdersByUserID(ctx context.Context, userID Login) ([]Order, error)
}

type OrderProcessorUseCase interface {
	ProcessOrder(ctx context.Context, number OrderNumber) error
}

type UserUseCase interface {
	Register(ctx context.Context, login Login, password string) error
	Authenticate(ctx context.Context, login Login, password string) (*User, error)
	GetUser(ctx context.Context, login Login) (*User, error)
}

type WithdrawalUseCase interface {
	Count(ctx context.Context, userID Login) (int, error)
	List(ctx context.Context, userID Login) ([]Withdrawal, error)
	Withdraw(ctx context.Context, userID Login, orderNumber OrderNumber, amount Currency) error
}
