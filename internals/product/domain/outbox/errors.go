package product

import "errors"

var (
	// Domain Validation Errors
	ErrProductOutboxIDEmpty                 = errors.New("product.outbox: id is empty")
	ErrProductOutboxAggrIDEmpty             = errors.New("product.outbox: aggr_id is empty")
	ErrProductOutboxAggrVersionLessthanZero = errors.New("product.outbox: aggr_version less than zero")
	ErrProductOutboxEventTypeEmpty          = errors.New("product.outbox: event_type is empty")
	ErrProductOutboxPayloadEmpty            = errors.New("product.outbox: payload is empty")
	ErrProductOutboxMetadataEmpty           = errors.New("product.outbox: metadata is empty")
	ErrProductOutboxStatustypeEmpty         = errors.New("product.outbox: status is empty")

	// Repository / State Errors
	ErrProductOutboxNotFound    = errors.New("product.outbox: not found")
	ErrProductOutboxConflict    = errors.New("product.outbox: optimistic lock conflict (version mismatch)")
	ErrProductOutboxPersistence = errors.New("product.outbox: persistence error (database issues)")
)
