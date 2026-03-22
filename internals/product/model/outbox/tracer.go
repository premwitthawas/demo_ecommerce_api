package product_outbox

type ProductTracerEvent string

const (
	TracerProductOutboxRepositoryCreated                            ProductTracerEvent = "repository.product.outbox.created"
	TracerProductOutboxRepositoryGetOutboxMessagesPendingOrRetrying ProductTracerEvent = "repository.product.outbox.get_messages_peding_or_retrying"
	TracerProductOutboxRepositoryUpdateMessage                      ProductTracerEvent = "repository.product.outbox.update_message"
)
