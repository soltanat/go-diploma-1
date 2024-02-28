// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecases/storager/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=internal/usecases/storager/interfaces.go -destination=internal/usecases/storager/mocks/mocks.go -package=mocks -self_package DBConn
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/soltanat/go-diploma-1/internal/entities"
	storager "github.com/soltanat/go-diploma-1/internal/usecases/storager"
	gomock "go.uber.org/mock/gomock"
)

// MockTx is a mock of Tx interface.
type MockTx struct {
	ctrl     *gomock.Controller
	recorder *MockTxMockRecorder
}

// MockTxMockRecorder is the mock recorder for MockTx.
type MockTxMockRecorder struct {
	mock *MockTx
}

// NewMockTx creates a new mock instance.
func NewMockTx(ctrl *gomock.Controller) *MockTx {
	mock := &MockTx{ctrl: ctrl}
	mock.recorder = &MockTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTx) EXPECT() *MockTxMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockTx) Begin(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Begin indicates an expected call of Begin.
func (mr *MockTxMockRecorder) Begin(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockTx)(nil).Begin), ctx)
}

// Commit mocks base method.
func (m *MockTx) Commit(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTxMockRecorder) Commit(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTx)(nil).Commit), ctx)
}

// Rollback mocks base method.
func (m *MockTx) Rollback(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTxMockRecorder) Rollback(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTx)(nil).Rollback), ctx)
}

// MockOrderStorager is a mock of OrderStorager interface.
type MockOrderStorager struct {
	ctrl     *gomock.Controller
	recorder *MockOrderStoragerMockRecorder
}

// MockOrderStoragerMockRecorder is the mock recorder for MockOrderStorager.
type MockOrderStoragerMockRecorder struct {
	mock *MockOrderStorager
}

// NewMockOrderStorager creates a new mock instance.
func NewMockOrderStorager(ctrl *gomock.Controller) *MockOrderStorager {
	mock := &MockOrderStorager{ctrl: ctrl}
	mock.recorder = &MockOrderStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderStorager) EXPECT() *MockOrderStoragerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockOrderStorager) Get(ctx context.Context, tx storager.Tx, number entities.OrderNumber) (*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, tx, number)
	ret0, _ := ret[0].(*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockOrderStoragerMockRecorder) Get(ctx, tx, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockOrderStorager)(nil).Get), ctx, tx, number)
}

// List mocks base method.
func (m *MockOrderStorager) List(ctx context.Context, tx storager.Tx, userID *entities.Login) ([]entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, tx, userID)
	ret0, _ := ret[0].([]entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockOrderStoragerMockRecorder) List(ctx, tx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockOrderStorager)(nil).List), ctx, tx, userID)
}

// Save mocks base method.
func (m *MockOrderStorager) Save(ctx context.Context, tx storager.Tx, order *entities.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, tx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockOrderStoragerMockRecorder) Save(ctx, tx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOrderStorager)(nil).Save), ctx, tx, order)
}

// Tx mocks base method.
func (m *MockOrderStorager) Tx(ctx context.Context) storager.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", ctx)
	ret0, _ := ret[0].(storager.Tx)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockOrderStoragerMockRecorder) Tx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockOrderStorager)(nil).Tx), ctx)
}

// Update mocks base method.
func (m *MockOrderStorager) Update(ctx context.Context, tx storager.Tx, order *entities.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockOrderStoragerMockRecorder) Update(ctx, tx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockOrderStorager)(nil).Update), ctx, tx, order)
}

// MockUserStorager is a mock of UserStorager interface.
type MockUserStorager struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoragerMockRecorder
}

// MockUserStoragerMockRecorder is the mock recorder for MockUserStorager.
type MockUserStoragerMockRecorder struct {
	mock *MockUserStorager
}

// NewMockUserStorager creates a new mock instance.
func NewMockUserStorager(ctrl *gomock.Controller) *MockUserStorager {
	mock := &MockUserStorager{ctrl: ctrl}
	mock.recorder = &MockUserStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorager) EXPECT() *MockUserStoragerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUserStorager) Get(ctx context.Context, tx storager.Tx, login entities.Login) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, tx, login)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserStoragerMockRecorder) Get(ctx, tx, login any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserStorager)(nil).Get), ctx, tx, login)
}

// Save mocks base method.
func (m *MockUserStorager) Save(ctx context.Context, tx storager.Tx, user *entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserStoragerMockRecorder) Save(ctx, tx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserStorager)(nil).Save), ctx, tx, user)
}

// Tx mocks base method.
func (m *MockUserStorager) Tx(ctx context.Context) storager.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", ctx)
	ret0, _ := ret[0].(storager.Tx)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockUserStoragerMockRecorder) Tx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockUserStorager)(nil).Tx), ctx)
}

// Update mocks base method.
func (m *MockUserStorager) Update(ctx context.Context, tx storager.Tx, user *entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserStoragerMockRecorder) Update(ctx, tx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserStorager)(nil).Update), ctx, tx, user)
}

// MockWithdrawalStorager is a mock of WithdrawalStorager interface.
type MockWithdrawalStorager struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawalStoragerMockRecorder
}

// MockWithdrawalStoragerMockRecorder is the mock recorder for MockWithdrawalStorager.
type MockWithdrawalStoragerMockRecorder struct {
	mock *MockWithdrawalStorager
}

// NewMockWithdrawalStorager creates a new mock instance.
func NewMockWithdrawalStorager(ctrl *gomock.Controller) *MockWithdrawalStorager {
	mock := &MockWithdrawalStorager{ctrl: ctrl}
	mock.recorder = &MockWithdrawalStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawalStorager) EXPECT() *MockWithdrawalStoragerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockWithdrawalStorager) List(ctx context.Context, tx storager.Tx, userID entities.Login) ([]entities.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, tx, userID)
	ret0, _ := ret[0].([]entities.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockWithdrawalStoragerMockRecorder) List(ctx, tx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockWithdrawalStorager)(nil).List), ctx, tx, userID)
}

// Save mocks base method.
func (m *MockWithdrawalStorager) Save(ctx context.Context, tx storager.Tx, withdraw *entities.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, tx, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockWithdrawalStoragerMockRecorder) Save(ctx, tx, withdraw any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockWithdrawalStorager)(nil).Save), ctx, tx, withdraw)
}

// Tx mocks base method.
func (m *MockWithdrawalStorager) Tx(ctx context.Context) storager.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx", ctx)
	ret0, _ := ret[0].(storager.Tx)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockWithdrawalStoragerMockRecorder) Tx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockWithdrawalStorager)(nil).Tx), ctx)
}

// MockAccrualOrderStorager is a mock of AccrualOrderStorager interface.
type MockAccrualOrderStorager struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualOrderStoragerMockRecorder
}

// MockAccrualOrderStoragerMockRecorder is the mock recorder for MockAccrualOrderStorager.
type MockAccrualOrderStoragerMockRecorder struct {
	mock *MockAccrualOrderStorager
}

// NewMockAccrualOrderStorager creates a new mock instance.
func NewMockAccrualOrderStorager(ctrl *gomock.Controller) *MockAccrualOrderStorager {
	mock := &MockAccrualOrderStorager{ctrl: ctrl}
	mock.recorder = &MockAccrualOrderStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualOrderStorager) EXPECT() *MockAccrualOrderStoragerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockAccrualOrderStorager) Get(ctx context.Context, number entities.OrderNumber) (*entities.AccrualOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, number)
	ret0, _ := ret[0].(*entities.AccrualOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAccrualOrderStoragerMockRecorder) Get(ctx, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAccrualOrderStorager)(nil).Get), ctx, number)
}
