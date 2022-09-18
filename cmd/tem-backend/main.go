package main

import (
	"context"
	"github.com/Coflnet/tem-backend/internal/api"
	"github.com/Coflnet/tem-backend/internal/mongo"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"os"
)

func main() {

	cleanup := initTracer()
	defer func(c func(context.Context) error) {
		err := c(context.Background())
		if err != nil {
			log.Error().Err(err).Msgf("failed to shutdown exporter")
			return
		}
	}(cleanup)

	mongo.Start()
	defer func() {
		err := mongo.Stop()
		if err != nil {
			log.Error().Err(err).Msgf("failed to stop mongo")
		}
	}()

	err := api.StartApi()

	if err != nil {
		log.Panic().Err(err).Msgf("Error starting API")
	}
}

func initTracer() func(context.Context) error {
	secureOption := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(getEnv("JAEGER_AGENT_HOST")),
		),
	)

	if err != nil {
		log.Error().Err(err).Msgf("failed to create exporter")
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "tem-backend"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func getEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panic().Msgf("%s not set", k)
	}
	return v
}
