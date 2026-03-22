package search

import "context"

type SearchConsumerMessageEvent interface {
	ConsumeMessage(ctx context.Context) error
}
