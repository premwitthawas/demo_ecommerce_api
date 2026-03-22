package product_kafka_message

import (
	"context"
	"encoding/json"
	"fmt"

	config "github.com/premwitthawas/demo_ecommerce_api/internals/product/port/config"
	messages "github.com/premwitthawas/demo_ecommerce_api/internals/product/port/messages"
	pkgs_kafka "github.com/premwitthawas/demo_ecommerce_api/pkgs/kafka"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type productKafkaMessage struct {
	cfg    config.Config
	writer *kafka.Writer
	tp     trace.Tracer
}

type ProducerEvent string

const (
	PublishedMessage ProducerEvent = "worker.product.publish_message"
)

func (p *productKafkaMessage) PublishMessage(ctx context.Context, topic string, msg []byte) error {
	ctx, sp := p.tp.Start(ctx, string(PublishedMessage))
	defer sp.End()
	if len(msg) == 0 {
		err := fmt.Errorf("cannot publish empty message")
		sp.RecordError(err)
		return err
	}
	if !json.Valid(msg) {
		err := fmt.Errorf("invalid json")
		sp.RecordError(err)
		return err
	}
	headerMap := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(headerMap))
	var kafakaHeader []kafka.Header
	for k, v := range headerMap {
		kafakaHeader = append(kafakaHeader, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}
	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic:   topic,
		Value:   msg,
		Headers: kafakaHeader,
	}); err != nil {
		return err
	}
	return nil
}

func (p *productKafkaMessage) Close() error {
	return p.writer.Close()
}

func NewProductKafkaMessage(cfg config.Config, tp trace.Tracer) messages.ProductMessage {
	writer, _ := pkgs_kafka.NewKafkaProducer(cfg.GetAPPConfig().KafkaAddresses)
	return &productKafkaMessage{
		cfg:    cfg,
		tp:     tp,
		writer: writer,
	}
}
