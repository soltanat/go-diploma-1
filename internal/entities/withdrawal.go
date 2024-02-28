package entities

import "time"

type Withdrawal struct {
	OrderNumber OrderNumber
	Sum         Currency
	ProcessedAt time.Time
	UserID      Login
}

func NewWithdrawal(order OrderNumber, sum Currency, userID Login) *Withdrawal {
	return &Withdrawal{
		OrderNumber: order,
		Sum:         sum,
		ProcessedAt: Now(),
		UserID:      userID,
	}
}
