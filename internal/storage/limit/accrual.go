package limit

import (
	"context"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type Accrual struct {
	storage storager.AccrualOrderStorager
	limit   chan struct{}
}

func NewLimitAccrualStorage(storage storager.AccrualOrderStorager, limit int) storager.AccrualOrderStorager {
	return &Accrual{
		storage: storage,
		limit:   make(chan struct{}, limit),
	}
}

func (l *Accrual) Get(ctx context.Context, number entities.OrderNumber) (*entities.AccrualOrder, error) {
	if err := number.Validate(); err != nil {
		return nil, err
	}

	l.limit <- struct{}{}
	defer func() {
		<-l.limit
	}()

	return l.storage.Get(ctx, number)
}
