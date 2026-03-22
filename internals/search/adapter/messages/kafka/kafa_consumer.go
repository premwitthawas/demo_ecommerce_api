package search_kafka_message

import (
	"context"
	"encoding/json"
	"log"

	search "github.com/premwitthawas/demo_ecommerce_api/internals/search/domain/product"
	port "github.com/premwitthawas/demo_ecommerce_api/internals/search/port"
	config "github.com/premwitthawas/demo_ecommerce_api/internals/search/port/config"
	portMessage "github.com/premwitthawas/demo_ecommerce_api/internals/search/port/messages"
	pkgs_kafka "github.com/premwitthawas/demo_ecommerce_api/pkgs/kafka"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type searchKafkaMessage struct {
	tp      trace.Tracer
	cfg     config.Config
	reader  *kafka.Reader
	usecsae port.SearchUsecase
}

func (s *searchKafkaMessage) ConsumeMessage(ctx context.Context) error {
	for {
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		headersMap := make(map[string]string)
		for _, h := range msg.Headers {
			headersMap[h.Key] = string(h.Value)
		}
		err = func(ctx context.Context, msg kafka.Message) error {
			childCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(headersMap))
			childCtx, childSp := s.tp.Start(childCtx, "worker.search.consume")
			defer childSp.End()
			payload := new(search.SearchProductMessage)
			if err := json.Unmarshal(msg.Value, payload); err != nil {
				return err
			}
			var processErr error
			switch msg.Topic {
			case "product.created", "product.updated":
				processErr = s.usecsae.SyncProduct(childCtx, payload)
			case "product.deleted":
				processErr = s.usecsae.DeleteProduct(childCtx, payload.ID)
			default:
				log.Printf("Unknown topic: %s", msg.Topic)
			}

			if processErr != nil {
				return processErr
			}
			if err := s.reader.CommitMessages(childCtx, msg); err != nil {
				return err
			}
			return nil
		}(ctx, msg)
		if err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}
}

func NewSearchKafkaMessage(cfg config.Config, tp trace.Tracer, topic, groudID string, usecsae port.SearchUsecase) portMessage.SearchConsumerMessageEvent {
	reader := pkgs_kafka.NewKafkaConsumer(cfg.GetAPPConfig().KafkaAddress, topic, groudID)
	return &searchKafkaMessage{
		tp:      tp,
		cfg:     cfg,
		reader:  reader,
		usecsae: usecsae,
	}
}
