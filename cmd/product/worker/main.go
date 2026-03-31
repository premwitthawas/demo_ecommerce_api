package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	dotnevx "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/config"
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres"
	product_kafka_message "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/messages/kafka"
	pkg_otel "github.com/premwitthawas/demo_ecommerce_api/pkgs/otel"
	pkg_postgres "github.com/premwitthawas/demo_ecommerce_api/pkgs/postgres"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("product-worker-service")

func main() {
	cfg := dotnevx.NewConfig()
	tp := pkg_otel.SetupOtelTracer(cfg.GetAPPConfig().OtelURL, "product-worker-service")
	defer func() {
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	poolCtx, stopPool := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopPool()
	pool, _ := pkg_postgres.NewPostgresPool(poolCtx, cfg.GetDBConfig().DatabaseURL)
	defer pool.Close()
	txRepo := product_postgresdb.NewTransactionManger(pool, tracer)
	outboxRepo := product_postgresdb.NewProductOutboxRepository(pool, tracer)
	producer := product_kafka_message.NewProductKafkaMessage(cfg, tracer)
	defer producer.Close()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
			defer stop()
			_ = txRepo.TransactionManager(ctx, func(ctx context.Context, tx any) error {
				repo := outboxRepo.WithTx(tx)
				messages, err := repo.GetOutboxMessagesPendingOrRetrying(ctx, 5, 50)
				if err != nil {
					log.Printf("Error fetching messages: %v", err)
					return err
				}
				if len(messages) == 0 {
					return nil
				}
				for _, msg := range messages {
					var metadataMap map[string]string
					_ = json.Unmarshal(msg.Metadata, &metadataMap)
					parentCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(metadataMap))
					spanContext := trace.SpanFromContext(parentCtx).SpanContext()
					if !spanContext.IsValid() {
						log.Println("⚠️ Warning: Trace Context is INVALID - Check Injector at API side")
					} else {
						log.Printf("🔗 Tracing Linked: TraceID=%s", spanContext.TraceID().String())
					}
					workerCtx, span := tracer.Start(parentCtx, "product-worker.publish-event")
					errPublish := producer.PublishMessage(workerCtx, msg.EventType, msg.Payload)
					if errPublish == nil { // ถ้าส่งสำเร็จ
						msg.UpdateOutboxPublished()
						log.Printf("✅ Published: %s", msg.ID)
					} else {
						log.Printf("❌ Failed to publish %s: %v", msg.ID, errPublish)
						if msg.RetryCount >= 5 {
							msg.UpdateOutboxDLQ(errPublish)
						} else {
							msg.UpdateOutboxRetrying(errPublish)
						}
					}
					updatedMsg, errUpdate := repo.UpdateOutboxMessage(workerCtx, msg)
					if errUpdate != nil {
						log.Printf("❌ DB Update Error ID %s: %v", msg.ID, errUpdate)
						return errUpdate
					}
					log.Printf("✅ Published & Updated DB ID: %s (Status: %s)", updatedMsg.ID, updatedMsg.Status)
					span.End()
				}
				return nil
			})
		case <-ctx.Done():
			log.Println("Worker shutting down gracefully...")
			return
		}
	}
}
