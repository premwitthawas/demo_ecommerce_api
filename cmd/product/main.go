package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	product_dotenvx "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/config"
	product_postgresdb "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/db/postgres"
	product_http "github.com/premwitthawas/demo_ecommerce_api/internals/product/adapter/http"
	product_usecase "github.com/premwitthawas/demo_ecommerce_api/internals/product/usecase/product"
	pkg_otel "github.com/premwitthawas/demo_ecommerce_api/pkgs/otel"
	pkg_postgres "github.com/premwitthawas/demo_ecommerce_api/pkgs/postgres"

	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
)

var tracer = otel.Tracer("product-service")

func main() {
	tp := pkg_otel.SetupOtelTracer("localhost:4317", "product-service")
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
	cfg := product_dotenvx.NewConfig()
	pool, _ := pkg_postgres.NewPostgresPool(context.Background(), cfg.GetDBConfig().DatabaseURL)
	txRepo := product_postgresdb.NewTransactionManger(pool, tracer)
	outboxRepo := product_postgresdb.NewProductOutboxRepository(pool, tracer)
	productRepo := product_postgresdb.NewProductRepository(pool, tracer)
	usecase := product_usecase.NewProductUsecase(tracer, productRepo, outboxRepo, txRepo)
	handler := product_http.NewHttpHandler(tracer, usecase)
	app := fiber.New()
	app.Use(fiberotel.Middleware())
	api := app.Group("/api/v1/products")
	api.Post("/", handler.CreateProduct)
	api.Get(":id", handler.GetProductByID)
	api.Delete(":id", handler.DeleteProductByID)
	g.Go(func() error {
		log.Printf("[product][info]: app listening at %s \n", cfg.GetAPPConfig().Address)
		return app.Listen(cfg.GetAPPConfig().Address)
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Println("[product][info]: signal graceful shutdown recived")
		gracfulCtx, gracefulStop := context.WithTimeout(context.Background(), time.Second*30)
		defer gracefulStop()
		return app.ShutdownWithContext(gracfulCtx)
	})
	if err := g.Wait(); err != nil {
		log.Printf("[product][error]: group exited with error: %v \n", err)
	}
	log.Println("[product][info]: graceful shutdown success.")
}
