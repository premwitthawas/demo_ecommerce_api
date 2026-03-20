package product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/outbox"
	"github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/product"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/product/port"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type productUsecase struct {
	tp                    trace.Tracer
	productRepository     port.ProductRepository
	outboxRepository      port.ProductOutboxMessageRepository
	transactionRepository port.ProductTransactionManagerRepository
}

func generateID() string {
	return uuid.NewString()
}

type EventTypeProduct string

const (
	Created EventTypeProduct = "product.created"
	Delete  EventTypeProduct = "product.deleted"
)

func (p *productUsecase) CreateProduct(ctx context.Context, dto *port.ProductCreateDTO) (*product.Product, error) {
	ctx, cancle := context.WithTimeout(ctx, 30*time.Second)
	defer cancle()
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductUsecaseCreated))
	defer sp.End()
	traceMetadata := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(traceMetadata))
	metadataJSON, _ := json.Marshal(traceMetadata)
	var res *product.Product
	if err := p.transactionRepository.TransactionManager(ctx, func(ctx context.Context, tx any) error {
		ptx := p.productRepository.WithTx(tx)
		otx := p.outboxRepository.WithTx(tx)
		product, err := product.NewProduct(generateID(), dto.Name, dto.Description, product.ProductCategoryType(dto.Category))
		if err != nil {
			return err
		}
		row, err := ptx.CreateProduct(ctx, product)
		if err != nil {
			sp.RecordError(err)
			return err
		}
		payload, err := json.Marshal(row)
		if err != nil {
			sp.RecordError(err)
			return err
		}
		outbox, err := outbox.NewProductOutbox(&outbox.ProductOutboxMessage{
			ID:          generateID(),
			EventType:   string(Created),
			AggrID:      row.ID,
			AggrVersion: row.Version,
			Payload:     payload,
			Metadata:    metadataJSON,
		})
		if err != nil {
			sp.RecordError(err)
			return err
		}
		outbox, err = otx.CreateProductOutboxMessage(ctx, outbox)
		if err != nil {
			sp.RecordError(err)
			return err
		}
		res = row
		return nil
	}); err != nil {
		sp.RecordError(err)
		return nil, err
	}
	return res, nil
}

// DeleteProductByID implements [product.ProductUsecase].
func (p *productUsecase) DeleteProductByID(ctx context.Context, id string) (*product.Product, error) {
	ctx, cancle := context.WithTimeout(ctx, 30*time.Second)
	defer cancle()
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductUsecaseDeleteByID))
	defer sp.End()
	traceMetadata := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(traceMetadata))
	metadataJSON, _ := json.Marshal(traceMetadata)
	var res *product.Product
	if err := p.transactionRepository.TransactionManager(ctx, func(ctx context.Context, tx any) error {
		ptx := p.productRepository.WithTx(tx)
		otx := p.outboxRepository.WithTx(tx)
		row, err := ptx.GetProductByID(ctx, id)
		if err != nil {
			sp.RecordError(err)
			return err
		}
		row, err = ptx.DeleteProductByID(ctx, row.ID, row.Version)
		if err != nil {
			sp.RecordError(err)
			return err
		}
		payload, err := json.Marshal(row)
		if err != nil {
			sp.RecordError(err)
			return product.ErrProductJsonParse
		}
		outbox, err := outbox.NewProductOutbox(&outbox.ProductOutboxMessage{
			ID:          generateID(),
			EventType:   string(Delete),
			AggrID:      row.ID,
			AggrVersion: row.Version,
			Payload:     payload,
			Metadata:    metadataJSON,
		})
		outbox, err = otx.CreateProductOutboxMessage(ctx, outbox)
		if err != nil {
			sp.RecordError(err)
			return err
		}
		res = row
		return nil
	}); err != nil {
		sp.RecordError(err)
		return nil, err
	}
	return res, nil
}

// GetProductByID implements [product.ProductUsecase].
func (p *productUsecase) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	ctx, cancle := context.WithTimeout(ctx, 30*time.Second)
	defer cancle()
	ctx, sp := p.tp.Start(ctx, string(product.TracerProductUsecaseGetByID))
	defer sp.End()
	res, err := p.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewProductUsecase(tp trace.Tracer,
	productRepository port.ProductRepository,
	outboxRepository port.ProductOutboxMessageRepository,
	transactionRepository port.ProductTransactionManagerRepository) port.ProductUsecase {
	return &productUsecase{
		tp:                    tp,
		productRepository:     productRepository,
		outboxRepository:      outboxRepository,
		transactionRepository: transactionRepository,
	}
}
