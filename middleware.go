package main

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type HDSTracingMw struct {
	tracer trace.Tracer
	next   HeavyDutyService
}

func NewHDSTracingMw(svc HeavyDutyService, t trace.Tracer) HeavyDutyService {
	return &HDSTracingMw{
		tracer: t,
		next:   svc,
	}
}

func (mw *HDSTracingMw) ExpansiveOpInBackground(c context.Context, t trace.Tracer) error {
  ctx, span := t.Start(c, "Expansive background op from middleware")
  defer span.End()
	return mw.next.ExpansiveOpInBackground(ctx, mw.tracer)
}

func (mw *HDSTracingMw) BlockingExpansiveOp(c context.Context) error {
  ctx, span := mw.tracer.Start(c, "Blocking heavy op from mw")
  defer span.End()
	return mw.next.BlockingExpansiveOp(ctx)
}
