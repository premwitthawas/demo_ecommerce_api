package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	search_dotnevx "github.com/premwitthawas/demo_ecommerce_api/internals/search/adapter/config"
	searc_product_handler "github.com/premwitthawas/demo_ecommerce_api/internals/search/adapter/http/handler"
	search_elasticsearch_repository "github.com/premwitthawas/demo_ecommerce_api/internals/search/adapter/repository/elasticsearch"
	usecase "github.com/premwitthawas/demo_ecommerce_api/internals/search/usecase"
	pkg_elasticsearch "github.com/premwitthawas/demo_ecommerce_api/pkgs/elasticsearch"
	pkg_otel "github.com/premwitthawas/demo_ecommerce_api/pkgs/otel"

	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
)

var tracer = otel.Tracer("search-service")

func main() {
	cfg := search_dotnevx.NewConfig()

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
	g, gCtx := errgroup.WithContext(ctx)
	es, err := pkg_elasticsearch.NewElasticsearch(cfg.GetAPPConfig().ElasticsearchAddress, cfg.GetAPPConfig().ElasticsearchUsername, cfg.GetAPPConfig().ElasticsearchPassword)
	repository := search_elasticsearch_repository.NewSearchElasticSearchRepository(es, tracer)
	usecase := usecase.NewSearchProductUsecase(tracer, repository)
	handler := searc_product_handler.NewSearchProductHandler(tracer, usecase)
	if err != nil {
		log.Panic(err)
	}
	app := fiber.New()
	app.Use(fiberotel.Middleware())
	api := app.Group("/api/v1/search")
	api.Get("/products", handler.GetProducts)
	g.Go(func() error {
		log.Printf("[search][info]: app listening at %s \n", cfg.GetAPPConfig().Address)
		return app.Listen(cfg.GetAPPConfig().Address)
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Println("[search][info]: signal graceful shutdown recived")
		gracfulCtx, gracefulStop := context.WithTimeout(context.Background(), time.Second*30)
		defer gracefulStop()
		return app.ShutdownWithContext(gracfulCtx)
	})
	if err := g.Wait(); err != nil {
		log.Printf("[search][error]: group exited with error: %v \n", err)
	}
	log.Println("[search][info]: graceful shutdown success.")
}
