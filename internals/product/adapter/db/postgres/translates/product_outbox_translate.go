package translate_product_postgresdb

import (
	"github.com/jackc/pgx/v5/pgtype"
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres/product/sqlc"
	outbox "github.com/premwitthawas/demo_ecommerce_api/internals/product/domain/outbox"
)

func ProductOutboxRepositoryTranslateCreate(entity *outbox.ProductOutboxMessage) *product_postgresdb.CreateProductOutboxParams {
	payload := &product_postgresdb.CreateProductOutboxParams{
		ID:          entity.ID,
		EventType:   entity.EventType,
		AggrID:      entity.AggrID,
		AggrVersion: entity.AggrVersion,
		Status:      entity.Status,
		Payload:     entity.Payload,
		Metadata:    entity.Metadata,
		NextRetryAt: entity.NextRetryAt,
		RetryCount:  entity.RetryCount,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
	if entity.Version > 0 {
		payload.Version = entity.Version
	}
	if entity.ErrText != nil {
		payload.ErrText = entity.ErrText
	}
	if entity.ConsumedAt != nil {
		payload.ConsumedAt = pgtype.Timestamptz{
			Time:  *entity.ConsumedAt,
			Valid: true,
		}
	}
	return payload
}

func ProductOutboxRepositoryTranslateRowToDomain(row *product_postgresdb.OutboxMessage) *outbox.ProductOutboxMessage {
	payload := &outbox.ProductOutboxMessage{
		ID:          row.ID,
		EventType:   row.EventType,
		AggrID:      row.AggrID,
		AggrVersion: row.AggrVersion,
		Status:      row.Status,
		Payload:     row.Payload,
		Metadata:    row.Metadata,
		RetryCount:  row.RetryCount,
		NextRetryAt: row.NextRetryAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Version:     row.Version,
	}
	if row.ErrText != nil {
		payload.ErrText = row.ErrText
	}
	if row.ConsumedAt.Valid {
		payload.ConsumedAt = &row.ConsumedAt.Time
	}
	return payload
}
