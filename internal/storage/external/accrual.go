package external

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/soltanat/go-diploma-1/internal/clients/accrual"
	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type AccrualStorage struct {
	client accrual.ClientWithResponsesInterface
}

func NewAccrualStorage(client accrual.ClientWithResponsesInterface) (storager.AccrualOrderStorager, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	return &AccrualStorage{
		client: client,
	}, nil
}

func (s *AccrualStorage) Get(ctx context.Context, number entities.OrderNumber) (*entities.AccrualOrder, error) {
	sNumber := strconv.Itoa(int(number))
	accrualOrder, err := s.client.GetOrderWithResponse(ctx, sNumber)
	if err != nil {
		return nil, entities.StorageError{Err: err}
	}
	if accrualOrder.StatusCode() == http.StatusNoContent {
		return nil, entities.NotFoundError{}
	}
	if accrualOrder.StatusCode() != http.StatusOK {
		return nil, entities.StorageError{Err: fmt.Errorf("unexpected status code: %d", accrualOrder.StatusCode())}
	}

	var currency *entities.Currency = nil
	if accrualOrder.JSON200.Accrual != nil {
		currencyPtr := entities.CurrencyFromFloat(*accrualOrder.JSON200.Accrual)
		currency = &currencyPtr
	}

	status := entities.AccrualOrderStatus(*accrualOrder.JSON200.Status)
	if !status.IsValid() {
		return nil, fmt.Errorf("accrual oreder invalid status: %s", *accrualOrder.JSON200.Status)
	}

	result := &entities.AccrualOrder{
		Number:  number,
		Status:  status,
		Accrual: currency,
	}

	return result, nil
}
