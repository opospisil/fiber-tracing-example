package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type HeavyDutyService interface {
	ExpansiveOpInBackground(context.Context, trace.Tracer) error
	BlockingExpansiveOp(context.Context) error
}

type DefaultHeavyDutyService struct{}

func (hds *DefaultHeavyDutyService) ExpansiveOpInBackground(c context.Context, t trace.Tracer) error {
	ctx, cancelFn := context.WithTimeout(c, time.Second*4)

	go func() {
		_, span := t.Start(c, "background_operation")
		err := ActualHeavyOp(ctx)
		if err != nil {
			fmt.Printf("goroutine finished with an error: %v\n", err)
		}
		if err == nil {
			fmt.Println("This one was success")
		}
		span.End()
		cancelFn()
	}()
	return nil
}

func (hds *DefaultHeavyDutyService) BlockingExpansiveOp(c context.Context) error {
	ctx, cancelFn := context.WithTimeout(c, time.Second*4)
	defer cancelFn()
	return ActualHeavyOp(ctx)
}

// simulates random response time as well as random failures for different reasons
func ActualHeavyOp(c context.Context) error {
	n := rand.Intn(10)
	fmt.Printf("Starting fire forget for number %v\n", n)
	for range n {
		select {
		case <-c.Done():
			return c.Err()
		default:
			time.Sleep(time.Second)
		}
	}
	if n%2 == 0 {
		return nil
	}
	return fmt.Errorf("received %v, failing", n)
}

func NewDefaultHeavHeavyDutyService() *DefaultHeavyDutyService {
	return &DefaultHeavyDutyService{}
}
