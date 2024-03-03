package entities

//go:generate go-enum --marshal

// ENUM(REGISTERED, INVALID, PROCESSING, PROCESSED)
type AccrualOrderStatus string

type AccrualOrder struct {
	Number  OrderNumber
	Status  AccrualOrderStatus
	Accrual *Currency
}
