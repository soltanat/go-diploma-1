package http

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/http/api"
	"github.com/soltanat/go-diploma-1/internal/logger"
)

type ServerInterfaceWrapper struct {
	userUseCase       entities.UserUseCase
	orderUseCase      entities.OrderUseCase
	withdrawalUseCase entities.WithdrawalUseCase
	tokenProvider     TokenProvider
}

func NewServerInterfaceWrapper(
	userUseCase entities.UserUseCase,
	orderUseCase entities.OrderUseCase,
	withdrawalUseCase entities.WithdrawalUseCase,
	tokenProvider TokenProvider,
) *ServerInterfaceWrapper {
	return &ServerInterfaceWrapper{
		userUseCase:       userUseCase,
		orderUseCase:      orderUseCase,
		withdrawalUseCase: withdrawalUseCase,
		tokenProvider:     tokenProvider,
	}
}

func (h *ServerInterfaceWrapper) RegisterUser(ctx context.Context, request api.RegisterUserRequestObject) (api.RegisterUserResponseObject, error) {
	l := logger.Get()

	err := h.userUseCase.Register(ctx, entities.Login(request.Body.Login), request.Body.Password)
	if err != nil {
		validationErr := &entities.ValidationError{}
		if errors.As(err, validationErr) {
			return api.RegisterUser400Response{}, nil
		}
		existErr := &entities.ExistUserError{}
		if errors.As(err, existErr) {
			return api.RegisterUser409Response{}, nil
		}
		l.Err(err).Msg("failed to register user")
		return api.RegisterUser500Response{}, nil
	}

	return api.RegisterUser200Response{}, nil
}

func (h *ServerInterfaceWrapper) LoginUser(ctx context.Context, request api.LoginUserRequestObject) (api.LoginUserResponseObject, error) {
	l := logger.Get()

	_, err := h.userUseCase.Authenticate(ctx, entities.Login(request.Body.Login), request.Body.Password)
	if err != nil {
		validationErr := &entities.ValidationError{}
		if errors.As(err, validationErr) {
			return api.LoginUser400Response{}, nil
		}
		notFoundErr := &entities.NotFoundError{}
		if errors.As(err, notFoundErr) {
			return api.LoginUser401Response{}, nil
		}
		pwdErr := &entities.InvalidPasswordError{}
		if errors.As(err, pwdErr) {
			return api.LoginUser401Response{}, nil
		}
		l.Err(err).Msg("failed to login user")
		return api.LoginUser500Response{}, nil
	}

	token, err := h.tokenProvider.GenerateToken(request.Body.Login)
	if err != nil {
		l.Err(err).Msg("failed to generate token")
		return api.LoginUser500Response{}, nil
	}

	return api.LoginUser200Response{
		Headers: api.LoginUser200ResponseHeaders{
			Authorization: "Bearer " + token,
		},
	}, nil
}

func (h *ServerInterfaceWrapper) CreateOrder(ctx context.Context, request api.CreateOrderRequestObject) (api.CreateOrderResponseObject, error) {
	l := logger.Get()

	userID := ctx.Value(userIDKeyStruct).(string)

	orderNumber, err := strconv.Atoi(*request.Body)
	if err != nil {
		return api.CreateOrder400Response{}, nil
	}

	err = h.orderUseCase.CreateOrder(ctx, entities.OrderNumber(orderNumber), entities.Login(userID))
	if err != nil {
		existErr := &entities.ExistOrderError{}
		if errors.As(err, existErr) {
			return api.CreateOrder200Response{}, nil
		}
		validationErr := &entities.ValidationError{}
		if errors.As(err, validationErr) {
			return api.CreateOrder400Response{}, nil
		}
		existAnotherErr := &entities.OrderIsCreatedByAnotherUserError{}
		if errors.As(err, existAnotherErr) {
			return api.CreateOrder409Response{}, nil
		}
		orderNumberErr := &entities.InvalidOrderNumberError{}
		if errors.As(err, orderNumberErr) {
			return api.CreateOrder422Response{}, nil
		}
		l.Err(err).Msg("failed to create order")
		return api.CreateOrder500Response{}, nil
	}

	return api.CreateOrder202Response{}, nil
}

