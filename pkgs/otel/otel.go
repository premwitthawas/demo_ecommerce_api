package pkg_otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	exporter "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func SetupOtelTracer(url, serviceName string) *sdktrace.TracerProvider {
	jaegerExporter, err := exporter.New(context.Background(), exporter.WithEndpoint(url), exporter.WithInsecure())
	if err != nil {
		log.Fatalf("[otel][error]: create exporter failure : % v \n", err)
	}
	resources, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("1.0.0"),
	))
	if err != nil {
		log.Fatalf("[otel][error]: create resources failure : % v \n", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithResource(resources),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
