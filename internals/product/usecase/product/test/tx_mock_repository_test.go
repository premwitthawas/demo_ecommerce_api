package product_test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TxRepositoryMock struct {
	mock.Mock
}

func (m *TxRepositoryMock) TransactionManager(ctx context.Context, handler func(ctx context.Context, tx any) error) error {
	args := m.Called(ctx, handler)
	if handler != nil {
		_ = handler(ctx, nil)
	}
	return args.Error(0)
}

func (m *TxRepositoryMock) MockSuccess() {
	m.On("TransactionManager", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context, any) error)
			_ = fn(args.Get(0).(context.Context), "mock_tx")
		}).
		Return(nil)
}

func (m *TxRepositoryMock) MockSuccessWithError(returnErr error) {
	m.On("TransactionManager", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context, any) error)
			_ = fn(args.Get(0).(context.Context), "mock_tx")
		}).
		Return(returnErr)
}

func (m *TxRepositoryMock) MockError(err error) {
	m.On("TransactionManager", mock.Anything, mock.Anything).Return(err)
}