func (h *ServerInterfaceWrapper) GetOrders(ctx context.Context, request api.GetOrdersRequestObject) (api.GetOrdersResponseObject, error) {
	l := logger.Get()

	userID := ctx.Value(userIDKeyStruct).(string)

	oo, err := h.orderUseCase.ListOrdersByUserID(ctx, entities.Login(userID))
	if err != nil {
		l.Err(err).Msg("failed to get orders")
		return api.GetOrders500Response{}, nil
	}

	if len(oo) == 0 {
		return api.GetOrders204Response{}, nil
	}

	response := api.GetOrders200JSONResponse{}
	for _, o := range oo {
		apiOrder := api.Order{
			Number:     int(o.Number),
			Status:     api.OrderStatus(o.Status.String()),
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}
		if o.IsProcessed() {
			accrual := o.Accrual.String()
			apiOrder.Accrual = &accrual
		}
		response = append(response, apiOrder)
	}

	return response, nil
}

func (h *ServerInterfaceWrapper) GetBalance(ctx context.Context, request api.GetBalanceRequestObject) (api.GetBalanceResponseObject, error) {
	l := logger.Get()

	userID := ctx.Value(userIDKeyStruct).(string)

	user, err := h.userUseCase.GetUser(ctx, entities.Login(userID))
	if err != nil {
		l.Err(err).Msg("failed to get balance")
		return api.GetBalance500Response{}, nil
	}

	countWithdrawals, err := h.withdrawalUseCase.Count(ctx, entities.Login(userID))
	if err != nil {
		l.Err(err).Msg("failed to get balance")
		return api.GetBalance500Response{}, nil
	}

	return api.GetBalance200JSONResponse{
		Current:     user.Balance.String(),
		Withdrawals: countWithdrawals,
	}, nil
}

func (h *ServerInterfaceWrapper) Withdraw(ctx context.Context, request api.WithdrawRequestObject) (api.WithdrawResponseObject, error) {
	l := logger.Get()
	userID := ctx.Value(userIDKeyStruct).(string)

	err := h.withdrawalUseCase.Withdraw(
		ctx, entities.Login(userID), entities.OrderNumber(request.Body.Order), entities.FromString(request.Body.Sum),
	)
	if err != nil {
		outOfBalanceErr := &entities.OutOfBalanceError{}
		if errors.As(err, outOfBalanceErr) {
			return api.Withdraw402Response{}, nil
		}
		orderNumberErr := &entities.InvalidOrderNumberError{}
		if errors.As(err, orderNumberErr) {
			return api.Withdraw422Response{}, nil
		}
		l.Err(err).Msg("failed to withdraw")
		return api.Withdraw500Response{}, nil
	}

	return api.Withdraw200Response{}, nil
}

func (h *ServerInterfaceWrapper) GetWithdrawals(ctx context.Context, request api.GetWithdrawalsRequestObject) (api.GetWithdrawalsResponseObject, error) {
	l := logger.Get()
	userID := ctx.Value(userIDKeyStruct).(string)

	ww, err := h.withdrawalUseCase.List(ctx, entities.Login(userID))
	if err != nil {
		l.Err(err).Msg("failed to get withdrawals")
		return api.GetWithdrawals500Response{}, nil
	}
	if len(ww) == 0 {
		return api.GetWithdrawals204Response{}, nil
	}

	response := api.GetWithdrawals200JSONResponse{}
	for _, w := range ww {
		apiWithdrawal := api.Withdrawal{
			Order:       int(w.OrderNumber),
			ProcessedAt: w.ProcessedAt,
			Sum:         w.Sum.String(),
		}
		response = append(response, apiWithdrawal)
	}

	return response, nil
}
