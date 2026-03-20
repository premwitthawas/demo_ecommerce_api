package product

import "context"

type ProductMessage interface {
	PublishMessage(ctx context.Context, topic string, msg []byte) error
}
