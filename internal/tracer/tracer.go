package tracer

import (
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func New(cfg *Config, namespace string, subsystem string) trace.Tracer {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if !cfg.Enabled {
		exporter, err = stdout.New(stdout.WithPrettyPrint())
	} else {
		exporter, err = jaeger.New(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(cfg.Host),
				jaeger.WithAgentPort(strconv.Itoa(cfg.Port)),
			),
		)
	}

	if err != nil {
		log.Fatalf("failed to initialize export pipeline: %v", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNamespaceKey.String(namespace),
			semconv.ServiceNameKey.String(subsystem),
		),
	)
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRate))),
	)

	otel.SetTracerProvider(tp)

	// register the TraceContext propagator globally.
	var tc propagation.TraceContext

	otel.SetTextMapPropagator(tc)

	return otel.Tracer("snapp-incubator/ghodrat")
}
