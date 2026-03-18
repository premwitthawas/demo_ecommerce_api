package product

import "errors"

var (
	// Domain Validation Errors
	ErrProductIDEmpty   = errors.New("product: id is empty")
	ErrProductNameEmpty = errors.New("product: name is empty")
	ErrProductDescEmpty = errors.New("product: description is empty")

	// Repository / State Errors
	ErrProductNotFound    = errors.New("product: not found")
	ErrProductConflict    = errors.New("product: optimistic lock conflict (version mismatch)")
	ErrProductPersistence = errors.New("product: persistence error (database issues)")
)
