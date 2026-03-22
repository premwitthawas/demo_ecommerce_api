package product_test

import (
	"context"

	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/model/outbox"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	"github.com/stretchr/testify/mock"
)

type OutboxProductMockRepository struct {
	mock.Mock
}

func (m *OutboxProductMockRepository) CreateProductOutboxMessage(ctx context.Context, entity *outbox.ProductOutboxMessage) (*outbox.ProductOutboxMessage, error) {
	args := m.Called(ctx, entity)
	var res *outbox.ProductOutboxMessage
	if args.Get(0) != nil {
		res = args.Get(0).(*outbox.ProductOutboxMessage)
	}
	return res, args.Error(1)
}

func (m *OutboxProductMockRepository) UpdateOutboxMessage(ctx context.Context, entity *outbox.ProductOutboxMessage) (*outbox.ProductOutboxMessage, error) {
	args := m.Called(ctx, entity)
	var res *outbox.ProductOutboxMessage
	if args.Get(0) != nil {
		res = args.Get(0).(*outbox.ProductOutboxMessage)
	}
	return res, args.Error(1)
}
func (m *OutboxProductMockRepository) GetOutboxMessagesPendingOrRetrying(ctx context.Context, retry, limit int32) ([]*outbox.ProductOutboxMessage, error) {
	args := m.Called(ctx, retry, limit)
	var res []*outbox.ProductOutboxMessage
	if args.Get(0) != nil {
		res = args.Get(0).([]*outbox.ProductOutboxMessage)
	}
	return res, args.Error(1)
}
func (m *OutboxProductMockRepository) WithTx(tx any) port.ProductOutboxMessageRepository {
	args := m.Called(tx)
	if args.Get(0) == nil {
		return m
	}
	return args.Get(0).(port.ProductOutboxMessageRepository)
}
