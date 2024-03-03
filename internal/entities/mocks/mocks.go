// Code generated by MockGen. DO NOT EDIT.
// Source: internal/entities/usecases.go
//
// Generated by this command:
//
//	mockgen -source=internal/entities/usecases.go -destination=internal/entities/mocks/mocks.go -package=usecasesmocks -self_package DBConn
//

// Package usecasesmocks is a generated GoMock package.
package usecasesmocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/soltanat/go-diploma-1/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockOrderUseCase is a mock of OrderUseCase interface.
type MockOrderUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUseCaseMockRecorder
}

// MockOrderUseCaseMockRecorder is the mock recorder for MockOrderUseCase.
type MockOrderUseCaseMockRecorder struct {
	mock *MockOrderUseCase
}

// NewMockOrderUseCase creates a new mock instance.
func NewMockOrderUseCase(ctrl *gomock.Controller) *MockOrderUseCase {
	mock := &MockOrderUseCase{ctrl: ctrl}
	mock.recorder = &MockOrderUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUseCase) EXPECT() *MockOrderUseCaseMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderUseCase) CreateOrder(ctx context.Context, orderNumber entities.OrderNumber, userID entities.Login) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, orderNumber, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderUseCaseMockRecorder) CreateOrder(ctx, orderNumber, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderUseCase)(nil).CreateOrder), ctx, orderNumber, userID)
}

// ListOrdersByUserID mocks base method.
func (m *MockOrderUseCase) ListOrdersByUserID(ctx context.Context, userID entities.Login) ([]entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrdersByUserID", ctx, userID)
	ret0, _ := ret[0].([]entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrdersByUserID indicates an expected call of ListOrdersByUserID.
func (mr *MockOrderUseCaseMockRecorder) ListOrdersByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrdersByUserID", reflect.TypeOf((*MockOrderUseCase)(nil).ListOrdersByUserID), ctx, userID)
}

// MockOrderProcessorUseCase is a mock of OrderProcessorUseCase interface.
type MockOrderProcessorUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderProcessorUseCaseMockRecorder
}

// MockOrderProcessorUseCaseMockRecorder is the mock recorder for MockOrderProcessorUseCase.
type MockOrderProcessorUseCaseMockRecorder struct {
	mock *MockOrderProcessorUseCase
}

// NewMockOrderProcessorUseCase creates a new mock instance.
func NewMockOrderProcessorUseCase(ctrl *gomock.Controller) *MockOrderProcessorUseCase {
	mock := &MockOrderProcessorUseCase{ctrl: ctrl}
	mock.recorder = &MockOrderProcessorUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderProcessorUseCase) EXPECT() *MockOrderProcessorUseCaseMockRecorder {
	return m.recorder
}

// ProcessOrder mocks base method.
func (m *MockOrderProcessorUseCase) ProcessOrder(ctx context.Context, number entities.OrderNumber) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessOrder", ctx, number)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessOrder indicates an expected call of ProcessOrder.
func (mr *MockOrderProcessorUseCaseMockRecorder) ProcessOrder(ctx, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessOrder", reflect.TypeOf((*MockOrderProcessorUseCase)(nil).ProcessOrder), ctx, number)
}

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockUserUseCase) Authenticate(ctx context.Context, login entities.Login, password string) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", ctx, login, password)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockUserUseCaseMockRecorder) Authenticate(ctx, login, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockUserUseCase)(nil).Authenticate), ctx, login, password)
}

// GetUser mocks base method.
func (m *MockUserUseCase) GetUser(ctx context.Context, login entities.Login) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, login)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserUseCaseMockRecorder) GetUser(ctx, login any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserUseCase)(nil).GetUser), ctx, login)
}

// Register mocks base method.
func (m *MockUserUseCase) Register(ctx context.Context, login entities.Login, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserUseCaseMockRecorder) Register(ctx, login, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserUseCase)(nil).Register), ctx, login, password)
}

// MockWithdrawalUseCase is a mock of WithdrawalUseCase interface.
type MockWithdrawalUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawalUseCaseMockRecorder
}

// MockWithdrawalUseCaseMockRecorder is the mock recorder for MockWithdrawalUseCase.
type MockWithdrawalUseCaseMockRecorder struct {
	mock *MockWithdrawalUseCase
}

// NewMockWithdrawalUseCase creates a new mock instance.
func NewMockWithdrawalUseCase(ctrl *gomock.Controller) *MockWithdrawalUseCase {
	mock := &MockWithdrawalUseCase{ctrl: ctrl}
	mock.recorder = &MockWithdrawalUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawalUseCase) EXPECT() *MockWithdrawalUseCaseMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockWithdrawalUseCase) Count(ctx context.Context, userID entities.Login) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, userID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockWithdrawalUseCaseMockRecorder) Count(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockWithdrawalUseCase)(nil).Count), ctx, userID)
}

// List mocks base method.
func (m *MockWithdrawalUseCase) List(ctx context.Context, userID entities.Login) ([]entities.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, userID)
	ret0, _ := ret[0].([]entities.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockWithdrawalUseCaseMockRecorder) List(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockWithdrawalUseCase)(nil).List), ctx, userID)
}

// Withdraw mocks base method.
func (m *MockWithdrawalUseCase) Withdraw(ctx context.Context, userID entities.Login, orderNumber entities.OrderNumber, amount entities.Currency) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", ctx, userID, orderNumber, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockWithdrawalUseCaseMockRecorder) Withdraw(ctx, userID, orderNumber, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockWithdrawalUseCase)(nil).Withdraw), ctx, userID, orderNumber, amount)
}

// MockPasswordHasher is a mock of PasswordHasher interface.
type MockPasswordHasher struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordHasherMockRecorder
}

// MockPasswordHasherMockRecorder is the mock recorder for MockPasswordHasher.
type MockPasswordHasherMockRecorder struct {
	mock *MockPasswordHasher
}

// NewMockPasswordHasher creates a new mock instance.
func NewMockPasswordHasher(ctrl *gomock.Controller) *MockPasswordHasher {
	mock := &MockPasswordHasher{ctrl: ctrl}
	mock.recorder = &MockPasswordHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPasswordHasher) EXPECT() *MockPasswordHasherMockRecorder {
	return m.recorder
}

// Compare mocks base method.
func (m *MockPasswordHasher) Compare(hashedPwd, plainPwd []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", hashedPwd, plainPwd)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Compare indicates an expected call of Compare.
func (mr *MockPasswordHasherMockRecorder) Compare(hashedPwd, plainPwd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockPasswordHasher)(nil).Compare), hashedPwd, plainPwd)
}

// Hash mocks base method.
func (m *MockPasswordHasher) Hash(pwd []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", pwd)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hash indicates an expected call of Hash.
func (mr *MockPasswordHasherMockRecorder) Hash(pwd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockPasswordHasher)(nil).Hash), pwd)
}
