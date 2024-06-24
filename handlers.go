package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type MyHandler struct {
	tracer trace.Tracer
	svc    HeavyDutyService
}

func NewMyHandler(t trace.Tracer, svc HeavyDutyService) *MyHandler {
	return &MyHandler{
		tracer: t,
		svc:    svc,
	}
}

func (h *MyHandler) HandleGetError(c *fiber.Ctx) error {
	return fmt.Errorf("AV - artificial vyjeb")
}

func (h *MyHandler) HandleGetBackgroundOp(c *fiber.Ctx) error {
	h.svc.ExpansiveOpInBackground(c.Context(), h.tracer)
	msg := "Request to the void fired"
	c.Write([]byte(msg))
	return nil
}

func (h *MyHandler) HandleGetBlocking(c *fiber.Ctx) error {
  err := h.svc.BlockingExpansiveOp(c.Context())
  if err != nil {
    return err
  }
  c.Write([]byte("Blocking op finished successfully"))
  return nil
}
