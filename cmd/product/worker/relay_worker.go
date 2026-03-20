package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	dotnevx "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/config"
	product_kafka_message "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/messages/kafka"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("product-worker-service")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cfg := dotnevx.NewConfig()
	producer := product_kafka_message.NewProductKafkaMessage(cfg, tracer)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("Worker working...")
				_ = producer
			}()
		case <-ctx.Done():
			log.Println("Worker shutting down gracefully...")
			return
		}
	}
}
