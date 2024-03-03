package external

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/soltanat/go-diploma-1/internal/usecases/storager"

	"github.com/soltanat/go-diploma-1/internal/clients/accrual"
	"github.com/soltanat/go-diploma-1/internal/entities"
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

	var currency *entities.Currency = nil
	if accrualOrder.JSON200.Accrual != nil {
		accrualString := fmt.Sprintf("%.2f", *accrualOrder.JSON200.Accrual)
		currencyPtr := entities.FromString(accrualString)
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
