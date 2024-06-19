package main

import (
	"context"
	"fmt"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func main() {
	config := fiber.Config{}

	tracer := otel.Tracer("test-server")
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	hdsvc := NewDefaultHeavHeavyDutyService()
	myHandler := NewMyHandler(tracer, hdsvc)

	app := fiber.New(config)
  app.Use(otelfiber.Middleware())
	app.Get("/error", myHandler.HandleGetError)
	app.Get("/ff", myHandler.HandleGetFireForget)

	app.Listen(":3333")
}
