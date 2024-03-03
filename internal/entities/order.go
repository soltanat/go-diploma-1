package entities

import (
	"time"

	"github.com/theplant/luhn"
	_ "github.com/theplant/luhn"
)

//go:generate go-enum --marshal

// ENUM(NEW, PROCESSING, INVALID, PROCESSED)
type OrderStatus string

type Order struct {
	Number     OrderNumber
	Status     OrderStatus
	Accrual    Currency
	UploadedAt time.Time
	UserID     Login
}

func NewOrder(number OrderNumber, userID Login) *Order {
	return &Order{
		Number:     number,
		Status:     OrderStatusNEW,
		Accrual:    Currency{0, 0},
		UserID:     userID,
		UploadedAt: Now(),
	}
}

func (o *Order) IsProcessed() bool {
	return o.Status == OrderStatusPROCESSED
}

func (o *Order) Validate() error {
	if err := o.Number.Validate(); err != nil {
		return err
	}
	if err := o.UserID.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *Order) UpdateWithAccrualOrder(accrualOrder *AccrualOrder) bool {
	var newStatus OrderStatus

	switch accrualOrder.Status {
	case AccrualOrderStatusREGISTERED:
		newStatus = OrderStatusNEW
	case AccrualOrderStatusINVALID:
		newStatus = OrderStatusINVALID
	case AccrualOrderStatusPROCESSING:
		newStatus = OrderStatusPROCESSING
	case AccrualOrderStatusPROCESSED:
		newStatus = OrderStatusPROCESSED
		if accrualOrder.Accrual != nil {
			o.Accrual = *accrualOrder.Accrual
		}
	}

	if newStatus == o.Status {
		return false
	}

	o.Status = newStatus

	return true
}

type OrderNumber int

func (n OrderNumber) Validate() error {
	if !luhn.Valid(int(n)) {
		return InvalidOrderNumberError{}
	}
	return nil
}
