package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	dotnevx "github.com/premwitthawas/demo_ecommerce_api/internals/auth/adapter/config"
	pkg_otel "github.com/premwitthawas/demo_ecommerce_api/pkgs/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/sync/errgroup"
)

var tracer = otel.Tracer("auth-service")

func main() {
	tp := pkg_otel.SetupOtelTracer("localhost:4317", "auth-service")
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
	cfg := dotnevx.NewConfig()
	app := fiber.New()
	api := app.Group("/api/v1/auth")
	api.Get("/", func(c fiber.Ctx) error {
		headersMap := c.GetReqHeaders()
		cleanHeaders := make(map[string]string)
		for k, v := range headersMap {
			if len(v) > 0 {
				cleanHeaders[strings.ToLower(k)] = v[0]
			}
		}
		ctx := otel.GetTextMapPropagator().Extract(c.Context(), propagation.MapCarrier(cleanHeaders))
		_, sp := tracer.Start(ctx, "handler.auth.health-check")
		defer sp.End()
		time.Sleep(50 * time.Millisecond)
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})

	g.Go(func() error {
		log.Printf("[auth][info]: app listening at %s \n", cfg.GetAppAddress())
		return app.Listen(cfg.GetAppAddress())
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Println("[auth][info]: signal graceful shutdown recived")
		gracfulCtx, gracefulStop := context.WithTimeout(context.Background(), time.Second*30)
		defer gracefulStop()
		return app.ShutdownWithContext(gracfulCtx)
	})
	if err := g.Wait(); err != nil {
		log.Printf("[auth][error]: group exited with error: %v \n", err)
	}
	log.Println("[auth][info]: graceful shutdown success.")
}
