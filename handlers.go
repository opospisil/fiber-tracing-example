package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	ott "go.opentelemetry.io/otel/trace"
)

type MyHandler struct {
	tracer ott.Tracer
	svc    HeavyDutyService
}

func NewMyHandler(t ott.Tracer, svc HeavyDutyService) *MyHandler {
	return &MyHandler{
		tracer: t,
		svc:    svc,
	}
}

func (h *MyHandler) HandleGetError(c *fiber.Ctx) error {
	return fmt.Errorf("AV - artificial vyjeb")
}

func (h *MyHandler) HandleGetFireForget(c *fiber.Ctx) error {
	_, span := h.tracer.Start(c.Context(), "getFireForget", ott.WithAttributes(attribute.String("testAttribute", "something")))
	defer span.End()
	h.svc.DoSomethingExpansive(c.Context())
	msg := "Request to the void fired"
	c.Write([]byte(msg))
	return nil
}
