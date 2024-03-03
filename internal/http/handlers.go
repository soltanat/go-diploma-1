package http

import (
	"context"
	"strconv"

	"github.com/lestrrat-go/jwx/jwt"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/http/api"
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

func (h *ServerInterfaceWrapper) GetBalance(ctx context.Context, request api.GetBalanceRequestObject) (api.GetBalanceResponseObject, error) {
	userID := ctx.Value(userIDKey).(string)

	user, err := h.userUseCase.GetUser(ctx, entities.Login(userID))
	if err != nil {
		//	TODO: что с ошибками
	}

	countWithdrawals, err := h.withdrawalUseCase.Count(ctx, entities.Login(userID))
	if err != nil {
		//	TODO: что с ошибками
	}

	return api.GetBalance200JSONResponse{
		Current:     user.Balance.String(),
		Withdrawals: countWithdrawals,
	}, nil
}

func (h *ServerInterfaceWrapper) Withdraw(ctx context.Context, request api.WithdrawRequestObject) (api.WithdrawResponseObject, error) {
	userID := ctx.Value(userIDKey).(string)

	err := h.withdrawalUseCase.Withdraw(
		ctx, entities.Login(userID), entities.OrderNumber(request.Body.Order), entities.FromString(request.Body.Sum),
	)
	if err != nil {
		//	TODO: что с ошибками
	}

	return api.Withdraw200Response{}, nil
}

func (h *ServerInterfaceWrapper) LoginUser(ctx context.Context, request api.LoginUserRequestObject) (api.LoginUserResponseObject, error) {
	_, err := h.userUseCase.Authenticate(ctx, entities.Login(request.Body.Login), request.Body.Password)
	if err != nil {
		//	TODO: что с ошибками
	}

	token, err := h.tokenProvider.GenerateToken(request.Body.Login)
	if err != nil {
		//	TODO: что с ошибками
	}

	return api.LoginUser200Response{
		Headers: api.LoginUser200ResponseHeaders{
			Authorization: "Bearer " + token,
		},
	}, nil
}

func (h *ServerInterfaceWrapper) GetOrders(ctx context.Context, request api.GetOrdersRequestObject) (api.GetOrdersResponseObject, error) {
	userID := ctx.Value(userIDKey).(string)

	oo, err := h.orderUseCase.ListOrdersByUserID(ctx, entities.Login(userID))
	if err != nil {
		//	TODO: что с ошибками
	}

	response := api.GetOrders200JSONResponse{}
	for _, o := range oo {
		apiOrder := api.Order{
			Number:     int(o.Number),
			Status:     api.OrderStatus(o.Status.String()),
			UploadedAt: o.UploadedAt,
		}
		if o.IsProcessed() {
			accrual := o.Accrual.String()
			apiOrder.Accrual = &accrual
		}
		response = append(response, apiOrder)
	}

	return response, nil
}

func (h *ServerInterfaceWrapper) CreateOrder(ctx context.Context, request api.CreateOrderRequestObject) (api.CreateOrderResponseObject, error) {
	userID := ctx.Value(userIDKey).(string)

	orderNumber, err := strconv.Atoi(*request.Body)
	if err != nil {
		//	TODO: что с ошибками
	}

	err = h.orderUseCase.CreateOrder(ctx, entities.OrderNumber(orderNumber), entities.Login(userID))
	if err != nil {
		//	TODO: что с ошибками
	}

	return api.CreateOrder200Response{}, nil
}

func (h *ServerInterfaceWrapper) RegisterUser(ctx context.Context, request api.RegisterUserRequestObject) (api.RegisterUserResponseObject, error) {
	err := h.userUseCase.Register(ctx, entities.Login(request.Body.Login), request.Body.Password)
	if err != nil {
		//	TODO: что с ошибками
	}
	return api.RegisterUser200JSONResponse{}, nil
}

func (h *ServerInterfaceWrapper) GetWithdrawals(ctx context.Context, request api.GetWithdrawalsRequestObject) (api.GetWithdrawalsResponseObject, error) {
	userID := ctx.Value(userIDKey).(string)

	ww, err := h.withdrawalUseCase.List(ctx, entities.Login(userID))
	if err != nil {
		//	TODO: что с ошибками
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

func GetSubFromJWT(token string) (string, error) {
	parsedToken, err := jwt.Parse([]byte(token))
	if err != nil {
		return "", err
	}
	return parsedToken.Subject(), nil
}
