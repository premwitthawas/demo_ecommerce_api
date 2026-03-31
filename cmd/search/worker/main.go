package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	dotnevx "github.com/premwitthawas/demo_ecommerce_api/internals/search/adapter/config"
	search_kafka_message "github.com/premwitthawas/demo_ecommerce_api/internals/search/adapter/messages/kafka"
	search_elasticsearch_repository "github.com/premwitthawas/demo_ecommerce_api/internals/search/adapter/repository/elasticsearch"
	usecase "github.com/premwitthawas/demo_ecommerce_api/internals/search/usecase"
	pkg_elasticsearch "github.com/premwitthawas/demo_ecommerce_api/pkgs/elasticsearch"
	pkg_otel "github.com/premwitthawas/demo_ecommerce_api/pkgs/otel"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("search-worker-service")

func main() {
	cfg := dotnevx.NewConfig()
	tp := pkg_otel.SetupOtelTracer(cfg.GetAPPConfig().OtelURL, "search-service")
	defer func() {
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	groupTopics := "product.created,product.updated,product.deleted"
	client, err := pkg_elasticsearch.NewElasticsearch(cfg.GetAPPConfig().ElasticsearchAddress, cfg.GetAPPConfig().ElasticsearchUsername, cfg.GetAPPConfig().ElasticsearchPassword)
	if err != nil {
		log.Panic(err)
	}
	repository := search_elasticsearch_repository.NewSearchElasticSearchRepository(client, tracer)
	usecase := usecase.NewSearchProductUsecase(tracer, repository)
	consumer := search_kafka_message.NewSearchKafkaMessage(cfg, tracer, groupTopics, "sync-search", usecase)
	go func() {
		if err := consumer.ConsumeMessage(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()
	<-ctx.Done()
	log.Println("Worker shutting down gracefully...")
}
