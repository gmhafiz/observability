package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	skdmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// SetupOTLPExporter reference in https://github.com/open-telemetry/opentelemetry-go/blob/main/example/dice/otel.go
func SetupOTLPExporter(ctx context.Context, cfg OpenTelemetry) func() {
	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.OtlpServiceName),
			semconv.ServiceVersionKey.String(cfg.OtlpServiceVersion),
		),
	)
	if err != nil {
		log.Println(fmt.Errorf("creating resource, %v", err))
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.OtlpEndpoint),
		//otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	sctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	traceExp, err := otlptrace.New(sctx, traceClient)
	if err != nil {
		log.Println(fmt.Errorf("creating trace client, %v", err))
	}
	bsp := trace.NewBatchSpanProcessor(traceExp)
	tracerProvider := trace.NewTracerProvider(
		//trace.WithSampler(trace.AlwaysSample()),
		trace.WithSampler(trace.TraceIDRatioBased(cfg.OtlpSamplerRatio)),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(cfg.OtlpEndpoint),
	)
	if err != nil {
		log.Println(fmt.Errorf("creating metric exporter, %v", err))
	}
	meterProvider := skdmetric.NewMeterProvider(
		skdmetric.WithResource(res),
		skdmetric.WithReader(
			skdmetric.NewPeriodicReader(
				metricExp,
				skdmetric.WithInterval(5*time.Second),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
		//pushes any last exports to the receiver
		if err := meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
