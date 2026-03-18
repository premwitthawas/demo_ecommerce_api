package product

type ProductTracerEvent string

const (
	TracerProductOutboxRepositoryCreated ProductTracerEvent = "repository.product.outbox.created"
)
