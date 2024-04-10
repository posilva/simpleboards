package telemetry

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func init() {
	// exporter, err := otlptracegrpc.New(ctx)
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	handleErr(err, "Failed to create exporter")

	openTelemetryURL := attribute.KeyValue{
		Key:   attribute.Key("opentelemetry.io/schemas"),
		Value: attribute.StringValue("1.24.0"),
	}

	resource, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			openTelemetryURL,
			semconv.ServiceName("Leaderboard Service"),
			semconv.ServiceVersion("v0.0.1"),
		))

	handleErr(err, "Failed to create resource")

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func handleErr(err error, message string) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", message, err))
	}
}

type DefaultTelemetryReporter struct{}

func NewDefaultTelemetryReporter() *DefaultTelemetryReporter {
	return &DefaultTelemetryReporter{}
}

func (t *DefaultTelemetryReporter) SetDefaultTags(tags map[string]string) {
}

func (t *DefaultTelemetryReporter) ReportGauge(name string, value float64, tags map[string]string) {
}

func (t *DefaultTelemetryReporter) ReportCounter(name string, value float64, tags map[string]string) {
}

func (t *DefaultTelemetryReporter) ReportHistogram(name string, value float64, tags map[string]string) {
}

func (t *DefaultTelemetryReporter) ReportSummary(name string, value float64, tags map[string]string) {
}
