package product

type ProductTracerEvent string

const (
	TracerProductRepositoryCreated    ProductTracerEvent = "repository.product.created"
	TracerProductRepositoryUpdateByID ProductTracerEvent = "repository.product.updated_by_id"
	TracerProductRepositoryDeleteByID ProductTracerEvent = "repository.product.delete_by_id"
	TracerProductRepositoryGetByID    ProductTracerEvent = "repository.product.get_by_id"
)
