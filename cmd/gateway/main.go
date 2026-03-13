package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/gofiber/fiber/v3/middleware/recover"
	dotenvx "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/adapter/config"
	keycloak "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/adapter/iam"
	gateway_middleware "github.com/premwitthawas/demo_ecommerce_api/internals/gateway/adapter/middleware"
	"github.com/premwitthawas/demo_ecommerce_api/internals/gateway/domain/iam"
	"github.com/premwitthawas/demo_ecommerce_api/internals/gateway/port/config"
	pkg_otel "github.com/premwitthawas/demo_ecommerce_api/pkgs/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

var tracer = otel.Tracer("http-handler")

type BootStrapApplication struct {
	app *fiber.App
	cfg config.GatewayConfigAdapter
	tp  trace.Tracer
}

func main() {
	tp := pkg_otel.SetupOtelTracer("localhost:4317", "gateway-service")
	defer func() {
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	g, gctx := errgroup.WithContext(ctx)
	cfg := dotenvx.NewGatewayConfig()
	gateway := NewApplication()
	bootstrap := NewBoostrapApplication(gateway, cfg, tracer)
	bootstrap.setup()
	g.Go(func() error {
		log.Printf("[gateway][info]: gateway listening at %s \n", cfg.GetAppConfig().Address)
		return gateway.Listen(cfg.GetAppConfig().Address)
	})
	g.Go(func() error {
		<-gctx.Done()
		log.Println("[gateway][info]: graceful signal recived.")
		sctx, stop := context.WithTimeout(context.Background(), time.Second*30)
		defer stop()
		return gateway.ShutdownWithContext(sctx)
	})
	if err := g.Wait(); err != nil {
		log.Fatalf("[gateway][error]: group error: %v \n", err)
	}
	log.Println("[gateway][info]: graceful shutdown success.")
}

func registerProxy(r fiber.Router, targer string) {
	r.Use(proxy.Balancer(proxy.Config{
		Servers:              []string{targer},
		KeepConnectionHeader: true,
		ModifyRequest: func(c fiber.Ctx) error {
			carrier := propagation.HeaderCarrier{}
			headers := c.Request().Header.All()
			for k, v := range headers {
				carrier[string(k)] = []string{string(v)}
			}
			otel.GetTextMapPropagator().Inject(c.Context(), carrier)
			for k, v := range carrier {
				if len(v) > 0 {
					c.Request().Header.Set(k, v[0])
				}
			}
			return nil
		},
		ModifyResponse: func(c fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			c.Response().Header.Del("x-powered-by")
			return nil
		},
	}))
}

func NewApplication() *fiber.App {
	app := fiber.New()
	return app
}

func NewBoostrapApplication(app *fiber.App, cfg config.GatewayConfigAdapter, tp trace.Tracer) *BootStrapApplication {
	return &BootStrapApplication{
		app: app,
		cfg: cfg,
		tp:  tp,
	}
}

func (b *BootStrapApplication) setup() {
	b.setupGlobalMiddleware()
	b.setupProtectRoutes()
}

func (b *BootStrapApplication) setupGlobalMiddleware() {
	b.app.Use(recover.New())
	b.app.Use(logger.New())
	b.app.Use(fiberotel.Middleware())
	b.app.Use(helmet.New())
}

func (b *BootStrapApplication) setupProtectRoutes() {
	keycloakIAM, err := keycloak.NewKeycloakAdapter(context.Background(), b.cfg)
	if err != nil {
		log.Printf("Error create Keycloak IAM failure: %v", err)
	}
	authAPI := b.app.Group("/api/v1/auth")
	authAPI.Use(gateway_middleware.AuthMiddleware(b.cfg, keycloakIAM, b.tp))
	authAPI.Use(gateway_middleware.RBACMiddleware(b.cfg, b.tp, "admin"))
	authAPI.All("/*", proxy.Balancer(proxy.Config{
		Servers: []string{"http://127.0.0.1:5001"},
		ModifyRequest: func(c fiber.Ctx) error {
			headers := make(propagation.HeaderCarrier)
			otel.GetTextMapPropagator().Inject(c.Context(), headers)
			for k, v := range headers {
				c.Request().Header.Set(k, string(v[0]))
			}
			if claims, ok := c.Locals("user_claims").(*iam.Claims); ok {
				c.Request().Header.Set("X-User-ID", claims.Subject)
			}
			return nil
		},
	}))
	// authAPI.All("/*", func(c fiber.Ctx) error {
	// 	headers := make(propagation.HeaderCarrier)
	// 	otel.GetTextMapPropagator().Inject(c.Context(), headers)
	// 	for k, v := range headers {
	// 		c.Request().Header.Set(k, string(v[0]))
	// 	}
	// 	targetURL := "http://127.0.0.1:5001" + c.Path()
	// 	return proxy.Do(c, targetURL)
	// })
}
